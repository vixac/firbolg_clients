package rest_bullet

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	util "github.com/vixac/firbolg_clients/bullet/util"
)

type TrackClient struct {
	*util.FirbolgClient
}

func (c *TrackClient) TrackDeleteMany(req bullet_interface.TrackDeleteMany) error {
	return errors.New("not impl")
}

func (c *TrackClient) TrackGetManyByPrefix(req bullet_interface.TrackGetItemsByPrefixRequest) (*bullet_interface.TrackGetManyResponse, error) {
	return nil, errors.New("not impl")
}

func (c *TrackClient) TrackGetMany(req bullet_interface.TrackGetManyRequest) (*bullet_interface.TrackGetManyResponse, error) {

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
