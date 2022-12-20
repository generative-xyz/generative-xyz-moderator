package middleware

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// AllowCORS
var AllowCORS = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// corsAllowOrigin := "*"
		// if origin := req.Header.Get("Origin"); origin != "" {
		// 	corsAllowOrigin = origin
		// }

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")


		spew.Dump(req.Header.Get("Access-Control-Request-Method"))
		spew.Dump(req.Method )
		if req.Method == "OPTIONS" && req.Header.Get("Access-Control-Request-Method") != "" {
			
			return
		}
		next.ServeHTTP(w, req)
	})
}
