package bullet

import (
	"encoding/json"
	"fmt"

	util "github.com/vixac/firbolg_clients/util"
)

type WayFinderClientInterface interface {
	WayFinderInsertOne(req WayFinderPutRequest) (int64, error)
	WayFinderQueryByPrefix(req WayFinderPrefixQueryRequest) ([]WayFinderQueryItem, error)
	WayFinderGetOne(req WayFinderGetOneRequest) (*WayFinderItem, error)
}

type WayFinderClient struct {
	Client *util.FirbolgClient
}

func NewWayFinderClient(baseURL string, appId int64) *WayFinderClient {
	firbolg := util.NewFirbolgClient(baseURL, appId)
	return &WayFinderClient{
		Client: firbolg,
	}
}

func (c *WayFinderClient) WayFinderInsertOne(req WayFinderPutRequest) (int64, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	if c.Client == nil {
		return 0, fmt.Errorf("FirbolgClient is nil")
	}
	resp, err := c.Client.PostReq("/wayfinder/insert-one", bodyBytes)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	var result struct {
		ItemId int64 `json:"itemId"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return result.ItemId, nil
}

// WayFinderGetOne sends a GET request to retrieve a WayFinder item by bucket/key
func (c *WayFinderClient) WayFinderGetOne(req WayFinderGetOneRequest) (*WayFinderItem, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if c.Client == nil {
		return nil, fmt.Errorf("FirbolgClient is nil")
	}
	fmt.Println("FC: about to send this request to get one", string(bodyBytes))
	resp, err := c.Client.PostReq("/wayfinder/get-one", bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var result struct {
		Item WayFinderItem `json:"item"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result.Item, nil
}

func (c *WayFinderClient) WayFinderQueryByPrefix(req WayFinderPrefixQueryRequest) ([]WayFinderQueryItem, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if c.Client == nil {
		return nil, fmt.Errorf("FirbolgClient is nil")
	}
	resp, err := c.Client.PostReq("/wayfinder/query-by-prefix", bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var result struct {
		Items []WayFinderQueryItem `json:"items"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return result.Items, nil
}
