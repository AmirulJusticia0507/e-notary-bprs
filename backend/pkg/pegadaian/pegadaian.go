package pegadaian

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client adalah klien HTTP untuk integrasi Pegadaian (Tahap 5).
// Dipakai untuk validasi agunan gadai (emas/kendaraan/elektronik)
// sebelum order legal diteruskan ke notaris.
type Client struct {
	baseURL    string
	apiKey     string
	partnerCode string
	httpClient *http.Client
}

type GadaiData struct {
	NoSBG         string `json:"no_sbg"`
	NamaNasabah   string `json:"nama_nasabah"`
	NIK           string `json:"nik"`
	JenisBarang   string `json:"jenis_barang"`
	Taksiran      int64  `json:"taksiran"`
	Pinjaman      int64  `json:"pinjaman"`
	StatusGadai   string `json:"status_gadai"`
	TanggalJatuhTempo string `json:"tanggal_jatuh_tempo"`
	Cabang        string `json:"cabang"`
}

type ValidasiResult struct {
	Valid           bool      `json:"valid"`
	Alasan          string    `json:"alasan"`
	Gadai           GadaiData `json:"gadai"`
	TanggalValidasi time.Time `json:"tanggal_validasi"`
}

func NewClient(baseURL, apiKey, partnerCode string) *Client {
	return &Client{
		baseURL:     baseURL,
		apiKey:      apiKey,
		partnerCode: partnerCode,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) ValidateGadai(ctx context.Context, noSBG string) (*ValidasiResult, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/gadai/validate?no_sbg="+noSBG, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Partner-Code", c.partnerCode)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pegadaian API returned status %d", resp.StatusCode)
	}

	var result ValidasiResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
