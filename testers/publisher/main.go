package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	data "github.com/fops9311/wbl0_231124/data"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

var NATSHOST = "nats://localhost:4223"
var PUBTOPIC = "TESTING"
var inmemcache = make(map[string][]byte)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	d := data.RawOrderData{}
	json.Unmarshal([]byte(TESTMESSAGE), &d)
	// Connect to a server
	nc, err := nats.Connect(NATSHOST)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Connect to a server

	sc, err := stan.Connect("TESTCLUSTER", "clientID33", stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, NATSHOST)
	}

	// Simple Publisher
	func() {
		//variate data a bit
		d.SmId = int(time.Now().UnixNano())
		d.DateCreated = time.Now().Format("2006-01-02T15:04:05.999999999Z07:00")
		b, err := json.Marshal(d)
		if err != nil {
			fmt.Println(err)
		}

		for nanosecs := 20000000; nanosecs > 10; nanosecs-- {
			select {
			case <-sigs:
				fmt.Println("INTERUPTED")
				return
				//case <-time.NewTimer(time.Millisecond * 1).C:
			default:
				//nanosecs--
				sc.Publish(PUBTOPIC, b)
				//fmt.Println(d.DateCreated, "SEND", len(b))
			}
		}

	}()
	nc.Close()
	fmt.Println("FINISHED")
}

var TESTMESSAGE = `{
	"order_uid": "b563feb7b2b84b6test",
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
	  "name": "Test Testov",
	  "phone": "+9720000000",
	  "zip": "2639809",
	  "city": "Kiryat Mozkin",
	  "address": "Ploshad Mira 15",
	  "region": "Kraiot",
	  "email": "test@gmail.com"
	},
	"payment": {
	  "transaction": "b563feb7b2b84b6test",
	  "request_id": "",
	  "currency": "USD",
	  "provider": "wbpay",
	  "amount": 1817,
	  "payment_dt": 1637907727,
	  "bank": "alpha",
	  "delivery_cost": 1500,
	  "goods_total": 317,
	  "custom_fee": 0
	},
	"items": [
	  {
		"chrt_id": 9934930,
		"track_number": "WBILMTESTTRACK",
		"price": 453,
		"rid": "ab4219087a764ae0btest",
		"name": "Mascaras",
		"sale": 30,
		"size": "0",
		"total_price": 317,
		"nm_id": 2389212,
		"brand": "Vivienne Sabo",
		"status": 202
	  }
	],
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
  }
`
