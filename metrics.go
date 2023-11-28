package main

import (
	"github.com/prometheus/client_golang/prometheus"
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
