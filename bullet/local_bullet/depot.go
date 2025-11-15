package localbullet

import (
	depot "github.com/vixac/firbolg_clients/bullet/depot"
)

func (l *LocalBullet) DepotInsertOne(req depot.DepotRequest) error {
	return l.store.DepotPut(l.appId, req.Key, req.Value)
}

func (l *LocalBullet) DepotGetMany(req depot.DepotGetManyRequest) (*depot.DepotGetManyResponse, error) {

	res, missing, err := l.store.DepotGetMany(l.appId, req.Keys)
	if err != nil {
		return nil, err
	}

	return &depot.DepotGetManyResponse{
		Values:  res,
		Missing: missing,
	}, nil
}
