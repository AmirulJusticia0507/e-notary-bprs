package emeterai

import (
	"context"
	"errors"
	"time"
)

type Client struct {
	baseURL string
	apiKey  string
}

type PurchaseRequest struct {
	DocumentID string `json:"document_id"`
	Purchaser  string `json:"purchaser"`
	Amount     int64  `json:"amount"`
}

type PurchaseResponse struct {
	MeteraiSN string    `json:"meterai_sn"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{baseURL: baseURL, apiKey: apiKey}
}

func (c *Client) Purchase(_ context.Context, _ PurchaseRequest) (*PurchaseResponse, error) {
	return nil, errors.New("e-Meterai provider is not configured")
}

func (c *Client) Validate(_ context.Context, _ string) (*PurchaseResponse, error) {
	return nil, errors.New("e-Meterai provider is not configured")
}
