package bullet

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type FirbolgClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AppId      string
}

func NewFirbolgClient(baseURL string, appId int) *FirbolgClient {
	return &FirbolgClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		AppId:      strconv.Itoa(appId), //VX:Note app id is part of all clients here. We could make it optional if needed.
	}
}

func (c *FirbolgClient) PostReq(urlSuffix string, body []byte) ([]byte, error) {
	httpReq, err := http.NewRequest("POST", c.BaseURL+urlSuffix, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-App-ID", c.AppId)
	if c.HTTPClient == nil {
		return nil, fmt.Errorf("HttpClient is nil")
	}
	// execute
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make resquest: %w", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}
