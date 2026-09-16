package bpn

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
	secretKey  string
	httpClient *http.Client
}

type SertifikatData struct {
	NoSertifikat   string `json:"no_sertifikat"`
	NomorHak       string `json:"nomor_hak"`
	AlamatTanah    string `json:"alamat_tanah"`
	StatusHak      string `json:"status_hak"`
	Kelurahan      string `json:"kelurahan"`
	Kecamatan      string `json:"kecamatan"`
	Kabupaten      string `json:"kabupaten"`
	Provinsi       string `json:"provinsi"`
	LuasTanah      float64 `json:"luas_tanah"`
	NilaiJual      int64   `json:"nilai_jual"`
	StatusAktif    bool    `json:"status_aktif"`
}

type ValidasiResult struct {
	Valid          bool     `json:"valid"`
	Alasan         string   `json:"alasan"`
	Sertifikat     SertifikatData `json:"sertifikat"`
	NilaiEstimasi  int64    `json:"nilai_estimasi"`
	TanggalValidasi time.Time `json:"tanggal_validasi"`
}

func NewClient(baseURL, apiKey, secretKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		secretKey:  secretKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) ValidateSertifikat(ctx context.Context, noSertifikat string) (*ValidasiResult, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/bpn/validate?no_sertifikat="+noSertifikat, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Secret-Key", c.secretKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BPN API returned status %d", resp.StatusCode)
	}

	var result ValidasiResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CekLuasTanah(ctx context.Context, noSertifikat string) (*SertifikatData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/bpn/luas?no_sertifikat="+noSertifikat, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Secret-Key", c.secretKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BPN API returned status %d", resp.StatusCode)
	}

	var data SertifikatData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}