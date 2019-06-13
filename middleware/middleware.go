package middleware

import (
	"container/list"
	"net/http"
)

var l *list.List

type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}

type middleware struct {
	handler Handler
	listElm *list.Element
}

func (m *middleware) serve(rw http.ResponseWriter, r *http.Request) {
	next := m.listElm.Next()
	if next != nil {
		m.handler.ServeHTTP(rw, r, next.Value.(*middleware).serve)
	}
}

func New(handlers ...Handler) *middleware {
	l = list.New()
	for _, h := range handlers {
		m := &middleware{h, nil}
		m.listElm = l.PushBack(m)
	}
	return new(middleware)
}

func (m *middleware) Add(handler Handler) {
	if handler == nil {
		panic("handler cannot be nil")
	}
	add := &middleware{handler, nil}
	add.listElm = l.PushBack(add)
}

func (m *middleware) AddFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	m.Add(HandlerFunc(handlerFunc))
}

func (m *middleware) AddHandler(handler http.Handler) {
	m.Add(Wrap(handler))
}

func (m *middleware) AddHandlerFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request)) {
	m.AddHandler(http.HandlerFunc(handlerFunc))
}

func (m *middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if l.Len() == 0 {
		panic("middleware's handler is empty")
	}
	// add empty Handler do nothing
	m.Add(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}))

	l.Front().Value.(*middleware).serve(rw, r)
}

func Wrap(handler http.Handler) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}
