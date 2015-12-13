package starfighter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// AuthHeader contains the Starfighter Authorization Header
	AuthHeader = "X-Starfighter-Authorization"
	// APILocation sets the Starfighter API Location
	APILocation = "https://api.stockfighter.io/ob/api"
)

// Client reflects a HTTP REST client to the Starfighter API.
type Client struct {
	// Your Starfighter API Token
	Token string
	// Location of the API
	Location string
	// The HTTP Client to use
	Client http.Client
}

// CallReq sets the authorization header and runs the request
func (c *Client) CallReq(req *http.Request) (*http.Response, error) {
	req.Header.Add(AuthHeader, c.Token)
	return c.Client.Do(req)
}

// Call hits a method, endpoint (without the location), with specified data (if necessary).
// It then returns the JSON response (with or without an error if necessary).
// If an error is returned and it is of type APIError, then the API has barfed on you.
// If it is not of type APIError, then your client has barfed on you.
func (c *Client) Call(method, endpoint string, data interface{}) (map[string]interface{}, *bytes.Buffer, error) {
	// set up the request
	req, err := http.NewRequest(method, c.Location+endpoint, nil)
	if data == nil {
		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)
		if err = encoder.Encode(data); err != nil {
			return nil, nil, err
		}

		req, err = http.NewRequest(method, c.Location+endpoint, buf)
	}

	if err != nil {
		return nil, nil, err
	}

	// call the request
	resp, err := c.CallReq(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	// keep a copy in case other methods do strange things
	copy := &bytes.Buffer{}
	reader := io.TeeReader(resp.Body, copy)

	// unmarshal
	body := map[string]interface{}{}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(body)
	if err != nil {
		return nil, copy, err
	}

	// and let's check for errors as a precaution
	var apiErr error
	if body["ok"] != false {
		apiErr = &APIError{
			Code:    resp.StatusCode,
			Message: body["error"].(string),
		}
	}

	return body, copy, apiErr
}

// Heartbeat checks if the API is up. Because maybe it isn't.
func (c *Client) Heartbeat() bool {
	_, _, err := c.Call("GET", "/heartbeat", nil)
	return err == nil
}

// VenueHealthCheck checks if a venue is up.
func (c *Client) VenueHealthCheck(venue string) bool {
	_, _, err := c.Call("GET", fmt.Sprintf("/venues/%s/heartbeat", venue), nil)
	return err == nil
}

// ListVenueStocks lists the stocks in a venue
func (c *Client) ListVenueStocks(venue string) ([]Stock, error) {
	resp, _, err := c.Call("GET", fmt.Sprintf("/venues/%s/stocks", venue), nil)
	if err != nil {
		return nil, err
	}

	stocks := make([]Stock, len(resp["symbols"].([]interface{})))

	for k, v := range resp["symbols"].([]map[string]interface{}) {
		stocks[k] = Stock{
			Name:   v["name"].(string),
			Symbol: v["symbol"].(string),
		}
	}

	return stocks, nil
}

// GetStockOrderbook retrieves the orderbook for the stock requested.
func (c *Client) GetStockOrderbook(venue, stock string) (*OrderBook, error) {
	_, copy, err := c.Call("GET", fmt.Sprintf("/venues/%s/stocks/%s", venue, stock), nil)
	if err != nil {
		return nil, err
	}

	orderBook := OrderBook{}

	decoder := json.NewDecoder(copy)
	err = decoder.Decode(&orderBook)

	return &orderBook, err
}

// PlaceStockOrder places an order for a stock.
func (c *Client) PlaceStockOrder(account, venue, stock string, price int64, qty int64, direction, ordertype string) (*OrderResult, error) {
	_, copy, err := c.Call("POST", fmt.Sprintf("/venues/%s/stocks/%s/orders", venue, stock), map[string]interface{}{
		"account":   account,
		"venue":     venue,
		"stock":     stock,
		"price":     price,
		"qty":       qty,
		"direction": direction,
		"orderType": ordertype,
	})

	if err != nil {
		return nil, err
	}

	orderResult := OrderResult{}

	decoder := json.NewDecoder(copy)
	err = decoder.Decode(&orderResult)

	return &orderResult, err
}

// QuoteStock shows you the most recent information. Which is probably outdated
// by the time you actually interpret it. So why are you even doing this?
func (c *Client) QuoteStock(venue, stock string) (*StockQuote, error) {
	_, copy, err := c.Call("GET", fmt.Sprintf("/venues/%s/stocks/%s/quote", venue, stock), nil)
	if err != nil {
		return nil, err
	}

	stockQuote := StockQuote{}

	decoder := json.NewDecoder(copy)
	err = decoder.Decode(&stockQuote)

	return &stockQuote, err
}

// GetOrderStatus retrieves the status for an existing order. Slowly.
func (c *Client) GetOrderStatus(venue, stock string, order int64) (*OrderResultAlt, error) {
	_, copy, err := c.Call("GET", fmt.Sprintf("/venues/%s/stocks/%s/orders/%d", venue, stock, order), nil)
	if err != nil {
		return nil, err
	}

	orderResult := OrderResultAlt{}

	decoder := json.NewDecoder(copy)
	err = decoder.Decode(&orderResult)

	return &orderResult, err
}

// CancelOrder attempts to cancel the order. Good luck, though.
func (c *Client) CancelOrder(venue, stock string, order int64) (*OrderResultAlt, error) {
	_, copy, err := c.Call("DELETE", fmt.Sprintf("/venues/%s/stocks/%s/orders/%d", venue, stock, order), nil)
	if err != nil {
		return nil, err
	}

	orderResult := OrderResultAlt{}

	decoder := json.NewDecoder(copy)
	err = decoder.Decode(&orderResult)

	return &orderResult, err
}

// ListVenueOrderStatus lists the status of all orders for the venue and account.
func (c *Client) ListVenueOrderStatus(venue, account string) (*OrderResultList, error) {
	_, copy, err := c.Call("GET", fmt.Sprintf("/venues/%s/accounts/%s/orders", venue, account), nil)
	if err != nil {
		return nil, err
	}

	orderResultList := OrderResultList{}

	decoder := json.NewDecoder(copy)
	err = decoder.Decode(&orderResultList)

	return &orderResultList, err
}

// ListVenueStockOrderStatus lists the status of all orders for the venue, stock, and account.
func (c *Client) ListVenueStockOrderStatus(venue, stock, account string) (*OrderResultList, error) {
	_, copy, err := c.Call("GET", fmt.Sprintf("/venues/%s/accounts/%s/stocks/%s/orders", venue, stock, account), nil)
	if err != nil {
		return nil, err
	}

	orderResultList := OrderResultList{}

	decoder := json.NewDecoder(copy)
	err = decoder.Decode(&orderResultList)

	return &orderResultList, err
}
