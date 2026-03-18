package local_bullet

import (
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

func (l *LocalBullet) DepotCreate(req bullet_interface.DepotCreateRequest) (*bullet_interface.DepotCreateResponse, error) {
	id, err := l.Store.DepotCreate(l.Space, req.BucketID, req.Value)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.DepotCreateResponse{ID: id}, nil
}

func (l *LocalBullet) DepotCreateMany(req bullet_interface.DepotCreateManyRequest) (*bullet_interface.DepotCreateManyResponse, error) {
	ids, err := l.Store.DepotCreateMany(l.Space, req.BucketID, req.Values)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.DepotCreateManyResponse{IDs: ids}, nil
}

func (l *LocalBullet) DepotUpdate(req bullet_interface.DepotUpdateRequest) error {
	return l.Store.DepotUpdate(l.Space, req.ID, req.Value)
}

func (l *LocalBullet) DepotGetOne(req bullet_interface.DepotGetRequest) (*bullet_interface.DepotGetResponse, error) {
	value, err := l.Store.DepotGet(l.Space, req.ID)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.DepotGetResponse{Value: value}, nil
}

func (l *LocalBullet) DepotGetMany(req bullet_interface.DepotGetManyRequest) (*bullet_interface.DepotGetManyResponse, error) {
	values, missing, err := l.Store.DepotGetMany(l.Space, req.IDs)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.DepotGetManyResponse{
		Values:  values,
		Missing: missing,
	}, nil
}

func (l *LocalBullet) DepotDelete(req bullet_interface.DepotDeleteRequest) error {
	return l.Store.DepotDelete(l.Space, req.ID)
}

func (l *LocalBullet) DepotDeleteByBucket(req bullet_interface.DepotBucketRequest) error {
	return l.Store.DepotDeleteByBucket(l.Space, req.BucketID)
}

func (l *LocalBullet) DepotGetAllByBucket(req bullet_interface.DepotBucketRequest) (*bullet_interface.DepotGetAllByBucketResponse, error) {
	values, err := l.Store.DepotGetAllByBucket(l.Space, req.BucketID)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.DepotGetAllByBucketResponse{Values: values}, nil
}
