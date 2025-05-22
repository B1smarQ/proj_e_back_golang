package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type authClient struct {
	baseURL    string
	httpClient *http.Client
}

type authResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func NewAuthClient(baseURL string) *authClient {
	return &authClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *authClient) ValidateToken(token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/validate", c.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp authResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		return nil, fmt.Errorf("auth service error: %s", errResp.Message)
	}

	var authResp authResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if authResp.Status != "success" {
		return nil, fmt.Errorf("auth service error: %s", authResp.Message)
	}

	return authResp.Data, nil
}
