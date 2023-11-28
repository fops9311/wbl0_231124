package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"log"

	data "github.com/fops9311/wbl0_231124/data"
	stan "github.com/nats-io/stan.go"
)

// функция для вызова горутин преобразования каналов
func TranformChan[T1 any, T2 any](in chan T1, f func(T1) (T2, error), size int) chan T2 {
	result := make(chan T2, size)

	go func() {
		for m := range in {
			r, err := f(m)
			if err != nil {
				log.Println(err)
				ordersErrors.Inc()
				continue
			}
			result <- r
		}
	}()

	return result
}

// функция для вызова горутин fanout
func ChanFanOut[T1 any](in chan T1, out []chan T1) {
	go func() {
		for t := range in {
			for i := range out {
				out[i] <- t
			}
		}
	}()
}

// функция для вызова горутин fanin
func ChanFanIn[T1 any](in []chan T1, out chan T1) {
	for i := range in {
		ind := i
		go func() {
			for v := range in[ind] {
				out <- v
			}
		}()
	}
}

// функция для вызова горутин потребителя канала
func ChanConsumer[T1 any](in chan T1, f func(T1) error) {
	go func() {
		for c := range in {
			err := f(c)
			if err != nil {
				ordersErrors.Inc()
				log.Println(err)
			}
		}
	}()
}

var ErrGotNilMessage = errors.New("got nil message")

func MessageToByteSlice(msg *stan.Msg) ([]byte, error) {
	ordersHandled.Inc()
	if msg == nil {
		return []byte{}, ErrGotNilMessage
	}
	//fmt.Println(strings.ToLower("MessageToByteSlice: "), len(msg.Data))
	return msg.Data, nil
}
func ByteSliceToOrderData(b []byte) (OrderWithKey, error) {
	var d *data.RawOrderData = &data.RawOrderData{}
	err := json.Unmarshal(b, d)
	if err != nil {
		return OrderWithKey{}, errors.New("bytes to order data:" + err.Error())
	}
	order := OrderWithKey{Val: *d}
	order.Key.Generate()
	order.Val.LocalID = order.Key.Get()

	return order, nil
}

// флаг регистрации gob
var isOrderTypeRegistered = false

func OrderDataToGob(d OrderWithKey) (GobWithKey, error) {
	order := d.Val
	if !isOrderTypeRegistered {
		gob.Register(order)
		isOrderTypeRegistered = true
	}
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(&order)
	if err != nil {
		return GobWithKey{}, errors.New("order data to gob:" + err.Error())
	}
	g := GobWithKey{Val: buff.Bytes()}
	g.Key.Set(d.Key.Get())
	return g, nil
}

type GobWithKey struct {
	Key Key
	Val []byte
}
