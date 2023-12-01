package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var TESTJSON = `{
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
var structs = "package data\n"

func apnd(s string) {
	structs += s
}
func c(m1 interface{}, id string) {
	switch m := m1.(type) {
	case map[string]interface{}:
		//defer apnd(fmt.Sprintf("//--------------map[string]interface{}-%s------------\n", id))

		defer apnd(fmt.Sprintf("}\n"))
		defer apnd(fmt.Sprintf("LocalID string `json:\"-\"`\n"))
		for k, v := range m {
			StructElemName := toCammelCase(k)
			StructElemType := fmt.Sprintf("%T", v)
			StructElemJsonName := k
			if StructElemType == "map[string]interface {}" {
				StructElemType = fmt.Sprintf("Raw%sData", StructElemName)
			}
			if StructElemType == "[]interface {}" {
				StructElemType = fmt.Sprintf("[]Raw%sData", StructElemName)
			}
			defer apnd(fmt.Sprintf("	%s %s `json:\"%s\"`\n", StructElemName, StructElemType, StructElemJsonName))
			c(v, k)
		}
		defer apnd(fmt.Sprintf("type Raw%sData struct{\n", toCammelCase(id)))
		defer apnd(fmt.Sprintf("//GENERATED json %s data\n", id))
	case []interface{}:
		for _, v := range m {
			c(v, id)
			break
		}
	}

}
func main() {

	var data interface{}
	json.Unmarshal([]byte(TESTJSON), &data)
	structs += `
import "encoding/gob"

func init() { gob.Register(RawOrderData{}) }
	`
	c(data, "Order")
	structs = strings.ReplaceAll(structs, "float64", "int")
	os.WriteFile(filepath.Join(".", "data", "jsonorder.go"), []byte(structs), 0644)

}
func capitalize(s string) string {
	b := []byte(s)
	if len(s) == 0 {
		return ""
	}
	b[0] = (strings.ToUpper(string(b[0])))[0]
	return string(b)
}
func toCammelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = capitalize(parts[i])
	}
	return strings.Join(parts, "")
}
