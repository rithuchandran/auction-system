package handler

import (
	"bidder/pkg/contract"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type bidHandler struct {
	service Bidder
}

func NewBidHandler(service Bidder) BidHandler {
	return &bidHandler{service: service}
}

type Bidder interface {
	Bid(ctx context.Context, request *contract.AdRequest) *contract.AdResponse
}

func (h *bidHandler) Bid(w http.ResponseWriter, r *http.Request) {
	var request contract.AdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Malformed request body: %v\n", err)
		http.Error(w, "Malformed request body", http.StatusBadRequest)
		return
	}
	adResponse := h.service.Bid(r.Context(), &request)
	err = json.NewEncoder(w).Encode(adResponse)
	if err != nil {
		log.Printf("JSON encode error: %v\n", err)
	}
}
