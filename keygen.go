package main

import (
	"fmt"
	"time"

	"github.com/fops9311/wbl0_231124/data"
)

// данные о заказе с ключем
type OrderWithKey struct {
	Key Key
	Val data.RawOrderData
}

// gob encoded data с ключем
type GobWithKey struct {
	Key Key
	Val []byte
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
