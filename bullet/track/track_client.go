package bullet

import (
	"encoding/json"
	"fmt"
	"strconv"

	util "github.com/vixac/firbolg_clients/util"
)

type TrackClient struct {
	*util.FirbolgClient
}

func (c *TrackClient) TrackGetMany(req TrackGetManyRequest) (*TrackGetManyResponse, error) {

	// marshal request body
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// execute
	resp, err := c.PostReq("/get-many", bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	// unmarshal
	var result TrackGetManyResponse
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
