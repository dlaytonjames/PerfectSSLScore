package main

import (
	"net/http"
)

func init() {
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"Home", "GET", "/", HomeHandler},
}
