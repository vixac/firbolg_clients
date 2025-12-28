package local_bullet

import (
	"github.com/vixac/bullet/model"
	bullet_interface "github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

// store provides a type which has an equivialent in firbolg clients.s
func bulletQueryItemToClient(model model.WayFinderQueryItem) bullet_interface.WayFinderQueryItem {
	return bullet_interface.WayFinderQueryItem{
		Key:     model.Key,
		ItemId:  model.ItemId,
		Tag:     model.Tag,
		Metric:  model.Metric,
		Payload: model.Payload,
	}
}

func bulletItemToClient(model *model.WayFinderGetResponse) *bullet_interface.WayFinderItem {
	if model == nil {
		return nil
	}
	return &bullet_interface.WayFinderItem{
		Item:    model.ItemId,
		Tag:     model.Tag,
		Metric:  model.Metric,
		Payload: model.Payload,
	}
}

func (l *LocalBullet) WayFinderInsertOne(req bullet_interface.WayFinderPutRequest) (int64, error) {
	return l.Store.WayFinderPut(l.Space, req.BucketId, req.Key, req.Payload, req.Tag, req.Metric)
}
func (l *LocalBullet) WayFinderQueryByPrefix(req bullet_interface.WayFinderPrefixQueryRequest) ([]bullet_interface.WayFinderQueryItem, error) {
	res, err := l.Store.WayFinderGetByPrefix(l.Space, req.BucketId, req.Prefix, req.Tags, req.Metric, req.MetricIsGt)
	if err != nil {
		return nil, err
	}

	var mappedRes []bullet_interface.WayFinderQueryItem
	for _, v := range res {
		mappedRes = append(mappedRes, bulletQueryItemToClient(v))
	}
	return mappedRes, nil
}

func (l *LocalBullet) WayFinderGetOne(req bullet_interface.WayFinderGetOneRequest) (*bullet_interface.WayFinderItem, error) {
	res, err := l.Store.WayFinderGetOne(l.Space, req.BucketId, req.Key)
	if err != nil {
		return nil, err
	}
	return bulletItemToClient(res), nil
}
