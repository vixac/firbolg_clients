package bullet_stl

import (
	"errors"
	"strings"

	track "github.com/vixac/firbolg_clients/bullet/track"
)

/*
*
A bullet agnostic data structure which allows insertions of subject object (key, value) pairs
*/
type OneWayList interface {
	Upsert(s ListSubject, o ListObject) error
	DeleteViaSub(s ListSubject) error
	//DeletePair(s ListSubject, o ListObject) error dont think we need it in the end.
	GetObject(s ListSubject) (*ListObject, error)
}

// The bullet client implementation of the OnewayList
type BulletOneWayList struct {
	TrackStore   track.TrackClientInterface
	BucketId     int32
	ListName     string // It's up to the caller to ensure this is unique across their app
	KeySeparator string //The delimiter used in the key. The caller must ensure this does not appear anwyhere else
}

func NewBulletOneWayList(store track.TrackClientInterface, bucketId int32, listName string, separator string) (*BulletOneWayList, error) {
	//VX:TODO check KeySeparator is not part of listName
	return &BulletOneWayList{
		TrackStore:   store,
		BucketId:     bucketId,
		ListName:     listName,
		KeySeparator: separator,
	}, nil
}

// VX:TODO test
// generates the key name. If the object is provided, is it appended
func buildKey(listName string, separator string, subject string, object *string) string {
	var key = listName + separator + subject + separator //note the separator on the end even if theres no object "a:b:c" or "a:b:"
	if object != nil {
		key = key + *object
	}
	return key
}

func (l *BulletOneWayList) Upsert(s ListSubject, o ListObject) error {
	//VX:TODO check keySepawrator is not used in names
	existing, err := l.GetObject(s)
	if err != nil {
		return nil
	}
	//delete the key if it exists.
	if existing != nil {
		err := l.DeleteViaSub(s)
		if err != nil {
			return err
		}
	}

	key := buildKey(l.ListName, l.KeySeparator, s.Value, &o.Value)
	return l.TrackStore.TrackInsertOne(l.BucketId, key, 0, nil, nil)
}

// VX:TODO test
func (l *BulletOneWayList) DeleteViaSub(s ListSubject) error {
	key := buildKey(l.ListName, l.KeySeparator, s.Value, nil)
	var values []track.TrackDeleteValue
	values = append(values, track.TrackDeleteValue{
		BucketID: l.BucketId,
		Key:      key,
	})
	return l.TrackStore.TrackDeleteMany(track.TrackDeleteMany{
		Values: values,
	})
}

// VX:TODO test
func (l *BulletOneWayList) GetObject(s ListSubject) (*ListObject, error) {
	prefixKey := buildKey(l.ListName, l.KeySeparator, s.Value, nil)

	req := track.TrackGetItemsByPrefixRequest{
		BucketID: l.BucketId,
		Prefix:   prefixKey,
	}
	res, err := l.TrackStore.TrackGetManyByPrefix(req)
	if err != nil {
		return nil, err
	}
	if len(res.Values) != 1 {
		return nil, errors.New("missing bucket")
	}
	itemsByBucket := res.Values[l.BucketId]
	if itemsByBucket == nil {
		return nil, nil // its ok to get and find nothing
	}

	//if we're here we assume the object exists, so its an error if its not where we expect it
	itemsInBucket := make([]string, 0, len(itemsByBucket))
	for k := range itemsByBucket {
		itemsInBucket = append(itemsInBucket, k)
	}

	if len(itemsInBucket) == 0 {
		return nil, errors.New("missing item in bucket")
	}
	if len(itemsByBucket) > 1 {
		return nil, errors.New("This two way store got more than 1 item for lookup")
	}
	resultKeyIncludingPrefix := itemsInBucket[0]
	object, found := strings.CutPrefix(resultKeyIncludingPrefix, prefixKey)
	if !found {
		return nil, errors.New("invalid result did not contain the prefix")
	}
	return &ListObject{
		Value: object,
	}, nil
}
