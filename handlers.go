package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func tenSecHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(10 * time.Second):
		w.WriteHeader(http.StatusOK)
		return
	case <-r.Context().Done():
		if r.Context().Err() == context.DeadlineExceeded {
			w.WriteHeader(http.StatusRequestTimeout)
			return
		} else if r.Context().Err() == context.Canceled {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}
	}
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	// parse input
	type requestParams struct {
		Items []string `json:"items"`
	}
	decoder := json.NewDecoder(r.Body)
	params := requestParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Print("Error decoding for /items")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// convert to "items" array
	items := params.Items

	// generate response
	type response struct {
		Items []string `json:"items"`
	}

	responseBody := response{
		Items: items,
	}
	data, err := json.Marshal(responseBody)
	if err != nil {
		log.Print("Unable to marshal items to JSON.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
