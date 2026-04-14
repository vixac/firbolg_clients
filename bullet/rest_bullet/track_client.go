package rest_bullet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	bullet_model "github.com/vixac/bullet/model"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"github.com/vixac/firbolg_clients/bullet/util"
)

func (c *RestClient) TrackDeleteMany(req bullet_interface.TrackDeleteMany) error {
	bodyBytes, err := util.MarshalJSONBody(req)
	if err != nil {
		return err
	}
	_, err = c.DeleteReq("/track/items", bodyBytes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("TrackDeleteMany request failed: %w", err)
	}
	return nil
}

func (c *RestClient) TrackGetManyByPrefix(req bullet_interface.TrackGetItemsByPrefixRequest) (*bullet_interface.TrackGetManyResponse, error) {
	bodyBytes, err := util.MarshalJSONBody(req)
	if err != nil {
		return nil, err
	}
	respBytes, err := c.PostReq("/track/query", bodyBytes, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("TrackGetManyByPrefix request failed: %w", err)
	}

	var apiResp struct {
		Items []bullet_model.TrackKeyValueItem `json:"items"`
	}
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TrackGetManyByPrefix response: %w", err)
	}

	return trackItemsResponse(req.BucketID, apiResp.Items), nil
}

func (c *RestClient) TrackGetMany(req bullet_interface.TrackGetManyRequest) (*bullet_interface.TrackGetManyResponse, error) {
	bodyBytes, err := util.MarshalJSONBody(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.PostReq("/track/items/batch-get", bodyBytes, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var apiResp bullet_model.TrackGetManyResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}

	return convertTrackGetManyResponse(apiResp)
}

func (c *RestClient) TrackInsertOne(bucketID int32, key string, value int64, tag *int64, metric *float64) error {
	reqBody := bullet_interface.TrackRequest{
		BucketID: bucketID,
		Key:      key,
		Value:    value,
		Tag:      tag,
		Metric:   metric,
	}
	bodyBytes, err := util.MarshalJSONBody(reqBody)
	if err != nil {
		return err
	}
	_, err = c.PostReq("/track/items", bodyBytes, http.StatusOK)
	if err != nil {
		return fmt.Errorf("TrackInsertOne request failed: %w", err)
	}
	return nil
}

func (c *RestClient) TrackGetByManyPrefixes(req bullet_interface.TrackGetItemsbyManyPrefixesRequest) (*bullet_interface.TrackGetManyResponse, error) {
	bodyBytes, err := util.MarshalJSONBody(req)
	if err != nil {
		return nil, err
	}
	respBytes, err := c.PostReq("/track/query/multi", bodyBytes, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("TrackGetByManyPrefixes request failed: %w", err)
	}

	var apiResp struct {
		Items []bullet_model.TrackKeyValueItem `json:"items"`
	}
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TrackGetByManyPrefixes response: %w", err)
	}

	return trackItemsResponse(req.BucketID, apiResp.Items), nil
}

func trackItemsResponse(bucketID int32, items []bullet_model.TrackKeyValueItem) *bullet_interface.TrackGetManyResponse {
	values := make(map[int32]map[string]bullet_interface.TrackValue)
	values[bucketID] = make(map[string]bullet_interface.TrackValue, len(items))
	for _, item := range items {
		values[bucketID][item.Key] = bullet_interface.TrackValue{
			Value:  item.Value.Value,
			Tag:    item.Value.Tag,
			Metric: item.Value.Metric,
		}
	}

	return &bullet_interface.TrackGetManyResponse{
		Values:  values,
		Missing: map[string][]string{},
	}
}

func convertTrackGetManyResponse(apiResp bullet_model.TrackGetManyResponse) (*bullet_interface.TrackGetManyResponse, error) {
	values := make(map[int32]map[string]bullet_interface.TrackValue, len(apiResp.Values))
	for bucketID, bucketValues := range apiResp.Values {
		parsedBucketID, err := strconv.ParseInt(bucketID, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid bucket ID %q in response: %w", bucketID, err)
		}

		values[int32(parsedBucketID)] = make(map[string]bullet_interface.TrackValue, len(bucketValues))
		for key, item := range bucketValues {
			values[int32(parsedBucketID)][key] = bullet_interface.TrackValue{
				Value:  item.Value,
				Tag:    item.Tag,
				Metric: item.Metric,
			}
		}
	}

	return &bullet_interface.TrackGetManyResponse{
		Values:  values,
		Missing: apiResp.Missing,
	}, nil
}
