package bullet_stl

import (
	"fmt"

	bullet_interface "github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

/*
This will resemble wayfinder in functionality.
*/
type Collection interface {
	CreateItemUnder(key string, payload string) (*CollectionId, error)
	AllItems() (map[CollectionId]string, error) //VX:Note this can be upgraded to have paging
	AllItemsUnderPrefix(prefix string) (map[CollectionId]string, error)
	DeleteItems(ids []CollectionId) error //VX:Note delete payload first as it has the less bad edge case
}

type CollectionId struct {
	Bucket  int32
	DepotId int64
	Key     string
}

type BulletCollection struct {
	TrackStore bullet_interface.TrackClientInterface
	DepotStore bullet_interface.DepotClientInterface
	BucketId   int32
}

func NewBulletCollection(bucket int32, track bullet_interface.TrackClientInterface, depot bullet_interface.DepotClientInterface) Collection {
	return &BulletCollection{
		TrackStore: track,
		DepotStore: depot,
		BucketId:   bucket,
	}
}

func (b *BulletCollection) CreateItemUnder(key string, payload string) (*CollectionId, error) {
	depotReq := bullet_interface.DepotCreateRequest{
		BucketID: b.BucketId,
		Value:    payload,
	}
	depotResponse, err := b.DepotStore.DepotCreate(depotReq)
	if err != nil {
		return nil, err
	}
	err = b.TrackStore.TrackInsertOne(b.BucketId, key, depotResponse.ID, nil, nil)
	if err != nil {
		fmt.Printf("VX: Warn this is an inconsistent state. We have an orphan depot item: %d\n", depotResponse.ID)
		return nil, err
	}
	collectionId := CollectionId{
		Bucket:  b.BucketId,
		DepotId: depotResponse.ID,
		Key:     key,
	}
	return &collectionId, nil
}

func (b *BulletCollection) AllItems() (map[CollectionId]string, error) {
	return b.AllItemsUnderPrefix("")
}

func (b *BulletCollection) AllItemsUnderPrefix(prefix string) (map[CollectionId]string, error) {
	trackReq := bullet_interface.TrackGetItemsByPrefixRequest{
		BucketID: b.BucketId,
		Prefix:   prefix,
	}

	//use track to get all the ids that fall under this collection.
	res, err := b.TrackStore.TrackGetManyByPrefix(trackReq)
	if err != nil || res == nil {
		return nil, err
	}

	bucket, ok := res.Values[b.BucketId]

	if !ok {
		return nil, nil
	}
	var depotIds []int64
	depotIdsToTrackKey := make(map[int64]string)
	for k, v := range bucket {
		depotIds = append(depotIds, v.Value)
		depotIdsToTrackKey[v.Value] = k //we need to find the track key given the depot id later to create the collection id
	}

	manyReq := bullet_interface.DepotGetManyRequest{
		IDs: depotIds,
	}
	depotRes, err := b.DepotStore.DepotGetMany(manyReq)
	if err != nil || depotRes == nil {
		return nil, err
	}
	if len(depotRes.Missing) > 0 {
		fmt.Printf("VX: WARN there are %d missing Ids in this collection. \n", len(depotRes.Missing))
	}

	result := make(map[CollectionId]string)
	for k, payload := range depotRes.Values {
		trackKey := depotIdsToTrackKey[k]
		col := CollectionId{
			Bucket:  b.BucketId,
			DepotId: k,
			Key:     trackKey,
		}
		result[col] = payload
	}
	return result, nil
}
func (b *BulletCollection) DeleteItems(ids []CollectionId) error {
	//VX:TODO delete many
	for _, v := range ids {
		req := bullet_interface.DepotDeleteRequest{
			ID: v.DepotId,
		}
		err := b.DepotStore.DepotDelete(req)
		if err != nil {
			return err
		}
	}
	return nil
}
