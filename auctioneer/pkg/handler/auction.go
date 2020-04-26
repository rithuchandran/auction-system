package handler

import (
	"auctioneer/pkg/contract"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const defaultTimeout = 200

type auctionHandler struct {
	service Auctioneer
}

func NewAuctionHandler(service Auctioneer) *auctionHandler {
	return &auctionHandler{service: service}
}

type Auctioneer interface {
	Auction(ctx context.Context, request *contract.AdRequest) *contract.AdResponse
	Register(ctx context.Context, request *contract.RegisterRequest)
	List(ctc context.Context) ([]byte, error)
}

func (h *auctionHandler) Auction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout*time.Millisecond)
	defer cancel()

	var request contract.AdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Malformed request body: %v\n", err)
		http.Error(w, "Malformed request body", http.StatusBadRequest)
		return
	}
	adResponse := h.service.Auction(ctx, &request)
	if adResponse == nil {
		log.Printf("Bidders did not respond within the deadline: %v\n", err)
		http.Error(w, "Bidders did not respond", http.StatusRequestTimeout)
		return
	}
	err = json.NewEncoder(w).Encode(adResponse)
	if err != nil {
		log.Printf("JSON encode error: %v\n", err)
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
		return
	}
}

func (h *auctionHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request contract.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Malformed request body: %v\n", err)
		http.Error(w, "Malformed request body", http.StatusBadRequest)
		return
	}
	h.service.Register(r.Context(), &request)
	w.WriteHeader(http.StatusOK)
}

func (h *auctionHandler) List(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.List(r.Context())
	if err != nil {
		log.Printf("Error fetching list: %v\n", err)
		http.Error(w, "Error fetching list", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Error writing response: %v\n", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}
