package localbullet

import (
	"fmt"

	bullet "github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

// TrackGetMany retrieves multiple keys from multiple buckets
func (l *LocalBullet) TrackGetMany(req bullet.TrackGetManyRequest) (*bullet.TrackGetManyResponse, error) {
	// build map[int32][]string for TrackStore interface
	requestMap := make(map[int32][]string)
	for _, bucket := range req.Buckets {
		requestMap[bucket.BucketID] = bucket.Keys
	}

	// call store
	found, missingMap, err := l.store.TrackGetMany(l.appId, requestMap)
	if err != nil {
		return nil, err
	}

	// convert to client response
	values := make(map[int32]map[string]bullet.TrackValue)
	missing := make(map[string][]string)

	for bucketID, kvMap := range found {
		values[bucketID] = make(map[string]bullet.TrackValue)
		for k, v := range kvMap {
			values[bucketID][k] = bullet.TrackValue{
				Value:  v.Value,
				Tag:    v.Tag,
				Metric: v.Metric,
			}
		}
	}

	for bucketID, keys := range missingMap {
		missing[fmt.Sprintf("%d", bucketID)] = keys
	}

	return &bullet.TrackGetManyResponse{
		Values:  values,
		Missing: missing,
	}, nil
}

// TrackInsertOne inserts a single key-value into a bucket
func (l *LocalBullet) TrackInsertOne(bucketID int32, key string, value int64, tag *int64, metric *float64) error {
	return l.store.TrackPut(l.appId, bucketID, key, value, tag, metric)
}

// TrackDeleteMany deletes multiple keys across buckets
func (l *LocalBullet) TrackDeleteMany(req bullet.TrackDeleteMany) error {
	for _, item := range req.Values {
		if err := l.store.TrackDelete(l.appId, item.BucketID, item.Key); err != nil {
			return err
		}
	}
	return nil
}

// TrackGetManyByPrefix queries by prefix, optionally filtering by tags and metric
func (l *LocalBullet) TrackGetManyByPrefix(req bullet.TrackGetItemsByPrefixRequest) (*bullet.TrackGetManyResponse, error) {
	var metricValue *float64
	var metricIsGt bool

	if req.Metric != nil {
		metricValue = &req.Metric.Value
		metricIsGt = req.Metric.Operator == "gt"
	}

	items, err := l.store.GetItemsByKeyPrefix(l.appId, req.BucketID, req.Prefix, req.Tags, metricValue, metricIsGt)
	if err != nil {
		return nil, err
	}

	values := make(map[int32]map[string]bullet.TrackValue)
	values[req.BucketID] = make(map[string]bullet.TrackValue)

	for _, item := range items {
		values[req.BucketID][item.Key] = bullet.TrackValue{
			Value:  item.Value.Value,
			Tag:    item.Value.Tag,
			Metric: item.Value.Metric,
		}
	}

	return &bullet.TrackGetManyResponse{
		Values:  values,
		Missing: map[string][]string{}, // can't know missing in prefix query
	}, nil
}
