package middleware

import (
	"container/list"
	"net/http"
)

type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}

type Middlewares struct {
	*list.List
}

type middleware struct {
	handler Handler
}

func (m *middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(rw,r,next)
}

func New(handlers ...Handler) *Middlewares {
	l := list.New()
	for _, h := range handlers {
		l.PushBack(h)
	}
	return &Middlewares{l}
}

func (m *Middlewares) Add(handler Handler) {
	if handler == nil {
		panic("handler cannot be nil")
	}

	m.List.PushBack(handler)
}

func (m *Middlewares) AddHandler(handler http.Handler) {
	m.Add(wrap(handler))
}

func (m *Middlewares) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if m.List.Len() == 0 {
		panic("middleware's handler is empty")
	}

	// add empty Handler do nothing
	m.Add(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}))

	for h := m.List.Front(); h != nil; h = h.Next() {
		h.Value.(Handler).ServeHTTP(rw, r, nil)
	}
}

func wrap(handler http.Handler) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}
