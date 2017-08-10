// PerfectSSlScore shows a webserver that gets perfect score by qualys ssllab tests.
//
// You need to have a RSA4096 server key/cert in order to get 100% in the key exchange tests.
// With RSA2048 you only get 90% in the key exchange category.
//
// build with the following arguments in order to maintain version and buildtime.
// go build -i -v -ldflags="-s -w -X main.version=$(git describe --always --long) -X 'main.buildtime=$(date -u '+%Y-%m-%d %H:%M:%S')'" ./
//
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/scusi/hsts"
	"log"
	"net/http"
	"os"
)

// version and buildtime must be set by the compiler
// use the following build ldflags to do so:
// -X main.version=$(git describe --always --long) -X 'main.buildtime=$(date -u '+%Y-%m-%d %H:%M:%S')'
var version = "undefined"
var buildtime = "undefined"

var listenAddr string
var certFile string
var keyFile string
var staticDir string

func init() {
	flag.StringVar(&listenAddr, "addr", ":443", "listen address (required)")
	flag.StringVar(&certFile, "cert", "", "TLS certificate file in PEM format (required)")
	flag.StringVar(&keyFile, "key", "", "TLS key file (unencrypted) in PEM format (required)")
	flag.StringVar(&staticDir, "staticDir", "", "directory with static content (optional)")
	fmt.Fprintf(os.Stderr, "compiled from git commit: %s, Buildtime: %s\n", version, buildtime)
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
	// if a staticDir was set use it.
	if staticDir != "" {
		router.PathPrefix("/static/").Handler(StaticWrapper(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))))
	}
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = hsts.NewHandler(Logger(handler, route.Name))
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
