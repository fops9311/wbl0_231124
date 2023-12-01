//go:generate go run ./scripts/analize_struct/.
//go:generate docker build -t wb_nats_demo:latest .
//go:generate docker-compose down
//go:generate docker-compose up --remove-orphans -d
package main

import (
	"fmt"
	"log"
	"net/http"
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

	memwritechan := make(chan OrderWithKey, 10)
	ChanFanIn(
		[]chan OrderWithKey{
			memchan,
			initmemchan,
		},
		memwritechan,
	)
	ChanConsumer(memwritechan, func(d OrderWithKey) error {
		StorageUpdate(d.Key.Get(), d.Val)
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
			h = get(HandleMemLen)
		case match(p, "/order_pages/last"):
			h = get(HandleGetOrderPage(HttpRequestParams{PageId: -1}))
		case match(p, "/order_pages/+", &pageId):
			h = get(HandleGetOrderPage(HttpRequestParams{PageId: pageId}))
		case match(p, "/order/+", &id):
			h = get(HandleGetOrder(HttpRequestParams{OrderId: id}))
		default:
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", Serve)
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
