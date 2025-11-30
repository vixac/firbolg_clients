package local_bullet

import (
	"github.com/vixac/bullet/model"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

func (l *LocalBullet) DepotInsertOne(req bullet_interface.DepotRequest) error {
	return l.Store.DepotPut(l.AppId, req.Key, req.Value)
}

func (l *LocalBullet) DepotGetMany(req bullet_interface.DepotGetManyRequest) (*bullet_interface.DepotGetManyResponse, error) {

	res, missing, err := l.Store.DepotGetMany(l.AppId, req.Keys)
	if err != nil {
		return nil, err
	}

	return &bullet_interface.DepotGetManyResponse{
		Values:  res,
		Missing: missing,
	}, nil
}

func (l *LocalBullet) DepotUpsertMany(req []bullet_interface.DepotRequest) error {
	var items []model.DepotKeyValueItem
	for _, v := range req {
		items = append(items, model.DepotKeyValueItem{
			Key:   v.Key,
			Value: v.Value,
		})
	}
	return l.Store.DepotPutMany(l.AppId, items)
}
