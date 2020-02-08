package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct{}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func (s Server) Start() error {
	var routes = []route{
		{
			"Healthcheck",
			"GET",
			"/healthcheck",
			s.HandleHealthcheck,
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(route.HandlerFunc)
	}

	return http.ListenAndServe(":3000", router)
}
