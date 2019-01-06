package middlewares

import (
	"net/http"
	"github.com/huangw1/go-web/util"
)

func Static(directories ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if len(directories) == 0 {
			directories = append(directories, "public")
		}
		n := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				next.ServeHTTP(w, r)
			} else {
				for _, dir := range directories {
					if util.IsDir(dir + r.URL.Path) {
						continue
					}
					if util.IsFile(dir + r.URL.Path) {
						http.ServeFile(w, r, dir + r.URL.Path)
						return
					}
				}

			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(n)
	}
}
