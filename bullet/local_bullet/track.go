package local_bullet

import (
	"errors"
	"fmt"
	"github.com/vixac/bullet/store/store_interface"

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
	found, missingMap, err := l.Store.TrackGetMany(l.Space, requestMap)
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
	return l.Store.TrackPut(l.Space, bucketID, key, value, tag, metric)
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
	return l.Store.TrackDeleteMany(l.Space, deleteItems)
}

func (l *LocalBullet) TrackGetByManyPrefixes(
	req bullet_interface.TrackGetItemsbyManyPrefixesRequest,
) (*bullet_interface.TrackGetManyResponse, error) {

	if len(req.Prefixes) == 0 {
		return nil, fmt.Errorf("prefixes must contain at least one prefix")
	}

	var metricValue *float64
	var metricIsGt bool

	if req.Metric != nil {
		metricValue = &req.Metric.Value
		metricIsGt = req.Metric.Operator == "gt"
	}

	items, err := l.Store.GetItemsByKeyPrefixes(
		l.Space,
		req.BucketID,
		req.Prefixes,
		req.Tags,
		metricValue,
		metricIsGt,
	)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	// Build response
	values := make(map[int32]map[string]bullet_interface.TrackValue)
	values[req.BucketID] = make(map[string]bullet_interface.TrackValue)

	for _, item := range items {
		values[req.BucketID][item.Key] = bullet_interface.TrackValue{
			Value:  item.Value.Value,
			Tag:    item.Value.Tag,
			Metric: item.Value.Metric,
		}
	}

	// Missing keys cannot be computed for prefix queries
	return &bullet_interface.TrackGetManyResponse{
		Values:  values,
		Missing: map[string][]string{},
	}, nil
}

// VX:TODO just call manyPrefixes?
// TrackGetManyByPrefix queries by prefix, optionally filtering by tags and metric
func (l *LocalBullet) TrackGetManyByPrefix(req bullet_interface.TrackGetItemsByPrefixRequest) (*bullet_interface.TrackGetManyResponse, error) {
	var metricValue *float64
	var metricIsGt bool

	if req.Metric != nil {
		metricValue = &req.Metric.Value
		metricIsGt = req.Metric.Operator == "gt"
	}

	items, err := l.Store.GetItemsByKeyPrefix(l.Space, req.BucketID, req.Prefix, req.Tags, metricValue, metricIsGt)
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

func (l *LocalBullet) TrackGet(bucketID int32, key string) (int64, error) {
	return l.Store.TrackGet(l.Space, bucketID, key)
}

func (l *LocalBullet) TrackPutMany(req bullet_interface.TrackPutManyRequest) error {
	items := make(map[int32][]model.TrackKeyValueItem)
	for _, bucket := range req.Buckets {
		for _, item := range bucket.Items {
			items[bucket.BucketID] = append(items[bucket.BucketID], model.TrackKeyValueItem{Key: item.Key, Value: model.TrackValue{Value: item.Value, Tag: item.Tag, Metric: item.Metric}})
		}
	}
	return l.Store.TrackPutMany(l.Space, items)
}

func (l *LocalBullet) TrackMutate(req bullet_interface.TrackMutation) (bullet_interface.TrackMutationResult, error) {
	mutation := store_interface.TrackMutation{MutationID: store_interface.MutationID(req.MutationID)}
	for _, put := range req.Puts {
		mutation.Puts = append(mutation.Puts, store_interface.TrackPut{Space: l.Space, BucketID: put.BucketID, Key: put.Key, Value: put.Value, Tag: put.Tag, Metric: put.Metric})
	}
	for _, key := range req.Deletes {
		mutation.Deletes = append(mutation.Deletes, store_interface.TrackKey{Space: l.Space, BucketID: key.BucketID, Key: key.Key})
	}
	result, err := l.Store.TrackMutate(mutation)
	if errors.Is(err, store_interface.ErrTrackMutationUnsupported) {
		return bullet_interface.TrackMutationResult{}, bullet_interface.ErrTrackMutationUnsupported
	}
	return bullet_interface.TrackMutationResult{Applied: result.Applied}, err
}
