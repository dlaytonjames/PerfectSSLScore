package main

import (
	"net/http"
	//"github.com/tsuru/config"
)

func init() {
	//StaticDir, err := config.GetString("AdminUI:Dir:Static")
	//checkFatal(err)
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
