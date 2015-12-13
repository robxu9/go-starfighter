package starfighter

import "time"

// Stock represents a symbol on the venue.
type Stock struct {
	Name   string
	Symbol string
}

// StockQuote shows a quote for a stock.
type StockQuote struct {
	Symbol    string    `json:"symbol"`
	Venue     string    `json:"venue"`
	Bid       int       `json:"bid"`
	Ask       int       `json:"ask"`
	BidSize   int       `json:"bidSize"`
	AskSize   int       `json:"askSize"`
	BidDepth  int       `json:"bidDepth"`
	AskDepth  int       `json:"askDepth"`
	Last      int       `json:"last"`
	LastSize  int       `json:"lastSize"`
	LastTrade time.Time `json:"lastTrade"`
	QuoteAt   time.Time `json:"quoteTime"`
}

// OrderBook represents the current state of an order.
type OrderBook struct {
	Asks []struct {
		IsBuy bool `json:"isBuy"`
		Price int  `json:"price"`
		Qty   int  `json:"qty"`
	} `json:"asks"`
	Bids []struct {
		IsBuy bool `json:"isBuy"`
		Price int  `json:"price"`
		Qty   int  `json:"qty"`
	} `json:"bids"`
	Symbol    string    `json:"symbol"`
	Timestamp time.Time `json:"ts"`
	Venue     string    `json:"venue"`
}

// OrderResult details the result of an order.
type OrderResult struct {
	Symbol      string    `json:"symbol"`
	Venue       string    `json:"venue"`
	Direction   string    `json:"direction"`
	OriginalQty int       `json:"originalQty"`
	Qty         int       `json:"qty"`
	Price       int       `json:"price"`
	Type        string    `json:"type"`
	ID          int       `json:"id"`
	Account     string    `json:"account"`
	Timestamp   time.Time `json:"ts"`
	Fills       []struct {
		Price     int       `json:"price"`
		Qty       int       `json:"qty"`
		Timestamp time.Time `json:"ts"`
	} `json:"fills"`
	TotalFilled int  `json:"totalFilled"`
	Open        bool `json:"open"`
}

// OrderResultAlt is the same thing as OrderResult, except uses orderType instead of type.
// goddamnit, api.
type OrderResultAlt struct {
	Symbol      string    `json:"symbol"`
	Venue       string    `json:"venue"`
	Direction   string    `json:"direction"`
	OriginalQty int       `json:"originalQty"`
	Qty         int       `json:"qty"`
	Price       int       `json:"price"`
	Type        string    `json:"orderType"`
	ID          int       `json:"id"`
	Account     string    `json:"account"`
	Timestamp   time.Time `json:"ts"`
	Fills       []struct {
		Price     int       `json:"price"`
		Qty       int       `json:"qty"`
		Timestamp time.Time `json:"ts"`
	} `json:"fills"`
	TotalFilled int  `json:"totalFilled"`
	Open        bool `json:"open"`
}

// OrderResultList shows a list of orders.
type OrderResultList struct {
	Orders []OrderResultAlt `json:"orders"`
}
