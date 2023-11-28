package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

var ErrNatsNotInited error = errors.New("nats not inited")
var ErrNatsSubNotFound error = errors.New("nats subscription not found")
var n nts = nts{
	subs: make(map[string]stan.Subscription),
}

type nts struct {
	nc   *nats.Conn
	sc   stan.Conn
	subs map[string]stan.Subscription
}

func (n *nts) Validate() error {
	if n.nc == nil {
		return ErrNatsNotInited
	}
	if n.subs == nil {
		return ErrNatsNotInited
	}
	return nil
}

func InitNatsSubscriber(url string, options ...nats.Option) error {

	var err error
	// Connect to a server
	n.nc, err = nats.Connect(url, options...)
	n.nc.Opts.AsyncErrorCB = func(c *nats.Conn, s *nats.Subscription, err error) {
		fmt.Println("nats: slow consumer, messages dropped")
		ordersErrors.Inc()
	}
	if err != nil {
		return err
	}

	n.sc, err = stan.Connect("TESTCLUSTER", "clientID", stan.NatsConn(n.nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, url)
		return err
	}

	return err
}

type NatsSubscribeOpts struct {
	Topic   string
	MsgChan chan *stan.Msg
}

func NatsSubscribe(opts NatsSubscribeOpts) error {
	if err := n.Validate(); err != nil {
		return err
	}
	sub, _ := n.sc.Subscribe(opts.Topic, func(m *stan.Msg) {
		opts.MsgChan <- m
	}, stan.DurableName("my-durable"))

	n.subs[opts.Topic] = sub
	return nil
}
func NatsUnsubscribe(topic string) error {
	if err := n.Validate(); err != nil {
		return err
	}
	if sub, ok := n.subs[topic]; ok {
		return sub.Unsubscribe()
	}
	return ErrNatsSubNotFound
}
