// go run middleware.go
// curl localhost:8080
// curl localhost:8080/panic

package main

import (
	"fmt"
	"github.com/dryyun/blend/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("mux"))
	})
	mux.HandleFunc("/panic", func(rw http.ResponseWriter, r *http.Request) {
		panic("trigger panic")
	})

	m := middleware.New()

	// recover middleware
	m.Add(middleware.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("recover %v", err)
			}
		}()
		fmt.Println("recover middleware before next")
		next(rw, r)
		fmt.Println("recover middleware after next")
	}))

	// log middleware
	m.Add(middleware.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var logger = log.New(os.Stderr, "[xxx]", log.LstdFlags|log.Lshortfile)
		startTime := time.Now()
		logger.Printf("log middleware receive request %s \n", r.Method)
		next(rw, r)
		logger.Printf("log middleware end request request_time %s \n", time.Since(startTime).String())
	}))

	// normal middleware
	m.Add(middleware.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fmt.Println("normal middleware before next")
		next(rw, r)
		fmt.Println("normal middleware after next")
	}))

	m.AddHandler(mux)

	if err := http.ListenAndServe("127.0.0.1:8080", m); err != nil {
		log.Fatal(err)
	}
}
