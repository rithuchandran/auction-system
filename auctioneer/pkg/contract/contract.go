package contract

type AdRequest struct {
	AuctionID string `json:"auction_id"`
}

type AdResponse struct {
	BidderID string  `json:"bidder_id"`
	Price    float64 `json:"price"`
}

type RegisterRequest struct {
	BidderID string `json:"bidder_id"`
	URL      string `json:"url"`
	Port     int    `json:"port"`
}