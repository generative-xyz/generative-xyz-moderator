package middleware

import (
	"net/http"
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
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, DELETE, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, X-Requested-With, param")
		if req.Method == "OPTIONS" && req.Header.Get("Access-Control-Request-Method") != "" {
			return
		}
		next.ServeHTTP(w, req)
	})
}
