package main

import (
	"fmt"
	"os"
	"strconv"
)

var NATSHOST = "nats://nats:4223"
var PUBTOPIC = "TESTING"

var (
	PGHOST     = "demo_postgres"
	PGPORT     = 5432
	PGUSER     = "postgres"
	PGPASSWORD = "password"
	PGDBNAME   = "demo1"
)

func InitGlobalVarsFromEnv() {
	NATSHOST = fmt.Sprintf("nats://%s:%s",
		getStringEnv("NATS_HOST", "nats"),
		getStringEnv("NATS_PORT", "4223"),
	)
	PUBTOPIC = getStringEnv("PUB_TOPIC", "TESTING")
	PGHOST = getStringEnv("PG_HOST", "demo_postgres")
	PGPORT = getIntEnv("PG_PORT", 5432)
	PGUSER = getStringEnv("PG_USER", "postgres")
	PGPASSWORD = getStringEnv("PG_PASSWORD", "password")
	PGDBNAME = getStringEnv("PG_DBNAME", "demo1")
}
func getStringEnv(name string, def string) string {
	s := os.Getenv(name)
	if s == "" {
		return def
	}
	return s
}
func getIntEnv(name string, def int) int {
	s := os.Getenv(name)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return int(i)
}
