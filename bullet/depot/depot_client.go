package bullet

import (
	"encoding/json"
	"fmt"

	util "github.com/vixac/firbolg_clients/util"
)

type DepotClient struct {
	*util.FirbolgClient
}

func (c *DepotClient) DepotInsertOne(req DepotRequest) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	_, err = c.PostReq("/insert-one", bodyBytes)
	if err != nil {
		return err
	}
	return nil
}

func (c *DepotClient) DepotGetMany(req DepotGetManyRequest) (*DepotGetManyResponse, error) {
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
	var result DepotGetManyResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}
