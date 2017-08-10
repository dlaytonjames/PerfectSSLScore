// PerfectSSlScore shows a webserver that gets perfect score by qualys ssllab tests.
//
package main

import (
	"crypto/tls"
	"flag"
	"github.com/gorilla/mux"
	"github.com/scusi/hsts"
	"log"
	"net/http"
)

const useTLS = true

var version = "undefined"
var listenAddr string
var certFile string
var keyFile string
var staticDir string

func init() {
	flag.StringVar(&listenAddr, "addr", ":8443", "listen address")
	flag.StringVar(&certFile, "cert", "", "TLS certificate file in PEM format")
	flag.StringVar(&keyFile, "key", "", "TLS key file (unencrypted) in PEM format")
	flag.StringVar(&staticDir, "staticDir", "", "directory with static content")
}

func main() {
	flag.Parse()
	router := NewRouter()
	// code block to get a nearly perfect ssllabs score
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(StaticWrapper(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))))
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		if useTLS == true {
			handler = hsts.NewHandler(Logger(handler, route.Name))
		} else {
			handler = Logger(handler, route.Name)
		}
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
