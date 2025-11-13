package server

import (
	"fmt"
	"log"
	"net/http"
)

func PanicRecoveryMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println(fmt.Sprint(rec))

				JSONInternalServerError(w)
			}
		}()

		h.ServeHTTP(w, r)
	}
}
