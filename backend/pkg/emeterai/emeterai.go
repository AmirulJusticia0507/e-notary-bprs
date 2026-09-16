package emeterai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client adalah klien HTTP untuk distributor e-Meterai resmi Peruri
// (mis. Peruri Digital Security, Mitra Pajakku, Mitracomm, DLI, Pos Indonesia)
// atau mitra API-nya. Isi EMETERAI_BASE_URL + EMETERAI_API_KEY dari
// kredensial yang diberikan distributor setelah pendaftaran reseller.
type Client struct {
	baseURL string
	apiKey  string
	httpClient *http.Client
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
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsConfigured mengembalikan true bila kredensial distributor sudah diisi.
func (c *Client) IsConfigured() bool {
	return c != nil && c.baseURL != "" && c.apiKey != ""
}

func (c *Client) Purchase(ctx context.Context, req PurchaseRequest) (*PurchaseResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("e-Meterai provider is not configured (isi EMETERAI_BASE_URL dan EMETERAI_API_KEY)")
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/meterai/purchase", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("e-Meterai API returned status %d", resp.StatusCode)
	}

	var result PurchaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) Validate(ctx context.Context, meteraiSN string) (*PurchaseResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("e-Meterai provider is not configured (isi EMETERAI_BASE_URL dan EMETERAI_API_KEY)")
	}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/meterai/validate?sn="+meteraiSN, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("e-Meterai API returned status %d", resp.StatusCode)
	}

	var result PurchaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
