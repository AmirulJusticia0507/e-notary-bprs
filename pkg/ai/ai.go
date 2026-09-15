package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Config menyimpan konfigurasi klien AI.
type Config struct {
	BaseURL string
	APIKey  string
	Model   string
}

// NewConfig membuat konfigurasi dari environment variable.
func NewConfig() Config {
	return Config{
		BaseURL: os.Getenv("AI_BASE_URL"),
		APIKey:  os.Getenv("AI_API_KEY"),
		Model:   os.Getenv("AI_MODEL"),
	}
}

// Client klien HTTP untuk komunikasi ke API AI.
type Client struct {
	cfg Config
}

// Chat mengirim prompt ke model AI dan mengembalikan respons teks.
func (c *Client) Chat(prompt string) (string, error) {
	url := c.cfg.BaseURL + "/chat/completions"
	reqBody := map[string]interface{}{
		"model":  c.cfg.Model,
		"messages": []map[string]interface{}{
			{"role": "user", "content": prompt},
		},
	}
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("gagal marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("gagal buat http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal kirim request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AI response status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gagal baca response body: %w", err)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("gagal parse JSON response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("tidak ada pilihan dari AI")
	}

	return result.Choices[0].Message.Content, nil
}