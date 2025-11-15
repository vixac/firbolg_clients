package localbullet

import (
	"github.com/vixac/bullet/model"
	wayfinder "github.com/vixac/firbolg_clients/bullet/wayfinder"
)

// store provides a type which has an equivialent in firbolg clients.s
func bulletQueryItemToClient(model model.WayFinderQueryItem) wayfinder.WayFinderQueryItem {
	return wayfinder.WayFinderQueryItem{
		Key:     model.Key,
		ItemId:  model.ItemId,
		Tag:     model.Tag,
		Metric:  model.Metric,
		Payload: model.Payload,
	}
}

func bulletItemToClient(model *model.WayFinderGetResponse) *wayfinder.WayFinderItem {
	if model == nil {
		return nil
	}
	return &wayfinder.WayFinderItem{
		Item:    model.ItemId,
		Tag:     model.Tag,
		Metric:  model.Metric,
		Payload: model.Payload,
	}
}

func (l *LocalBullet) WayFinderInsertOne(req wayfinder.WayFinderPutRequest) (int64, error) {
	return l.store.WayFinderPut(l.appId, req.BucketId, req.Key, req.Payload, req.Tag, req.Metric)
}
func (l *LocalBullet) WayFinderQueryByPrefix(req wayfinder.WayFinderPrefixQueryRequest) ([]wayfinder.WayFinderQueryItem, error) {
	res, err := l.store.WayFinderGetByPrefix(l.appId, req.BucketId, req.Prefix, req.Tags, req.Metric, req.MetricIsGt)
	if err != nil {
		return nil, err
	}

	var mappedRes []wayfinder.WayFinderQueryItem
	for _, v := range res {
		mappedRes = append(mappedRes, bulletQueryItemToClient(v))
	}
	return mappedRes, nil
}

func (l *LocalBullet) WayFinderGetOne(req wayfinder.WayFinderGetOneRequest) (*wayfinder.WayFinderItem, error) {
	res, err := l.store.WayFinderGetOne(l.appId, req.BucketId, req.Key)
	if err != nil {
		return nil, err
	}
	return bulletItemToClient(res), nil
}
