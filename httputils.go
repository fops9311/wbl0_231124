package main

import (
	"net/http"
	"strconv"
	"strings"
)

//Набор утилит для организации http сервера
//Мне показался такой подход интересным, плюс он довольно быстрый согласно бенчмаркам

// get takes a HandlerFunc and wraps it to only allow the GET method
func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "GET")
}

// post takes a HandlerFunc and wraps it to only allow the POST method
func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "POST")
}

// allowMethod takes a HandlerFunc and wraps it in a handler that only
//
// responds if the request method is the given method, otherwise it
//
// responds with HTTP 405 Method Not Allowed.
func allowMethod(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.Header().Set("Allow", method)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

// match reports whether path matches the given pattern, which is a
//
// path with '+' wildcards wherever you want to use a parameter. Path
//
// parameters are assigned to the pointers in vars (len(vars) must be
//
// the number of wildcards), which must be of type *string or *int.
func match(path, pattern string, vars ...interface{}) bool {
	for ; pattern != "" && path != ""; pattern = pattern[1:] {
		switch pattern[0] {
		case '+':
			// '+' matches till next slash in path
			slash := strings.IndexByte(path, '/')
			if slash < 0 {
				slash = len(path)
			}
			segment := path[:slash]
			path = path[slash:]
			switch p := vars[0].(type) {
			case *string:
				*p = segment
			case *int:
				n, err := strconv.Atoi(segment)
				if err != nil || n < 0 {
					return false
				}
				*p = n
			default:
				panic("vars must be *string or *int")
			}
			vars = vars[1:]
		case path[0]:
			// non-'+' pattern byte must match path byte
			path = path[1:]
		default:
			return false
		}
	}
	return path == "" && pattern == ""
}
