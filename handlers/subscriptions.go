package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tailored-style/subscriptions-svc/subscriptions"
)

const FETCH_LIMIT = 20

type outputSubscriptionList struct {
	Subscriptions []*outputSubscription `json:"subscriptions"`
	HasMore       bool                  `json:"hasMore"`
	LastKey       *string               `json:"lastKey"`
}

type outputSubscription struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        string `json:"size"`
	Email       string `json:"email"`
	StripeToken string `json:"stripeToken"`
}

func SubscriptionsIndexHandler(w http.ResponseWriter, r *http.Request) {
	result, err := subscriptions.FetchAllSubscriptions(FETCH_LIMIT, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := formatIndexResponse(result)

	writeJsonResponse(w, response)
}

func SubscriptionReadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := subscriptions.FetchSubscription(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := &outputSubscription{
		ID:          res.ID,
		Name:        res.Name,
		Email:       res.Email,
		Size:        res.Size,
		StripeToken: res.StripeToken,
	}

	writeJsonResponse(w, response)
}

func SubscriptionsCreateHandler(w http.ResponseWriter, r *http.Request) {
	type reqFormat struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Size  string `json:"size"`
	}

	// Parse the inputs
	var input *reqFormat = &reqFormat{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if input.Name == "" {
		http.Error(w, "Name must be provided", http.StatusBadRequest)
		return
	}

	if input.Email == "" {
		http.Error(w, "Email must be provided", http.StatusBadRequest)
		return
	}

	if input.Size == "" {
		http.Error(w, "Size must be provided", http.StatusBadRequest)
		return
	}

	result, err := subscriptions.CreateSubscription(input.Name, input.Email, input.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := formatCreateResponse(result)

	writeJsonResponse(w, response)
}

func writeJsonResponse(w http.ResponseWriter, marshallable interface{}) {
	jsonBytes, err := json.Marshal(marshallable)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(jsonBytes))
}

func formatCreateResponse(sub *subscriptions.Subscription) *outputSubscription {
	if sub == nil {
		panic("Subscription is nil!")
	}

	return &outputSubscription{
		ID:          sub.ID,
		Name:        sub.Name,
		Email:       sub.Email,
		Size:        sub.Size,
		StripeToken: sub.StripeToken,
	}
}

func formatIndexResponse(result *subscriptions.FetchAllSubscriptionsResult) *outputSubscriptionList {
	if result == nil {
		panic("Result is nil!")
	}

	var lastKey *string = nil
	if result.LastEvaluatedKey != nil {
		s := *(result.LastEvaluatedKey)
		lastKey = &s
	}

	subs := make([]*outputSubscription, len(result.Subscriptions))
	for i := 0; i < len(result.Subscriptions); i++ {
		cur := result.Subscriptions[i]
		subs[i] = &outputSubscription{
			ID:          cur.ID,
			Name:        cur.Name,
			Size:        cur.Size,
			Email:       cur.Email,
			StripeToken: cur.StripeToken,
		}
	}

	return &outputSubscriptionList{
		HasMore:       result.HasMore,
		Subscriptions: subs,
		LastKey:       lastKey,
	}
}
