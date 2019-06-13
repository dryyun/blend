package main

import (
	"fmt"
	"github.com/dryyun/blend/middleware"
	"net/http"
)

func main() {
	m := middleware.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("mux"))
	})

	m.Add(middleware.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fmt.Println("before next")
		next(rw, r)
		fmt.Println("after next")
	}))
	m.AddHandler(mux)

	http.ListenAndServe("127.0.0.1:8080", m)
}
