package main

import (
	"log"
	"net/http"
	"time"
)

// StaticWrapper does logging as Logger.go would do for http.FileServer
func StaticWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r) // call original
		log.Printf("%s\t%s\t%s\t%s\t%s",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			"staticWrapper",
			time.Since(start))
	})
}
