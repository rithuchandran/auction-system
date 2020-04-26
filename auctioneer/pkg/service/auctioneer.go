package service

import (
	"auctioneer/pkg/contract"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Auctioneer struct {
	sync.Mutex
	Bidders map[string]string
	Client  *http.Client
}

func (svc *Auctioneer) Auction(ctx context.Context, request *contract.AdRequest) *contract.AdResponse {
	svc.Lock()
	defer svc.Unlock()

	ch := make(chan *http.Response)
	b, _ := json.Marshal(request)
	for id := range svc.Bidders {
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/bid", svc.Bidders[id]), bytes.NewBuffer(b))
		go func() {
			res, err := svc.Client.Do(request.WithContext(ctx))
			if err != nil || res.StatusCode != http.StatusOK {
				log.Printf("Error fetching price: %v", err)
				return
			}
			ch <- res
		}()
	}

	maxBidderID := ""
	maxPrice := 0.0
	count := 0
	done := make(chan struct{})

	for {
		select {
		case res := <-ch:
			bidResponse := contract.AdResponse{}
			json.NewDecoder(res.Body).Decode(&bidResponse)
			if bidResponse.Price > maxPrice {
				maxBidderID = bidResponse.BidderID
				maxPrice = bidResponse.Price
			}
			count++
			if count == len(svc.Bidders) {
				close(done)
			}
		case <-or(ctx.Done(), done): //return if all bidders have been processed or timeout is exceeded
			if maxBidderID != "" {
				return &contract.AdResponse{BidderID: maxBidderID, Price: maxPrice}
			}
			return nil
		}
	}
}

func (svc *Auctioneer) Register(ctx context.Context, request *contract.RegisterRequest) {
	svc.Lock()
	svc.Bidders[request.BidderID] = fmt.Sprintf("%s:%d", request.URL, request.Port)
	svc.Unlock()
}

func (svc *Auctioneer) List(ctx context.Context) ([]byte, error) {
	svc.Lock()
	b, err := json.Marshal(svc.Bidders)
	svc.Unlock()
	return b, err
}


func or(channels ...<-chan struct{}) <-chan struct{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan struct{})
	go func() {
		defer close(orDone)
		select {
		case <-channels[0]:
		case <-channels[1]:
		case <-or(append(channels[2:], orDone)...):
		}
	}()
	return orDone
}
