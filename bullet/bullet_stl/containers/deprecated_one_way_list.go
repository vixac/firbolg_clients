package bullet_stl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vixac/bullet/client"
	"github.com/vixac/bullet/model"
)

//VX:TODO test oneway list.
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
	TrackStore   client.Track
	BucketId     int32
	ListName     string // It's up to the caller to ensure this is unique across their app
	KeySeparator string //The delimiter used in the key. The caller must ensure this does not appear anwyhere else
}

func NewBulletOneWayList(store client.Track, bucketId int32, listName string, separator string) (*BulletOneWayList, error) {
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
	return l.TrackStore.TrackPut(l.BucketId, key, model.TrackValue{})
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
	var values []model.TrackKey
	values = append(values, model.TrackKey{
		BucketID: l.BucketId,
		Key:      key,
	})
	return l.TrackStore.TrackDeleteMany(values)
}

func (l *BulletOneWayList) GetObjectForMany(subjects []ListSubject) (map[ListSubject]*ListObject, error) {
	var keys []string
	for _, s := range subjects {
		prefixKey := buildKey(l.ListName, l.KeySeparator, s.Value, nil, false)
		keys = append(keys, prefixKey)

	}

	items, err := l.TrackStore.GetItemsByKeyPrefixes(l.BucketId, keys, nil, nil, false)

	if err != nil {
		return nil, err
	}
	resMap := make(map[ListSubject]*ListObject)
	for _, item := range items {
		//ok dammit this is not simple. Its all in the prefix key but we dont know which
		//so we need to trim based on the separator
		split := strings.Split(item.Key, l.KeySeparator)
		if len(split) != 3 {
			fmt.Printf("VX:Error, key = %s, separator is %s, len is %d\n", item.Key, l.KeySeparator, len(split))
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
	items, err := l.TrackStore.GetItemsByKeyPrefix(l.BucketId, prefixKey, nil, nil, false)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	if len(items) > 1 {
		return nil, errors.New("this two way store got more than 1 item for lookup")
	}
	resultKeyIncludingPrefix := items[0].Key
	object, found := strings.CutPrefix(resultKeyIncludingPrefix, prefixKey)
	if !found {
		return nil, errors.New("invalid result did not contain the prefix")
	}
	return &ListObject{
		Value: object,
	}, nil
}
