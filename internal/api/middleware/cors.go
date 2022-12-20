package middleware

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// AllowCORS
var AllowCORS = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		corsAllowOrigin := "*"
		if origin := req.Header.Get("Origin"); origin != "" {
			corsAllowOrigin = origin
		}
		w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")


		spew.Dump(req.Header.Get("Access-Control-Request-Method"))
		if req.Method == "OPTIONS" && req.Header.Get("Access-Control-Request-Method") != "" {
			spew.Dump("Run here")
			return
		}
		next.ServeHTTP(w, req)
	})
}
