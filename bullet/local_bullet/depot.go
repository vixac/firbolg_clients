package localbullet

import (
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

func (l *LocalBullet) DepotInsertOne(req bullet_interface.DepotRequest) error {
	return l.store.DepotPut(l.appId, req.Key, req.Value)
}

func (l *LocalBullet) DepotGetMany(req bullet_interface.DepotGetManyRequest) (*bullet_interface.DepotGetManyResponse, error) {

	res, missing, err := l.store.DepotGetMany(l.appId, req.Keys)
	if err != nil {
		return nil, err
	}

	return &bullet_interface.DepotGetManyResponse{
		Values:  res,
		Missing: missing,
	}, nil
}
