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

	r.HandleFunc("/subscriptions", handlers.SubscriptionsIndexHandler).
		Methods("GET")

	r.HandleFunc("/subscriptions/{id:[0-9a-zA-Z\\-]+}", handlers.SubscriptionReadHandler).
		Methods("GET")

	r.HandleFunc("/subscriptions", handlers.SubscriptionsCreateHandler).
		Methods("POST")

	r.HandleFunc("/panic", func (w http.ResponseWriter, r *http.Request) {
		panic("THROWING A PANIC")
	}).Methods("GET")

	return r
}

