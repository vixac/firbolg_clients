package rest_bullet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"github.com/vixac/firbolg_clients/bullet/util"
)

func (c *RestClient) DepotCreate(req bullet_interface.DepotCreateRequest) (*bullet_interface.DepotCreateResponse, error) {
	bodyBytes, err := util.MarshalJSONBody(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.PostReq("/depot/items", bodyBytes, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	var result bullet_interface.DepotCreateResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}

func (c *RestClient) DepotCreateMany(req bullet_interface.DepotCreateManyRequest) (*bullet_interface.DepotCreateManyResponse, error) {
	bodyBytes, err := util.MarshalJSONBody(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.PostReq("/depot/items/batch", bodyBytes, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	var result bullet_interface.DepotCreateManyResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}

func (c *RestClient) DepotUpdate(req bullet_interface.DepotUpdateRequest) error {
	bodyBytes, err := util.MarshalJSONBody(struct {
		Value string `json:"value"`
	}{
		Value: req.Value,
	})
	if err != nil {
		return err
	}
	_, err = c.PutReq("/depot/items/"+strconv.FormatInt(req.ID, 10), bodyBytes, http.StatusOK)
	return err
}

func (c *RestClient) DepotGetOne(req bullet_interface.DepotGetRequest) (*bullet_interface.DepotGetResponse, error) {
	resp, err := c.GetReq("/depot/items/"+strconv.FormatInt(req.ID, 10), http.StatusOK)
	if err != nil {
		return nil, err
	}
	var result bullet_interface.DepotGetResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}

func (c *RestClient) DepotGetMany(req bullet_interface.DepotGetManyRequest) (*bullet_interface.DepotGetManyResponse, error) {
	bodyBytes, err := util.MarshalJSONBody(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.PostReq("/depot/items/batch-get", bodyBytes, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var result bullet_interface.DepotGetManyResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}

func (c *RestClient) DepotDelete(req bullet_interface.DepotDeleteRequest) error {
	_, err := c.DeleteReq("/depot/items/"+strconv.FormatInt(req.ID, 10), nil, http.StatusNoContent)
	return err
}

func (c *RestClient) DepotDeleteByBucket(req bullet_interface.DepotBucketRequest) error {
	_, err := c.DeleteReq("/depot/bucket/"+strconv.FormatInt(int64(req.BucketID), 10), nil, http.StatusNoContent)
	return err
}

func (c *RestClient) DepotGetAllByBucket(req bullet_interface.DepotBucketRequest) (*bullet_interface.DepotGetAllByBucketResponse, error) {
	resp, err := c.GetReq("/depot/bucket/"+strconv.FormatInt(int64(req.BucketID), 10), http.StatusOK)
	if err != nil {
		return nil, err
	}
	var result bullet_interface.DepotGetAllByBucketResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result, nil
}
