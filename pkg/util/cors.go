package util

import (
	"github.com/go-chi/cors"
	"net/http"
)

func AllowOriginFunc(r *http.Request, origin string) bool {
	/*if origin == "https://notes.biosave.org"  {
		return true
	}*/
	return true
}

var CorsOptions = cors.Options{
	AllowOriginFunc:  AllowOriginFunc,
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300, // Maximum value not ignored by any of major browsers
}
