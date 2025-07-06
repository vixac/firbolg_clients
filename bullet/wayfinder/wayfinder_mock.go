package bullet

import (
	"fmt"
	"sync"
)

type Bucket struct {
	id int32
}
type WayFinderMockClient struct {
	Data map[Bucket]map[string]WayFinderItem
	mu   sync.Mutex
}

var nextId = int64(1)

func (m *WayFinderMockClient) WayFinderInsertOne(req WayFinderPutRequest) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := int64(len(m.Data) + 1)
	bucket := Bucket{id: req.BucketId}
	if m.Data == nil {
		m.Data = make(map[Bucket]map[string]WayFinderItem)
	}
	if _, ok := m.Data[bucket]; !ok {
		m.Data[bucket] = make(map[string]WayFinderItem)
	}

	m.Data[bucket][req.Key] = WayFinderItem{
		Tag:     req.Tag,
		Metric:  req.Metric,
		Payload: req.Payload,
		Item:    nextId,
	}
	nextId = nextId + 1
	return id, nil
}

func (m *WayFinderMockClient) WayFinderGetOne(req WayFinderGetOneRequest) (*WayFinderItem, error) {
	fmt.Printf("data is %+v", m.Data)
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

func (c *WayFinderMockClient) WayFinderQueryByPrefix(req WayFinderPrefixQueryRequest) ([]WayFinderQueryItem, error) {
	return nil, nil
}
