package bullet_stl

import (
	"fmt"
	"time"

	bullet_interface "github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

/*
This will resemble wayfinder in functionality.
*/
type Collection interface {
	//the updatedTime is stored in the metric field of track. If you want the created time, consider
	//baking that into your key or your payload.
	CreateItemUnder(key string, payload string, updateTime *time.Time) (*CollectionId, error)
	EditPayload(id CollectionId, payload string, updateTime *time.Time) error
	AllItems() (map[CollectionId]string, error) //VX:Note this can be upgraded to have paging
	AllItemsUnderPrefix(prefix string) (map[CollectionId]CollectionItem, error)
	AllItemsUnderPrefixes(prefixes []string) (map[CollectionId]CollectionItem, error)
	ItemsForKeys(keys []string) (map[CollectionId]CollectionItem, error)
	DeleteItems(ids []CollectionId) error //VX:Note delete payload first as it has the less bad edge case
}

// the edited time is baked into the
type CollectionItem struct {
	Payload string
	Updated time.Time
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

func (b *BulletCollection) CreateItemUnder(key string, payload string, updateTime *time.Time) (*CollectionId, error) {
	depotReq := bullet_interface.DepotCreateRequest{
		BucketID: b.BucketId,
		Value:    payload,
	}
	depotResponse, err := b.DepotStore.DepotCreate(depotReq)
	if err != nil {
		return nil, err
	}
	var updateTimeUnix *float64 = nil
	if updateTime != nil {
		updateUnix := float64(updateTime.Unix())
		updateTimeUnix = &updateUnix
	}
	err = b.TrackStore.TrackInsertOne(b.BucketId, key, depotResponse.ID, nil, updateTimeUnix)
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

func (b *BulletCollection) EditPayload(id CollectionId, payload string, updateTime *time.Time) error {
	err := b.DepotStore.DepotUpdate(bullet_interface.DepotUpdateRequest{
		ID:    id.DepotId,
		Value: payload,
	})
	if err != nil {
		return err
	}

	var updateTimeUnix *float64 = nil
	if updateTime != nil {
		updateUnix := float64(updateTime.Unix())
		updateTimeUnix = &updateUnix
	}
	return b.TrackStore.TrackInsertOne(b.BucketId, id.Key, id.DepotId, nil, updateTimeUnix)
}

func (b *BulletCollection) AllItems() (map[CollectionId]string, error) {
	items, err := b.AllItemsUnderPrefix("")
	if err != nil || items == nil {
		return nil, err
	}
	result := make(map[CollectionId]string)
	for k, v := range items {
		result[k] = v.Payload
	}
	return result, nil
}

type trackMeta struct {
	Key    string
	Tag    *int64
	Metric *float64
}

func (b *BulletCollection) fetchItemsFor(res *bullet_interface.TrackGetManyResponse) (map[CollectionId]CollectionItem, error) {
	bucket, ok := res.Values[b.BucketId]
	if !ok {
		return nil, nil
	}
	var depotIds []int64
	depotIdsToMeta := make(map[int64]trackMeta)
	for k, v := range bucket {
		depotIds = append(depotIds, v.Value)
		depotIdsToMeta[v.Value] = trackMeta{Key: k, Tag: v.Tag, Metric: v.Metric}
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

	result := make(map[CollectionId]CollectionItem)
	for k, payload := range depotRes.Values {
		meta := depotIdsToMeta[k]
		col := CollectionId{
			Bucket:  b.BucketId,
			DepotId: k,
			Key:     meta.Key,
		}
		item := CollectionItem{Payload: payload}
		if meta.Metric != nil {
			item.Updated = time.Unix(int64(*meta.Metric), 0)
		}
		result[col] = item
	}
	return result, nil
}

func (b *BulletCollection) AllItemsUnderPrefixes(prefixes []string) (map[CollectionId]CollectionItem, error) {
	trackReq := bullet_interface.TrackGetItemsbyManyPrefixesRequest{
		BucketID: b.BucketId,
		Prefixes: prefixes,
	}

	//use track to get all the ids that fall under this collection.
	res, err := b.TrackStore.TrackGetByManyPrefixes(trackReq)
	if err != nil || res == nil {
		return nil, err
	}
	return b.fetchItemsFor(res)
}

func (b *BulletCollection) AllItemsUnderPrefix(prefix string) (map[CollectionId]CollectionItem, error) {
	trackReq := bullet_interface.TrackGetItemsByPrefixRequest{
		BucketID: b.BucketId,
		Prefix:   prefix,
	}

	//use track to get all the ids that fall under this collection.
	res, err := b.TrackStore.TrackGetManyByPrefix(trackReq)
	if err != nil || res == nil {
		return nil, err
	}

	return b.fetchItemsFor(res)
}

func (b *BulletCollection) ItemsForKeys(keys []string) (map[CollectionId]CollectionItem, error) {
	req := bullet_interface.TrackGetManyRequest{
		Buckets: []bullet_interface.TrackGetKeys{
			{BucketID: b.BucketId, Keys: keys},
		},
	}
	res, err := b.TrackStore.TrackGetMany(req)
	if err != nil || res == nil {
		return nil, err
	}
	return b.fetchItemsFor(res)
}

func (b *BulletCollection) DeleteItems(ids []CollectionId) error {
	// Delete depot first: if track delete fails, orphaned track entries are the less bad edge case
	for _, v := range ids {
		req := bullet_interface.DepotDeleteRequest{
			ID: v.DepotId,
		}
		err := b.DepotStore.DepotDelete(req)
		if err != nil {
			return err
		}
	}
	var trackDeletes []bullet_interface.TrackDeleteValue
	for _, v := range ids {
		trackDeletes = append(trackDeletes, bullet_interface.TrackDeleteValue{
			BucketID: b.BucketId,
			Key:      v.Key,
		})
	}
	return b.TrackStore.TrackDeleteMany(bullet_interface.TrackDeleteMany{Values: trackDeletes})
}
