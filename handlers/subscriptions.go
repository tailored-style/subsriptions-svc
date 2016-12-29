package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tailored-style/subscriptions-svc/subscriptions"
)

const FETCH_LIMIT = 20

type subscriptionList struct {
	Subscriptions []*subscription `json:"subscriptions"`
	HasMore       bool            `json:"hasMore"`
	LastKey       *string         `json:"lastKey"`
}

type subscription struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Size  string `json:"size"`
	Email string `json:"email"`
}

func SubscriptionsIndexHandler(w http.ResponseWriter, r *http.Request) {
	result, err := subscriptions.FetchAllSubscriptions(FETCH_LIMIT, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	response := formatResponse(result)

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, string(jsonBytes))
}

func formatResponse(result *subscriptions.FetchAllSubscriptionsResult) *subscriptionList {
	if result == nil {
		panic("Result is nil!")
	}

	var lastKey *string = nil
	if result.LastEvaluatedKey != nil {
		s := *(result.LastEvaluatedKey)
		lastKey = &s
	}

	subs := make([]*subscription, len(result.Subscriptions))
	for i := 0; i < len(result.Subscriptions); i++ {
		cur := result.Subscriptions[i]
		subs[i] = &subscription{
			ID:    cur.ID,
			Name:  cur.Name,
			Size:  cur.Size,
			Email: cur.Email}
	}

	return &subscriptionList{
		HasMore:       result.HasMore,
		Subscriptions: subs,
		LastKey:       lastKey}
}
