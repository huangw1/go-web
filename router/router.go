package router

import (
	"net/http"
	"strings"
	"../context"
)

// https://github.com/acmacalister/helm
type Handler http.Handler

type Router struct {
	tree *node
	rootHandler Handler
}

const GET = "GET"
const POST = "POST"
const PUT = "PUT"
const DELETE = "DELETE"

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

func (r *Router) Get(path string, h Handler) {
	r.Handle(GET, path, h)
}

func (r *Router) POST(path string, h Handler) {
	r.Handle(POST, path, h)
}

func (r *Router) PUT(path string, h Handler) {
	r.Handle(PUT, path, h)
}

func (r *Router) DELETE(path string, h Handler) {
	r.Handle(DELETE, path, h)
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