package cbs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	clientCode string
	httpClient *http.Client
}

type FinancingData struct {
	CustomerName      string  `json:"customer_name"`
	CustomerNIK       string  `json:"customer_nik"`
	FinancingAmount   float64 `json:"financing_amount"`
	CollateralType    string  `json:"collateral_type"`
	CollateralDetails string  `json:"collateral_details"`
	Status            string  `json:"status"`
}

type SyncResponse struct {
	Synced  int `json:"synced"`
	Updated int `json:"updated"`
	Errors  int `json:"errors"`
}

func NewClient(baseURL, apiKey, clientCode string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		clientCode: clientCode,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetFinancings(ctx context.Context) ([]FinancingData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/financings", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Client-Code", c.clientCode)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cbs API returned status %d", resp.StatusCode)
	}

	var result struct {
		Data []FinancingData `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) SyncFinancings(ctx context.Context) (*SyncResponse, error) {
	financings, err := c.GetFinancings(ctx)
	if err != nil {
		return nil, err
	}

	// In a real implementation, this would sync to the local database
	// For now, we just return the count
	return &SyncResponse{
		Synced:  len(financings),
		Updated: 0,
		Errors:  0,
	}, nil
}
