package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// параметры http запроса
type HttpRequestParams struct {
	PageId  int
	OrderId string
}

// хэндлер возвращает json с общим количеством сохраненных значений
func HandleMemLen(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var l int = len(StorageGetKeys())
	w.Write([]byte(fmt.Sprintf("{\"mem_len\": %d}", l)))
}

// хэндлер возвращает json отображение запрошенного заказа
func HandleGetOrder(p HttpRequestParams) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res, ok := StorageGet(p.OrderId)
		if !ok {
			http.Error(w, "404 resource not found", http.StatusNotFound)
			return
		}
		b, err := json.Marshal(&res)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}

// хэндлер возвращает json страницы заказов pagination
func HandleGetOrderPage(p HttpRequestParams) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, code := GetPage(p.PageId)
		if code != http.StatusOK {
			http.Error(w, http.StatusText(code), code)
			return
		}
		w.Write(body)
	}
}

// хэлпер функция для получения страницы. Если страницы не существует, возвращает код http отличный от 200
func GetPage(pageId int) ([]byte, int) {
	keys := StorageGetKeys()
	var pageVolume int = 50
	len := len(keys)
	var lastPage int = len/pageVolume + 1
	if pageId < 0 {
		pageId = lastPage
	}
	if len == 0 {
		return []byte{}, http.StatusNotFound
	}
	if pageId > lastPage {
		return []byte{}, http.StatusNotFound
	}
	if pageId < 1 {
		return []byte{}, http.StatusBadRequest
	}
	page := make([]string, pageVolume)
	if (pageId)*pageVolume < len {
		page = keys[(pageId-1)*pageVolume : (pageId)*pageVolume]
	} else {
		page = keys[(pageId-1)*pageVolume:]
	}
	s := strings.Join(page, "\",\"")
	s = fmt.Sprintf("[\"%s\"]", s)
	return []byte(s), http.StatusOK
}
