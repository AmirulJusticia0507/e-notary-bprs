package esign

import (
	"context"
	"errors"
	"time"
)

type Client struct {
	baseURL string
	apiKey  string
}

type SignRequest struct {
	DocumentID  string `json:"document_id"`
	SignerName  string `json:"signer_name"`
	SignerEmail string `json:"signer_email"`
}

type SignResponse struct {
	SignID   string     `json:"sign_id"`
	Status   string     `json:"status"`
	SignedAt *time.Time `json:"signed_at,omitempty"`
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{baseURL: baseURL, apiKey: apiKey}
}

func (c *Client) Sign(_ context.Context, _ SignRequest) (*SignResponse, error) {
	return nil, errors.New("e-Sign provider is not configured")
}

func (c *Client) GetStatus(_ context.Context, _ string) (*SignResponse, error) {
	return nil, errors.New("e-Sign provider is not configured")
}
