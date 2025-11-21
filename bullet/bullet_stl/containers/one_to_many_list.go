package bullet_stl

import (
	"errors"
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

type ForwardMesh interface {
	AppendPairs(pairs []ManyToManyPair) error
	RemovePairs(pairs []ManyToManyPair) error
	RemoveSubject(subject ListSubject) error
	AllPairsForSubject(subject ListSubject) (*PairFetchResponse, error)
	AllPairsForPrefixSubject(subject ListSubject) (*PairFetchResponse, error)
}

//I *think* this can be handled with twoWay lists? not sure. not.

type BulletForwardMesh struct {
	TrackStore bullet.TrackClientInterface
	BucketId   int32
	MeshName   string
	Separator  string
}

func NewBulletForwardMesh(store bullet.TrackClientInterface, bucketId int32, meshName string, separator string) (ForwardMesh, error) {
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
		err := b.TrackStore.TrackInsertOne(b.BucketId, key, 0, nil, &floatMetric)
		if err != nil {
			//VX:Note partial fail, some may have inserted.
			return err
		}
	}
	return nil
}

func (b *BulletForwardMesh) RemovePairs(pairs []ManyToManyPair) error {
	var values []bullet.TrackDeleteValue
	for _, pair := range pairs {
		objectValue := pair.Object.Value
		key := buildKey(b.MeshName, b.Separator, pair.Subject.Value, &objectValue, false)
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

func (b *BulletForwardMesh) RemoveSubject(subject ListSubject) error {
	allPairs, err := b.AllPairsForSubject(subject)
	if err != nil {
		return nil
	}
	return b.RemovePairs(allPairs.Pairs)
}
func (b *BulletForwardMesh) AllPairsForPrefixSubject(subject ListSubject) (*PairFetchResponse, error) {
	return b.allPairsForSubjectImpl(subject, true)
}

func (b *BulletForwardMesh) AllPairsForSubject(subject ListSubject) (*PairFetchResponse, error) {
	return b.allPairsForSubjectImpl(subject, false)
}
func (b *BulletForwardMesh) allPairsForSubjectImpl(subject ListSubject, subjectIsActuallyAPrefix bool) (*PairFetchResponse, error) {
	prefixKey := buildKey(b.MeshName, b.Separator, subject.Value, nil, subjectIsActuallyAPrefix)
	req := bullet.TrackGetItemsByPrefixRequest{
		BucketID: b.BucketId,
		Prefix:   prefixKey,
	}
	res, err := b.TrackStore.TrackGetManyByPrefix(req)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("VX:Fetch res is %+v\n", res)
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
