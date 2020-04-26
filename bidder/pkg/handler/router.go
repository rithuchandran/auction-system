package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	BidPath = "/bid"
)

type BidHandler interface {
	Bid(w http.ResponseWriter, r *http.Request)
}

func New(bidder BidHandler) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc(BidPath, bidder.Bid).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 page not found", http.StatusNotFound)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	})

	return router
}
