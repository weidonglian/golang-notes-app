package middleware

import (
	"github.com/go-chi/cors"
	"net/http"
)

func allowOriginFunc(r *http.Request, origin string) bool {
	return CheckOriginFunc(r)
}

func CheckOriginFunc(r *http.Request) bool {
	return true
}

func Cors() func(handler http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowOriginFunc:  allowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
