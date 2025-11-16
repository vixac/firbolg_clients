package rest_bullet

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"github.com/vixac/firbolg_clients/bullet/util"
)

type TrackClient struct {
	*util.FirbolgClient
}

func (c *TrackClient) TrackDeleteMany(req bullet_interface.TrackDeleteMany) error {

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal TrackDeleteMany request: %w", err)
	}

	_, err = c.PostReq("/delete-many", bodyBytes)
	if err != nil {
		return fmt.Errorf("TrackDeleteMany request failed: %w", err)
	}

	return nil
}

func (c *TrackClient) TrackGetManyByPrefix(req bullet_interface.TrackGetItemsByPrefixRequest) (*bullet_interface.TrackGetManyResponse, error) {

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TrackGetManyByPrefix request: %w", err)
	}

	respBytes, err := c.PostReq("/get-query", bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("TrackGetManyByPrefix request failed: %w", err)
	}

	// The Bullet API returns:
	//   { "items": { bucketId -> (key -> TrackValue) } }
	// but `TrackGetManyResponse` expects:
	//   { "values": ..., "missing": ... }
	//
	// Which are different models.
	//
	// So we need a struct matching the API response.
	var apiResp struct {
		Items map[int32]map[string]bullet_interface.TrackValue `json:"items"`
	}

	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TrackGetManyByPrefix response: %w", err)
	}
	// Convert to TrackGetManyResponse:
	out := &bullet_interface.TrackGetManyResponse{
		Values:  apiResp.Items,
		Missing: map[string][]string{}, // prefix query never returns missing entries
	}

	return out, nil
}

func (c *TrackClient) TrackGetMany(req bullet_interface.TrackGetManyRequest) (*bullet_interface.TrackGetManyResponse, error) {

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.PostReq("/get-many", bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var result bullet_interface.TrackGetManyResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}

	return &result, nil
}

func (c *TrackClient) TrackInsertOne(bucketID int32, key string, value int, tag *int64, metric *float64) error {
	reqBody := map[string]interface{}{
		"bucketId": bucketID,
		"key":      key,
		"value":    strconv.Itoa(value),
	}
	if tag != nil {
		reqBody["tag"] = *tag
	}
	if metric != nil {
		reqBody["metric"] = *metric
	}
	bodyBytes, _ := json.Marshal(reqBody)
	_, err := c.PostReq("/insert-one", bodyBytes)
	if err != nil {
		return err
	}
	return nil
}
