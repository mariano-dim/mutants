package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Route estructura basica para encapsular los elementos basicos de Route
type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

// App es una estructura para abstraer el >Router y usarlo desde los Tests
type App struct {
	Router *mux.Router
}

// Routes Array de
type Routes []Route

// Run ejecuta el servicio
func (a *App) Run() {
	server := http.ListenAndServe(":8080", a.Router)
	log.Fatal(server)
}

// Initialize genera un nuevo Router
func (a *App) Initialize() {

	fmt.Println("Inicializando APP")
	a.Router = mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		a.Router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)
	}
}

var routes = Routes{
	Route{
		"Ping",
		"GET",
		"/ping",
		Ping,
	},
	Route{
		"MutantCheck",
		"POST",
		"/mutant",
		MutantCheck,
	},
	Route{
		"Stats",
		"GET",
		"/stats",
		StatsMutants,
	},
}
