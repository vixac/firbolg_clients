package bullet_stl

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	bullet "github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

// / a way for many to many relationships. It's still using subject -> object notation
// / so the subjet namespaces and object namespaces are considered separate, but they
// can contain the same ids for example a->b, a->c, and a->a are fine,
// they have complimentary keys b<-a, c<-a, and a<-a. Mesh takes twice the storage of ForwardMesh for that reason.
type Mesh interface {
	AppendPairs(pairs []ManyToManyPair) error
	RemovePairs(pairs []ManyToManyPair) error
	RemoveSubject(subject ListSubject) error
	RemoveObject(object ListObject) error
	AllPairsForSubject(subject ListSubject) (*PairFetchResponse, error)
	AllPairsForObject(object ListObject) (*PairFetchResponse, error)
}

type BulletMesh struct {
	TrackStore        bullet.TrackClientInterface
	BucketId          int32
	MeshName          string
	ForwardSeparator  string
	BackwardSeparator string
}

func NewBulletMesh(store bullet.TrackClientInterface, bucketId int32, meshName string, forwardSeparator string, backwardSeparator string) (Mesh, error) {
	//VX:TODO check meshName and upward and downward are all valid wrt eachother
	return &BulletMesh{
		TrackStore:        store,
		BucketId:          bucketId,
		MeshName:          meshName,
		ForwardSeparator:  forwardSeparator,
		BackwardSeparator: backwardSeparator,
	}, nil
}

func (b *BulletMesh) AppendPairs(pairs []ManyToManyPair) error {
	//VX:Note can I not bulk insert? oh well.
	for _, pair := range pairs {
		objectValue := pair.Object.Value
		forwardKey := buildKey(b.MeshName, b.ForwardSeparator, pair.Subject.Value, &objectValue)
		backwardKey := buildKey(b.MeshName, b.BackwardSeparator, pair.Object.Value, &pair.Subject.Value)
		floatMetric := float64(pair.Rank)
		err := b.TrackStore.TrackInsertOne(b.BucketId, forwardKey, 0, nil, &floatMetric)
		if err != nil {
			//VX:Note partial fail, some may have inserted.
			return err
		}
		err = b.TrackStore.TrackInsertOne(b.BucketId, backwardKey, 0, nil, &floatMetric)
		if err != nil {
			//VX:Note partial fail, some may have inserted.
			return err
		}
	}
	return nil
}

func (b *BulletMesh) RemoveObject(object ListObject) error {
	pairs, err := b.AllPairsForObject(object)
	if err != nil {
		return err
	}
	return b.RemovePairs(pairs.Pairs)
}

func (b *BulletMesh) RemovePairs(pairs []ManyToManyPair) error {
	var values []bullet.TrackDeleteValue
	for _, pair := range pairs {
		objectValue := pair.Object.Value
		forwardKey := buildKey(b.MeshName, b.ForwardSeparator, pair.Subject.Value, &objectValue)
		values = append(values, bullet.TrackDeleteValue{
			BucketID: b.BucketId,
			Key:      forwardKey,
		})
		backwardKey := buildKey(b.MeshName, b.BackwardSeparator, pair.Object.Value, &pair.Subject.Value)
		values = append(values, bullet.TrackDeleteValue{
			BucketID: b.BucketId,
			Key:      backwardKey,
		})
	}

	req := bullet.TrackDeleteMany{
		Values: values,
	}
	return b.TrackStore.TrackDeleteMany(req)
}

func (b *BulletMesh) RemoveSubject(subject ListSubject) error {
	allPairs, err := b.AllPairsForSubject(subject)
	if err != nil {
		return nil
	}

	return b.RemovePairs(allPairs.Pairs)
}

func (b *BulletMesh) AllPairsForObject(object ListObject) (*PairFetchResponse, error) {
	prefixKey := buildKey(b.MeshName, b.BackwardSeparator, object.Value, nil)
	req := bullet.TrackGetItemsByPrefixRequest{
		BucketID: b.BucketId,
		Prefix:   prefixKey,
	}
	res, err := b.TrackStore.TrackGetManyByPrefix(req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("VX:Fetch res is %+v\n", res)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	if len(res.Values) != 1 {
		return nil, errors.New("missing bucket")
	}
	itemsByBucket := res.Values[b.BucketId]
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
	sort.Strings(itemsInBucket)

	var pairs []ManyToManyPair
	//here the suffix is the subject because we're going backwards.
	for _, itemIncludingPrefix := range itemsInBucket {
		subject, found := strings.CutPrefix(itemIncludingPrefix, prefixKey)
		if !found {
			return nil, errors.New("invalid result did not contain the prefix")
		}
		pairs = append(pairs, ManyToManyPair{
			Subject: ListSubject{Value: subject},
			Object:  object,
		})
	}

	return &PairFetchResponse{
		Pairs: pairs,
	}, nil
}
func (b *BulletMesh) AllPairsForSubject(subject ListSubject) (*PairFetchResponse, error) {
	prefixKey := buildKey(b.MeshName, b.ForwardSeparator, subject.Value, nil)
	req := bullet.TrackGetItemsByPrefixRequest{
		BucketID: b.BucketId,
		Prefix:   prefixKey,
	}
	res, err := b.TrackStore.TrackGetManyByPrefix(req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("VX:Fetch res is %+v\n", res)
	if res == nil {
		return nil, nil
	}
	if len(res.Values) != 1 {
		return nil, errors.New("missing bucket")
	}
	itemsByBucket := res.Values[b.BucketId]
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
	sort.Strings(itemsInBucket)

	var pairs []ManyToManyPair
	for _, itemIncludingPrefix := range itemsInBucket {
		object, found := strings.CutPrefix(itemIncludingPrefix, prefixKey)
		if !found {
			return nil, errors.New("invalid result did not contain the prefix")
		}
		pairs = append(pairs, ManyToManyPair{
			Subject: subject,
			Object:  ListObject{Value: object},
		})
	}

	return &PairFetchResponse{
		Pairs: pairs,
	}, nil
}
