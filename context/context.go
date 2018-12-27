package context

import (
	"sync"
	"net/http"
	"time"
)

var mutex sync.RWMutex
var data = make(map[*http.Request]map[interface{}]interface{})
var ts = make(map[*http.Request]int64)

func Set(r *http.Request, key, val interface{})  {
	mutex.Lock()
	if data[r] == nil {
		data[r] = make(map[interface{}]interface{})
		ts[r] = time.Now().Unix()
	}
	data[r][key] = val
	mutex.Unlock()
}

func Get(r *http.Request, key interface{}) interface{} {
	mutex.RLock()
	if ctx := data[r]; ctx != nil {
		val := ctx[key]
		mutex.RUnlock()
		return val
	}
	mutex.RUnlock()
	return nil
}

func Delete(r *http.Request, key interface{})  {
	mutex.Lock()
	if ctx := data[r]; ctx != nil {
		delete(ctx, key)
	}
	mutex.Unlock()
}

func Clear(r *http.Request)  {
	mutex.Lock()
	delete(data, r)
	delete(ts, r)
	mutex.Unlock()
}

func ClearHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			Clear(r)
		}()
		h.ServeHTTP(w, r)
	})
}

