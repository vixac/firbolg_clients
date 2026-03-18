package rest_bullet

import (
	"encoding/json"
	"fmt"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	util "github.com/vixac/firbolg_clients/bullet/util"
)

type DepotClient struct {
	*util.FirbolgClient
}

func (c *DepotClient) DepotCreate(req bullet_interface.DepotCreateRequest) (*bullet_interface.DepotCreateResponse, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	resp, err := c.PostReq("/create-one", bodyBytes)
	if err != nil {
		return nil, err
	}
	var result bullet_interface.DepotCreateResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}

func (c *DepotClient) DepotCreateMany(req bullet_interface.DepotCreateManyRequest) (*bullet_interface.DepotCreateManyResponse, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	resp, err := c.PostReq("/create-many", bodyBytes)
	if err != nil {
		return nil, err
	}
	var result bullet_interface.DepotCreateManyResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}

func (c *DepotClient) DepotUpdate(req bullet_interface.DepotUpdateRequest) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	_, err = c.PostReq("/update", bodyBytes)
	return err
}

func (c *DepotClient) DepotGetOne(req bullet_interface.DepotGetRequest) (*bullet_interface.DepotGetResponse, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	resp, err := c.PostReq("/get-one", bodyBytes)
	if err != nil {
		return nil, err
	}
	var result bullet_interface.DepotGetResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}

func (c *DepotClient) DepotGetMany(req bullet_interface.DepotGetManyRequest) (*bullet_interface.DepotGetManyResponse, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	resp, err := c.PostReq("/get-many", bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var result bullet_interface.DepotGetManyResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}

func (c *DepotClient) DepotDelete(req bullet_interface.DepotDeleteRequest) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	_, err = c.PostReq("/delete-one", bodyBytes)
	return err
}

func (c *DepotClient) DepotDeleteByBucket(req bullet_interface.DepotBucketRequest) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	_, err = c.PostReq("/delete-by-bucket", bodyBytes)
	return err
}

func (c *DepotClient) DepotGetAllByBucket(req bullet_interface.DepotBucketRequest) (*bullet_interface.DepotGetAllByBucketResponse, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	resp, err := c.PostReq("/get-all-by-bucket", bodyBytes)
	if err != nil {
		return nil, err
	}
	var result bullet_interface.DepotGetAllByBucketResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}
