package main

import (
	"net/http"
	"fmt"
	"./chain"
	"./middlewares"
	"./router"
	"./context"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "首页！")
}

func list(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "列表！")
}

func detail(w http.ResponseWriter, r *http.Request) {
	reqParams := context.Get(r, "reqParams").(map[string]string)
	fmt.Fprintf(w, fmt.Sprintf("详情【%s】！", reqParams["id"]))
}

func main() {
	middles := chain.New(middlewares.RecoverMiddleware, middlewares.LoggingMiddleware)
	r := router.New(middles.ThenFunc(index))
	r.Handle("GET", "/list", middles.ThenFunc(list))
	r.Handle("GET", "/list/:id", middles.ThenFunc(detail))
	http.ListenAndServe(":8080", r)
}