package esign

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client adalah klien HTTP untuk PSrE (Penyelenggara Sertifikasi Elektronik)
// yang diakui Komdigi — daftar resmi di https://tte.komdigi.go.id.
// Isi ESIGN_BASE_URL + ESIGN_API_KEY dari kredensial yang diberikan
// PSrE setelah PKS. Hanya TTE tersertifikasi yang sah penuh sebagai
// alat bukti (UU ITE + PP PSTE).
type Client struct {
	baseURL string
	apiKey  string
	httpClient *http.Client
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
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsConfigured mengembalikan true bila kredensial PSrE sudah diisi.
func (c *Client) IsConfigured() bool {
	return c != nil && c.baseURL != "" && c.apiKey != ""
}

func (c *Client) Sign(ctx context.Context, req SignRequest) (*SignResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("e-Sign provider is not configured (isi ESIGN_BASE_URL dan ESIGN_API_KEY)")
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/sign/request", bytes.NewReader(body))
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
		return nil, fmt.Errorf("e-Sign API returned status %d", resp.StatusCode)
	}

	var result SignResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetStatus(ctx context.Context, signID string) (*SignResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("e-Sign provider is not configured (isi ESIGN_BASE_URL dan ESIGN_API_KEY)")
	}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/sign/status?sign_id="+signID, nil)
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
		return nil, fmt.Errorf("e-Sign API returned status %d", resp.StatusCode)
	}

	var result SignResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
