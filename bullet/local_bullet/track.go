package local_bullet

import (
	"fmt"

	"github.com/vixac/bullet/model"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

// TrackGetMany retrieves multiple keys from multiple buckets
func (l *LocalBullet) TrackGetMany(req bullet_interface.TrackGetManyRequest) (*bullet_interface.TrackGetManyResponse, error) {
	// build map[int32][]string for TrackStore interface
	requestMap := make(map[int32][]string)
	for _, bucket := range req.Buckets {
		requestMap[bucket.BucketID] = bucket.Keys
	}

	// call store
	found, missingMap, err := l.Store.TrackGetMany(l.AppId, requestMap)
	if err != nil {
		return nil, err
	}

	// convert to client response
	values := make(map[int32]map[string]bullet_interface.TrackValue)
	missing := make(map[string][]string)

	for bucketID, kvMap := range found {
		values[bucketID] = make(map[string]bullet_interface.TrackValue)
		for k, v := range kvMap {
			values[bucketID][k] = bullet_interface.TrackValue{
				Value:  v.Value,
				Tag:    v.Tag,
				Metric: v.Metric,
			}
		}
	}

	for bucketID, keys := range missingMap {
		missing[fmt.Sprintf("%d", bucketID)] = keys
	}

	return &bullet_interface.TrackGetManyResponse{
		Values:  values,
		Missing: missing,
	}, nil
}

// TrackInsertOne inserts a single key-value into a bucket
func (l *LocalBullet) TrackInsertOne(bucketID int32, key string, value int64, tag *int64, metric *float64) error {
	return l.Store.TrackPut(l.AppId, bucketID, key, value, tag, metric)
}

// TrackDeleteMany deletes multiple keys across buckets
func (l *LocalBullet) TrackDeleteMany(req bullet_interface.TrackDeleteMany) error {

	var deleteItems []model.TrackBucketKeyPair
	for _, item := range req.Values {
		deleteItems = append(deleteItems, model.TrackBucketKeyPair{
			BucketID: item.BucketID,
			Key:      item.Key,
		})

	}
	return l.Store.TrackDeleteMany(l.AppId, deleteItems)
}

// TrackGetManyByPrefix queries by prefix, optionally filtering by tags and metric
func (l *LocalBullet) TrackGetManyByPrefix(req bullet_interface.TrackGetItemsByPrefixRequest) (*bullet_interface.TrackGetManyResponse, error) {
	var metricValue *float64
	var metricIsGt bool

	if req.Metric != nil {
		metricValue = &req.Metric.Value
		metricIsGt = req.Metric.Operator == "gt"
	}

	items, err := l.Store.GetItemsByKeyPrefix(l.AppId, req.BucketID, req.Prefix, req.Tags, metricValue, metricIsGt)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	values := make(map[int32]map[string]bullet_interface.TrackValue)
	values[req.BucketID] = make(map[string]bullet_interface.TrackValue)

	for _, item := range items {
		values[req.BucketID][item.Key] = bullet_interface.TrackValue{
			Value:  item.Value.Value,
			Tag:    item.Value.Tag,
			Metric: item.Value.Metric,
		}
	}

	return &bullet_interface.TrackGetManyResponse{
		Values:  values,
		Missing: map[string][]string{}, // can't know missing in prefix query
	}, nil
}
