package router

import (
	"net/http"
	"strings"
	"../context"
)

type Handler http.Handler

type Router struct {
	tree *node
	rootHandler Handler
}

func New(h Handler) *Router {
	return &Router{
		&node{
			component: "/",
			methods: make(map[string]Handler),
		},
		h,
	}
}

func (r *Router) Handle(method, path string, h Handler) {
	if path[0] != '/' {
		panic("Path must starts with a /.")
	}
	r.tree.addNode(method, path, h)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	params := make(map[string]string)
	context.Set(req, "reqParams", params)
	node, _, _ := r.tree.traverse(strings.Split(req.URL.Path, "/")[1:], params)
	if handler := node.methods[req.Method]; handler != nil {
		handler.ServeHTTP(w, req)
	} else {
		r.rootHandler.ServeHTTP(w, req)
	}
}