package iex

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultBaseURL = "https://cloud.iexapis.com/beta"

type Client struct {
	HTTPClient *http.Client
	Token      string
	BaseURL    string
}

//TODO: Prices should not be floats.  We should deserialize to a custom price type
type Quote struct {
	Symbol           string  `json:"symbol"`
	CompanyName      string  `json:"companyName"`
	CalculationPrice string  `json:"calculationPrice"`
	Open             float64 `json:"open"`
	OpenTime         int64   `json:'openTime"`
	Close            float64 `json:"close"`
	CloseTime        int64   `json:"closeTime"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	LatestPrice      float64 `json:"latestPrice"`
	LatestSource     string  `json:"latestSource"`
	LatestTime       string  `json:"latestTime"`
	LatestUpdate     int64   `json:"latestUpdate"`
	LatestVolume     int64   `json:"latestVolume"`
	Change           float64 `json:"change"`
	ChangePercent    float64 `json:"changePercent"`
}

// Do is a proxy for the HTTPClient Do but sets the query param for the token.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	queryParams := req.URL.Query()
	queryParams.Set("token", c.Token)
	req.URL.RawQuery = queryParams.Encode()

	return c.HTTPClient.Do(req)
}

func (c *Client) GetQuote(symbol string) (*Quote, error) {
	// Injection attack?
	url := c.BaseURL + "/stock/" + symbol + "/quote"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %s", res.Status)
	}

	dec := json.NewDecoder(res.Body)
	var quote Quote
	err = dec.Decode(&quote)
	if err != nil {
		return nil, nil
	}

	return &quote, nil
}
