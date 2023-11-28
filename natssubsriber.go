package main

import (
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
)

var ErrNatsNotInited error = errors.New("nats not inited")
var ErrNatsSubNotFound error = errors.New("nats subscription not found")
var n nts = nts{
	subs: make(map[string]*nats.Subscription),
}

type nts struct {
	nc   *nats.Conn
	subs map[string]*nats.Subscription
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
	return err
}

type NatsSubscribeOpts struct {
	Topic   string
	MsgChan chan *nats.Msg
}

func NatsSubscribe(opts NatsSubscribeOpts) error {
	if err := n.Validate(); err != nil {
		return err
	}
	sub, err := n.nc.ChanSubscribe(opts.Topic, opts.MsgChan)
	if err != nil {
		return err
	}
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
