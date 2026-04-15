package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Logger is a minimal interface satisfied by *log.Logger, *testing.T, and similar.
type Logger interface {
	Printf(format string, v ...any)
}

type FirbolgClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AppId      string
	TenancyId  string
	Logger     Logger
}

func NewFirbolgClient(baseURL string, appId int64, tenancyId int64) *FirbolgClient {
	return &FirbolgClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		AppId:      strconv.FormatInt(appId, 10),
		TenancyId:  strconv.FormatInt(tenancyId, 10),
	}
}

func (c *FirbolgClient) DoJSON(method string, urlSuffix string, body []byte, okStatuses ...int) ([]byte, error) {
	endpoint := c.BaseURL + urlSuffix
	httpReq, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-App-Id", c.AppId)
	httpReq.Header.Set("X-Tenancy-Id", c.TenancyId)
	if c.HTTPClient == nil {
		return nil, fmt.Errorf("HttpClient is nil")
	}
	if c.Logger != nil {
		c.Logger.Printf(">> %s %s body=%s", method, endpoint, string(body))
	}
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make resquest: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if c.Logger != nil {
		c.Logger.Printf("<< %d body=%s", resp.StatusCode, string(respBody))
	}

	if len(okStatuses) == 0 {
		okStatuses = []int{http.StatusOK}
	}
	for _, status := range okStatuses {
		if resp.StatusCode == status {
			return respBody, nil
		}
	}

	return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
}

func (c *FirbolgClient) PostReq(urlSuffix string, body []byte, okStatuses ...int) ([]byte, error) {
	return c.DoJSON(http.MethodPost, urlSuffix, body, okStatuses...)
}

func (c *FirbolgClient) PutReq(urlSuffix string, body []byte, okStatuses ...int) ([]byte, error) {
	return c.DoJSON(http.MethodPut, urlSuffix, body, okStatuses...)
}

func (c *FirbolgClient) PatchReq(urlSuffix string, body []byte, okStatuses ...int) ([]byte, error) {
	return c.DoJSON(http.MethodPatch, urlSuffix, body, okStatuses...)
}

func (c *FirbolgClient) DeleteReq(urlSuffix string, body []byte, okStatuses ...int) ([]byte, error) {
	return c.DoJSON(http.MethodDelete, urlSuffix, body, okStatuses...)
}

func (c *FirbolgClient) GetReq(urlSuffix string, okStatuses ...int) ([]byte, error) {
	return c.DoJSON(http.MethodGet, urlSuffix, nil, okStatuses...)
}

func MarshalJSONBody(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte{}, nil
	}
	body, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	return body, nil
}

func MustStatus(defaultStatuses ...int) []int {
	if len(defaultStatuses) == 0 {
		return []int{http.StatusOK}
	}
	return defaultStatuses
}
