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

// variables used to store commandline flags
var listenAddr string
var certFile string
var keyFile string
var staticDir string

// init initializes flags and prints the version and buildtime
func init() {
	flag.StringVar(&listenAddr, "addr", ":443", "listen address (required)")
	flag.StringVar(&certFile, "cert", "", "TLS certificate file in PEM format (required)")
	flag.StringVar(&keyFile, "key", "", "TLS key file (unencrypted) in PEM format (required)")
	flag.StringVar(&staticDir, "staticDir", "", "directory with static content (optional)")
	fmt.Fprintf(os.Stderr, "compiled from git commit: %s, Buildtime: %s\n", version, buildtime)
}

func main() {
	flag.Parse()
	// check required flags
	if certFile == "" {
		err := fmt.Errorf("No certificate file set, use -cert flag.")
		log.Fatal(err)
	}
	if keyFile == "" {
		err := fmt.Errorf("No key file set, use -key flag.")
		log.Fatal(err)
	}
	router := NewRouter()
	// create a TLS config
	cfg := &tls.Config{
		// only allow TLS version 1.2
		MinVersion: tls.VersionTLS12,
		// only use secure curves
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		// prefer chipher suites of the server
		PreferServerCipherSuites: true,
		// define server cipher suites to be used
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	// define a http.Server that uses the above TLS config and our router
	srv := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	// start the http server
	log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}

// NewRouter creates a new mux router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// if a staticDir was set, use it.
	if staticDir != "" {
		// start a http.Fileserver for the given staticDir and log it's requests (via StaticWrapper)
		router.PathPrefix("/static/").Handler(StaticWrapper(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))))
	}
	// add all routes from Routes.go
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
