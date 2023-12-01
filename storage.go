package main

import (
	"sync"

	data "github.com/fops9311/wbl0_231124/data"
)

var m sync.Mutex

var inmem = make(map[string]data.RawOrderData)
var inmemkeys = make([]string, 0)

func StorageGet(id string) (data.RawOrderData, bool) {
	m.Lock()
	defer m.Unlock()
	v, ok := inmem[id]
	return v, ok
}

func StorageUpdate(id string, val data.RawOrderData) {
	m.Lock()
	defer m.Unlock()
	inmem[id] = val
	inmemkeys = append(inmemkeys, id)
}
func StorageGetKeys() []string {
	return inmemkeys
}
