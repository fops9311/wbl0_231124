package main

import (
	"fmt"
	"time"

	storage "github.com/fops9311/wbl0_231124/internal/memorystorage"
)

func main() {
	s := storage.NewInmemStorage[bool]()
	go s.StorageUpdate("test", true)
	go s.StorageUpdate("test1", true)
	go s.StorageUpdate("test", true)
	go s.StorageUpdate("test1", true)
	<-time.After(time.Second)
	fmt.Println(s.StorageGetKeys())
	fmt.Println(s.StorageGetAll())
	fmt.Println(s.StorageGet("test"))
	fmt.Println(s.StorageGet("test23"))
}
