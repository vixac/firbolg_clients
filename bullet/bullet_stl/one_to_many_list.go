package bullet_stl

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	bullet "github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

type ManyToManyPair struct {
	Subject ListSubject //the item above
	Object  ListObject  //the item below
	Rank    int32       //just metadata for the user
}

type PairFetchResponse struct {
	Pairs []ManyToManyPair
}

type Mesh interface {
	AppendPairs(pairs []ManyToManyPair) error
	RemovePairs(pairs []ManyToManyPair) error
	RemoveSubject(subject ListSubject) error
	AllPairsForSubject(subject ListSubject) (*PairFetchResponse, error)
}

//I *think* this can be handled with twoWay lists? not sure. not.

type BulletMesh struct {
	TrackStore bullet.TrackClientInterface
	BucketId   int32
	MeshName   string
	Separator  string
}

func NewBulletMesh(store bullet.TrackClientInterface, bucketId int32, meshName string, separator string) (*BulletMesh, error) {
	//VX:TODO check meshName and upward and downward are all valid wrt eachother
	return &BulletMesh{
		TrackStore: store,
		BucketId:   bucketId,
		MeshName:   meshName,
		Separator:  separator,
	}, nil
}

func (b *BulletMesh) AppendPairs(pairs []ManyToManyPair) error {

	//VX:Note can I not bulk insert? oh well.
	for _, pair := range pairs {
		objectValue := pair.Object.Value
		key := buildKey(b.MeshName, b.Separator, pair.Subject.Value, &objectValue)
		floatMetric := float64(pair.Rank)
		err := b.TrackStore.TrackInsertOne(b.BucketId, key, 0, nil, &floatMetric)
		if err != nil {
			//VX:Note partial fail, some may have inserted.
			return err
		}
	}
	return nil
}
func (b *BulletMesh) RemovePairs(pairs []ManyToManyPair) error {
	var values []bullet.TrackDeleteValue
	for _, pair := range pairs {
		objectValue := pair.Object.Value
		key := buildKey(b.MeshName, b.Separator, pair.Subject.Value, &objectValue)
		values = append(values, bullet.TrackDeleteValue{
			BucketID: b.BucketId,
			Key:      key,
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

func (b *BulletMesh) AllPairsForSubject(subject ListSubject) (*PairFetchResponse, error) {
	prefixKey := buildKey(b.MeshName, b.Separator, subject.Value, nil)
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
