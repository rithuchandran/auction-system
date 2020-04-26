package service

import (
	"bidder/pkg/contract"
	"context"
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"time"
)

type Bidder struct {
	Delay time.Duration
	ID    string
}

func (svc *Bidder) Bid(ctx context.Context, request *contract.AdRequest) *contract.AdResponse {
	rand.Seed(hash(fmt.Sprintf("%s:%s", svc.ID, request.AuctionID)))
	response := &contract.AdResponse{
		BidderID: svc.ID,
		Price:    math.Round(float64(rand.Int31())+rand.Float64()*100) / 100,
	}
	time.Sleep(svc.Delay)
	return response
}

func hash(s string) int64 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int64(h.Sum32())
}
