package bullet_stl

import (
	"errors"
	"sort"
	"strings"

	"github.com/vixac/bullet/client"
	"github.com/vixac/bullet/model"
)

type PairFetchResponse struct {
	Pairs []ManyToManyPair
}

type ForwardMesh interface {
	AppendPairs(pairs []ManyToManyPair) error
	RemovePairs(pairs []ManyToManyPair) error
	AllPairsForSubject(subject ListSubject) (*PairFetchResponse, error)
	AllPairsForPrefixSubject(subject ListSubject) (*PairFetchResponse, error)
}

//I *think* this can be handled with twoWay lists? not sure. not.

type BulletForwardMesh struct {
	TrackStore client.Track
	BucketId   int32
	MeshName   string
	Separator  string
}

func NewBulletForwardMesh(store client.Track, bucketId int32, meshName string, separator string) (ForwardMesh, error) {
	//VX:TODO check meshName and upward and downward are all valid wrt eachother
	return &BulletForwardMesh{
		TrackStore: store,
		BucketId:   bucketId,
		MeshName:   meshName,
		Separator:  separator,
	}, nil
}

func (b *BulletForwardMesh) AppendPairs(pairs []ManyToManyPair) error {

	//VX:Note can I not bulk insert? oh well.
	for _, pair := range pairs {
		objectValue := pair.Object.Value
		key := buildKey(b.MeshName, b.Separator, pair.Subject.Value, &objectValue, false)
		floatMetric := float64(pair.Rank)
		err := b.TrackStore.TrackPut(b.BucketId, key, 0, nil, &floatMetric)
		if err != nil {
			//VX:Note partial fail, some may have inserted.
			return err
		}
	}
	return nil
}

func (b *BulletForwardMesh) RemovePairs(pairs []ManyToManyPair) error {
	var values []model.TrackKey
	for _, pair := range pairs {
		objectValue := pair.Object.Value
		key := buildKey(b.MeshName, b.Separator, pair.Subject.Value, &objectValue, false)
		values = append(values, model.TrackKey{
			BucketID: b.BucketId,
			Key:      key,
		})
	}

	return b.TrackStore.TrackDeleteMany(values)
}

func (b *BulletForwardMesh) AllPairsForPrefixSubject(subject ListSubject) (*PairFetchResponse, error) {
	return b.allPairsForSubjectImpl(subject, true)
}

func (b *BulletForwardMesh) AllPairsForSubject(subject ListSubject) (*PairFetchResponse, error) {
	return b.allPairsForSubjectImpl(subject, false)
}

func (b *BulletForwardMesh) allPairsForSubjectImpl(subject ListSubject, subjectIsActuallyAPrefix bool) (*PairFetchResponse, error) {
	prefixKey := buildKey(b.MeshName, b.Separator, subject.Value, nil, subjectIsActuallyAPrefix)
	items, err := b.TrackStore.GetItemsByKeyPrefix(b.BucketId, prefixKey, nil, nil, false)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	itemsInBucket := make([]string, 0, len(items))
	for _, item := range items {
		itemsInBucket = append(itemsInBucket, item.Key)
	}
	sort.Strings(itemsInBucket)

	var pairs []ManyToManyPair
	for _, itemIncludingPrefix := range itemsInBucket {

		split := strings.Split(itemIncludingPrefix, b.Separator)
		if len(split) != 3 {
			return nil, errors.New("expected <listname><separator><subject><separator><object")
		}

		subjectValue := split[1]
		objectValue := split[2]

		pairs = append(pairs, ManyToManyPair{
			Subject: ListSubject{Value: subjectValue},
			Object:  ListObject{Value: objectValue},
		})
	}

	sort.Slice(pairs, func(i, j int) bool {
		a := pairs[i]
		b := pairs[j]

		if a.Subject.Value == b.Subject.Value {
			return a.Object.Value < b.Object.Value
		} else {
			return a.Subject.Value < b.Subject.Value
		}
	})
	return &PairFetchResponse{
		Pairs: pairs,
	}, nil
}
