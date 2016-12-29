package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tailored-style/subscriptions-svc/handlers"
)

func buildRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.IndexHandler).
		Methods("GET")

	r.HandleFunc("/subscriptions", handlers.SubscriptionsIndexHandler)

	return r
}
