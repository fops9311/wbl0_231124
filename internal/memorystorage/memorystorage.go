package memorystorage

import (
	"sync"

	storage "github.com/fops9311/wbl0_231124/internal/storage"
)

type memcache struct {
	m    sync.Mutex
	data map[string]interface{}
}

func (mc *memcache) StorageGet(id string) (interface{}, bool) {
	mc.m.Lock()
	defer mc.m.Unlock()
	v, ok := mc.data[id]
	return v, ok

}
func (mc *memcache) StorageGetAll() []interface{} {
	res := make([]interface{}, len(mc.data))
	count := 0
	mc.m.Lock()
	defer mc.m.Unlock()
	for k := range mc.data {
		res[count] = mc.data[k]
		count++
	}
	return res
}
func (mc *memcache) StorageUpdate(id string, val interface{}) {
	mc.m.Lock()
	defer mc.m.Unlock()
	mc.data[id] = val
}
func (mc *memcache) StorageGetKeys() []string {
	res := make([]string, len(mc.data))
	count := 0
	mc.m.Lock()
	defer mc.m.Unlock()
	for k := range mc.data {
		res[count] = k
		count++
	}
	return res
}

func (mc *memcache) StorageInit() {
	mc.m.Lock()
	defer mc.m.Unlock()
	mc.data = make(map[string]interface{})
}

func NewInmemStorage[T any]() *storage.AnyStorage[T] {
	return storage.NewStorage[T](&memcache{})
}
