//go:generate go run ./scripts/analize_struct/.
//go:generate docker build -t wb_nats_demo:latest .
//go:generate docker-compose down
//go:generate docker-compose up --remove-orphans -d
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	data "github.com/fops9311/wbl0_231124/data"
	"github.com/nats-io/nats.go"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ordersHandled = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_orders_handled_total",
		Help: "The total number of orders handled",
	})
	ordersWriteDb = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_orders_write_db",
		Help: "The number of orders write to db",
	})
	ordersReadDb = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_orders_read_db",
		Help: "The number of orders read to db",
	})
	ordersErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_orders_errors",
		Help: "The number of orders with errors",
	})
	ordersInMem = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_orders_inmemcache",
		Help: "The number of orders in memory cache",
	})
)
var NATSHOST = "nats://nats:4223"
var PUBTOPIC = "TESTING"

type OrderWithKey struct {
	Key Key
	Val data.RawOrderData
}
type Key string

func (k *Key) Get() string {
	return string(*k)
}
func (k *Key) Set(s string) {
	*k = Key(s)
}
func (k *Key) Generate() {
	k.Set(fmt.Sprintf("%d", time.Now().UnixNano()))
}

func main() {
	fmt.Println("INIT FROM DB")

	initmemchan := make(chan OrderWithKey, 10)
	go selectOrders(initmemchan)
	err := InitNatsSubscriber(NATSHOST)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("INITED")
	natschan := make(chan *nats.Msg, 20)
	err = NatsSubscribe(NatsSubscribeOpts{
		Topic:   PUBTOPIC,
		MsgChan: natschan,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SUBSCRIBED")

	//											   initmemchan -> |(fin)
	//															  |
	//natschan -> msgdatachan -> orderschan |(fout) -> memchan -> |(fin) -> memwritechan (inmemcacheConsumer)
	//										|
	//									  	|(fout) -> encodechan -> gobchan (databaseConsumer)
	msgdatachan := TranformChan(natschan, MessageToByteSlice, 10)
	orderschan := TranformChan(msgdatachan, ByteSliceToOrderData, 10)

	memchan := make(chan OrderWithKey, 10)
	encodechan := make(chan OrderWithKey, 10)
	ChanFanOut(orderschan, []chan OrderWithKey{
		memchan,
		encodechan,
	})
	var m sync.Mutex
	inmem := make(map[string]data.RawOrderData)
	inmemkeys := make([]string, 0)
	getInmemById := func(id string) (data.RawOrderData, bool) {
		m.Lock()
		v, ok := inmem[id]
		m.Unlock()
		return v, ok
	}
	updateInmemOrder := func(d OrderWithKey) {
		ordersInMem.Inc()
		m.Lock()
		defer m.Unlock()
		inmem[d.Key.Get()] = d.Val
		inmemkeys = append(inmemkeys, d.Key.Get())
	}

	memwritechan := make(chan OrderWithKey, 10)
	ChanFanIn(
		[]chan OrderWithKey{
			memchan,
			initmemchan,
		},
		memwritechan,
	)
	ChanConsumer(memwritechan, func(d OrderWithKey) error {
		updateInmemOrder(d)
		return nil
	})
	gobchan := TranformChan(encodechan, OrderDataToGob, 10)
	ChanConsumer(gobchan, func(d GobWithKey) error {
		return insertOrder(d.Key.Get(), string(d.Val))
	})

	memLenHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		m.Lock()

		var l int = len(inmem)
		m.Unlock()
		w.Write([]byte(fmt.Sprintf("{\"mem_len\": %d}", l)))
	}
	memLastHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		m.Lock()
		var s string
		if len(inmemkeys) > 10 {
			s = strings.Join(inmemkeys[len(inmemkeys)-10:], "\",\"")
		}
		m.Unlock()
		w.Write([]byte(fmt.Sprintf("[\"%s\"]", s)))
	}

	Serve := func(w http.ResponseWriter, r *http.Request) {
		var h http.Handler
		var id string

		p := r.URL.Path

		switch {
		case match(p, "/mem_len"):
			h = get(memLenHandler)
		case match(p, "/mem_last"):
			h = get(memLastHandler)
		case match(p, "/order/+", &id):
			h = get(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					res, ok := getInmemById(id)
					if !ok {
						http.Error(w, "404 resource not found", http.StatusNotFound)
						return
					}
					b, err := json.Marshal(&res)
					if err != nil {
						http.Error(w, "500 internal server error", http.StatusInternalServerError)
						return
					}
					w.Write(b)
				})
		default:
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}
	http.HandleFunc("/", Serve)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
