package local_bullet

import (
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
