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
	stan "github.com/nats-io/stan.go"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// данные о заказе с ключем
type OrderWithKey struct {
	Key Key
	Val data.RawOrderData
}

// ключ служит для хранения в бд и кэше
type Key string

// получить значение
func (k *Key) Get() string {
	return string(*k)
}

// установить значение
func (k *Key) Set(s string) {
	*k = Key(s)
}

// сгенерировать ключ
func (k *Key) Generate() {
	k.Set(fmt.Sprintf("%d", time.Now().UnixNano()))
}

func main() {
	//инициализация из env
	InitGlobalVarsFromEnv()
	//инизацилация таблице в бд
	fmt.Println("INIT TABLE RESULT", initTable())
	fmt.Println("INIT FROM DB")
	//инизациализация подписки на топик
	initmemchan := make(chan OrderWithKey, 10)
	go selectOrders(initmemchan)
	err := InitNatsSubscriber(NATSHOST)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("INITED")
	natschan := make(chan *stan.Msg, 20)
	err = NatsSubscribe(NatsSubscribeOpts{
		Topic:   PUBTOPIC,
		MsgChan: natschan,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SUBSCRIBED")
	//формирование конвейера
	//											   initmemchan -> |(fin)
	//															  |
	//natschan -> msgdatachan -> orderschan |(fout) -> memchan -> |(fin) -> memwritechan (inmemcacheConsumer)
	//										|
	//									  	|(fout) -> encodechan -> gobchan (databaseConsumer)
	//преобразование каналов
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
	//http сервер
	Serve := func(w http.ResponseWriter, r *http.Request) {
		var h http.Handler
		var id string
		var pageId int

		p := r.URL.Path

		switch {
		case match(p, "/mem_len"):
			h = get(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				m.Lock()

				var l int = len(inmem)
				m.Unlock()
				w.Write([]byte(fmt.Sprintf("{\"mem_len\": %d}", l)))
			})
		case match(p, "/order_pages/last"):
			h = get(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")

					m.Lock()
					defer m.Unlock()
					var pageVolume int = 50
					len := len(inmemkeys)
					var lastPage int = len/pageVolume + 1
					pageId = lastPage
					page := make([]string, pageVolume)
					if len == 0 {
						http.Error(w, "404 resource not found", http.StatusNotFound)
						return
					}
					if (pageId)*pageVolume < len {
						page = inmemkeys[(pageId-1)*pageVolume : (pageId)*pageVolume]
					} else {
						page = inmemkeys[(pageId-1)*pageVolume:]
					}
					s := strings.Join(page, "\",\"")
					w.Write([]byte(fmt.Sprintf("[\"%s\"]", s)))
				})
		case match(p, "/order_pages/+", &pageId):
			h = get(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")

					m.Lock()
					defer m.Unlock()
					var pageVolume int = 50
					len := len(inmemkeys)
					var lastPage int = len/pageVolume + 1
					if len == 0 {
						http.Error(w, "404 resource not found", http.StatusNotFound)
						return
					}
					if pageId > lastPage {
						http.Error(w, "404 resource not found", http.StatusNotFound)
						return
					}
					if pageId < 1 {
						http.Error(w, "400 bad request", http.StatusBadRequest)
						return
					}
					page := make([]string, pageVolume)
					if (pageId)*pageVolume < len {
						page = inmemkeys[(pageId-1)*pageVolume : (pageId)*pageVolume]
					} else {
						page = inmemkeys[(pageId-1)*pageVolume:]
					}
					s := strings.Join(page, "\",\"")
					w.Write([]byte(fmt.Sprintf("[\"%s\"]", s)))
				})
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
