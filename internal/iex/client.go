package iex

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

const defaultBaseURL = "https://cloud.iexapis.com/v1"

type Client struct {
	HTTPClient *http.Client
	Token      string
	BaseURL    string
}

type Quote struct {
	Symbol           string          `json:"symbol"`
	CompanyName      string          `json:"companyName"`
	CalculationPrice string          `json:"calculationPrice"`
	Open             decimal.Decimal `json:"open"`
	OpenTime         int64           `json:'openTime"`
	Close            decimal.Decimal `json:"close"`
	CloseTime        int64           `json:"closeTime"`
	High             decimal.Decimal `json:"high"`
	Low              decimal.Decimal `json:"low"`
	LatestPrice      decimal.Decimal `json:"latestPrice"`
	LatestSource     string          `json:"latestSource"`
	LatestTime       string          `json:"latestTime"`
	LatestUpdate     int64           `json:"latestUpdate"`
	LatestVolume     int64           `json:"latestVolume"`
	Change           decimal.Decimal `json:"change"`
	ChangePercent    decimal.Decimal `json:"changePercent"`
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
