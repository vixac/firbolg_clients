package bullet_stl

import (
	"errors"
	"fmt"
	"strings"

	bullet_interface "github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

/*
A bullet agnostic data structure which allows insertions of subject object (key, value) pairs
*/
type OneWayList interface {
	Upsert(s ListSubject, o ListObject) error
	DeletePair(s ListSubject, o ListObject) error
	DeleteBySub(s ListSubject) error
	GetObject(s ListSubject) (*ListObject, error)
	GetObjectForMany(s []ListSubject) (map[ListSubject]*ListObject, error)
}

// The bullet client implementation of the OnewayList
type BulletOneWayList struct {
	TrackStore   bullet_interface.TrackClientInterface
	BucketId     int32
	ListName     string // It's up to the caller to ensure this is unique across their app
	KeySeparator string //The delimiter used in the key. The caller must ensure this does not appear anwyhere else
}

func NewBulletOneWayList(store bullet_interface.TrackClientInterface, bucketId int32, listName string, separator string) (*BulletOneWayList, error) {
	//VX:TODO check KeySeparator is not part of listName
	return &BulletOneWayList{
		TrackStore:   store,
		BucketId:     bucketId,
		ListName:     listName,
		KeySeparator: separator,
	}, nil
}

// generates the key name. If the object is provided, is it appended
func buildKey(listName string, separator string, subject string, object *string, subjectIsActuallyAPrefix bool) string {
	var key = listName + separator + subject

	//how this works is that the separator at the end acts as a delimter of the end of the key, as in subject:object
	//so if you're looking for all keys that use "sub", you don't want to look for "sub:"
	if !subjectIsActuallyAPrefix {
		key += separator
	}
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
		err := l.DeletePair(s, *existing)
		if err != nil {
			return err
		}
	}

	key := buildKey(l.ListName, l.KeySeparator, s.Value, &o.Value, false)
	return l.TrackStore.TrackInsertOne(l.BucketId, key, 0, nil, nil)
}

func (l *BulletOneWayList) DeleteBySub(s ListSubject) error {
	existing, err := l.GetObject(s)
	if err != nil {
		return err
	}
	if existing != nil {
		return l.DeletePair(s, *existing)
	}
	//nothing to delete.
	return nil
}

func (l *BulletOneWayList) DeletePair(s ListSubject, o ListObject) error {
	key := buildKey(l.ListName, l.KeySeparator, s.Value, &o.Value, false)
	var values []bullet_interface.TrackDeleteValue
	values = append(values, bullet_interface.TrackDeleteValue{
		BucketID: l.BucketId,
		Key:      key,
	})
	return l.TrackStore.TrackDeleteMany(bullet_interface.TrackDeleteMany{
		Values: values,
	})
}

func (l *BulletOneWayList) GetObjectForMany(subjects []ListSubject) (map[ListSubject]*ListObject, error) {
	var keys []string
	for _, s := range subjects {
		prefixKey := buildKey(l.ListName, l.KeySeparator, s.Value, nil, false)
		keys = append(keys, prefixKey)

	}

	prefixReq := bullet_interface.TrackGetItemsbyManyPrefixesRequest{
		BucketID: l.BucketId,
		Prefixes: keys,
	}

	resp, err := l.TrackStore.TrackGetByManyPrefixes(prefixReq)

	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	if _, ok := resp.Values[l.BucketId]; !ok {
		return nil, nil
	}
	values := resp.Values[l.BucketId]
	resMap := make(map[ListSubject]*ListObject)
	for k, _ := range values {
		//ok dammit this is not simple. Its all in the prefix key but we dont know which
		//so we need to trim based on the separator
		split := strings.Split(k, l.KeySeparator)
		if len(split) != 3 {
			fmt.Printf("VX:Error, k = %s, separator is %s, len is %d\n", k, l.KeySeparator, len(split))
			return nil, errors.New("this string did not split into 2")
		}
		subjectValue := split[1]
		objectValue := split[2]
		resMap[ListSubject{Value: subjectValue}] = &ListObject{Value: objectValue}

	}
	return resMap, nil

}
func (l *BulletOneWayList) GetObject(s ListSubject) (*ListObject, error) {
	prefixKey := buildKey(l.ListName, l.KeySeparator, s.Value, nil, false)
	req := bullet_interface.TrackGetItemsByPrefixRequest{
		BucketID: l.BucketId,
		Prefix:   prefixKey,
	}
	res, err := l.TrackStore.TrackGetManyByPrefix(req)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
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
		return nil, errors.New("this two way store got more than 1 item for lookup")
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
