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
type TwoWayList interface {
	Upsert(s ListSubject, o ListObject) error
	DeleteViaSub(s ListSubject) error
	DeleteViaObj(o ListObject) error
	GetObjectViaSubject(s ListSubject) (*ListObject, error)
	GetObjectViaObject(s ListSubject) (*ListObject, error)
}

// The bullet client implementation of the OnewayList
type BulletTwoWayList struct {
	TrackStore        track.TrackClientInterface
	BucketId          int32
	ListName          string // It's up to the caller to ensure this is unique across their app
	ForwardSeparator  string //Used for S -> O mappings
	BackwardSeparator string //Used for O -> S mappings
}

func NewBulletTwoWayList(store track.TrackClientInterface, bucketId int32, listName string, forwardSeparator string, backwardSeparator string) (*BulletTwoWayList, error) {
	//VX:TODO check separators are unique, and are not in the listName
	return &BulletTwoWayList{
		TrackStore:        store,
		BucketId:          bucketId,
		ListName:          listName,
		ForwardSeparator:  forwardSeparator,
		BackwardSeparator: backwardSeparator,
	}, nil
}

func (l *BulletTwoWayList) buildForwardKey(subject ListSubject, object *ListObject) string {

	var key = l.ListName + l.ForwardSeparator + subject.Value + l.ForwardSeparator //note the separator on the end even if theres no object "a>b>c" or "a>b>"
	if object != nil {
		key = key + *&object.Value
	}
	return key
}

func (l *BulletTwoWayList) buildBackwardKey(subject *ListSubject, object ListObject) string {

	var key = l.ListName + l.ForwardSeparator + object.Value + l.BackwardSeparator //note the separator on the end even if theres no object "a<c<b" or "a<c<"
	if subject != nil {
		key = key + *&subject.Value
	}
	return key
}

func (l *BulletTwoWayList) Upsert(s ListSubject, o ListObject) error {
	//VX:TODO check keySepawrator is not used in names
	forwardKey := l.buildForwardKey(s, &o)
	backwardKey := l.buildBackwardKey(&s, o)

	err := l.TrackStore.TrackInsertOne(l.BucketId, forwardKey, 0, nil, nil)
	if err != nil {
		return err
	}
	err = l.TrackStore.TrackInsertOne(l.BucketId, backwardKey, 0, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// VX:TODO test
func (l *BulletOneWayList) DeleteViaSub(s ListSubject) error {
	//VX:TODO WRITE FOR TWOLIST
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
	//VX:TODO WRITE FOR TWOLIST
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
		return nil, errors.New("missing bucket")
	}

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
	return &OneWayObject{
		Value: object,
	}, nil
}
