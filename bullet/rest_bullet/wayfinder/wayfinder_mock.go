package rest_bullet

import (
	"sync"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

type Bucket struct {
	id int32
}

// VX:TODO not used
type WayFinderMockClient struct {
	Data map[Bucket]map[string]bullet_interface.WayFinderItem
	mu   sync.Mutex
}

func (m *WayFinderMockClient) WayFinderInsertOne(req bullet_interface.WayFinderPutRequest) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := int64(len(m.Data) + 1)
	bucket := Bucket{id: req.BucketId}
	if m.Data == nil {
		m.Data = make(map[Bucket]map[string]bullet_interface.WayFinderItem)
	}
	if _, ok := m.Data[bucket]; !ok {
		m.Data[bucket] = make(map[string]bullet_interface.WayFinderItem)
	}

	m.Data[bucket][req.Key] = bullet_interface.WayFinderItem{
		Tag:     req.Tag,
		Metric:  req.Metric,
		Payload: req.Payload,
		Item:    id,
	}
	return id, nil
}

func (m *WayFinderMockClient) WayFinderGetOne(req bullet_interface.WayFinderGetOneRequest) (*bullet_interface.WayFinderItem, error) {
	if m.Data == nil {
		return nil, nil
	}
	keys, ok := m.Data[Bucket{id: req.BucketId}]
	if !ok {
		return nil, nil
	}
	if keys == nil {
		return nil, nil
	}
	item, ok := keys[req.Key]
	if !ok {
		return nil, nil
	}
	return &item, nil
}

func (c *WayFinderMockClient) WayFinderQueryByPrefix(req bullet_interface.WayFinderPrefixQueryRequest) ([]bullet_interface.WayFinderQueryItem, error) {
	return nil, nil
}
