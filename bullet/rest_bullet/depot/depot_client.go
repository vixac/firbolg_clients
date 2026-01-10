package rest_bullet

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	util "github.com/vixac/firbolg_clients/bullet/util"
)

type DepotClient struct {
	*util.FirbolgClient
}

func (c *DepotClient) DepotInsertOne(req bullet_interface.DepotRequest) error {
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

func (c *DepotClient) DepotGetMany(req bullet_interface.DepotGetManyRequest) (*bullet_interface.DepotGetManyResponse, error) {
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
	var result bullet_interface.DepotGetManyResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}

func (c *DepotClient) DepotUpsertMany(req []bullet_interface.DepotRequest) error {
	return errors.New("not implemented")
}

func (c *DepotClient) DepotDeleteOne(req bullet_interface.DepotDeleteRequest) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	_, err = c.PostReq("/delete-one", bodyBytes)
	if err != nil {
		return err
	}
	return nil
}
