package main

import (
	"log"
	"net/http"
	"time"
)

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

func StaticServer(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("r.URL.Path: %s\n", r.URL.Path)
	log.Printf("[StaticServer] staticDir: %s\n", staticDir)
	h := http.StripPrefix("/static", http.FileServer(http.Dir(staticDir)))
	//h := http.FileServer(http.Dir(staticDir))
	log.Printf("handler: %+v\n", h.ServeHTTP)
	log.Printf("r.URL.Path: %s\n", r.URL.Path)
	h.ServeHTTP(w, r)
	log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, "staticServer", time.Since(start))
}
