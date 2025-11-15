package bullet_stl

import track "github.com/vixac/firbolg_clients/bullet/track"

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
	AllPairsForObject(object ListObject) (PairFetchResponse, error)
	AllPairsForSubject(object ListObject) (PairFetchResponse, error)
}

//I *think* this can be handled with twoWay lists? not sure. not.

type BulletMesh struct {
	TrackStore        track.TrackClientInterface
	BucketId          int32
	MeshName          string
	UpwardSeparator   string
	DownwardSeparator string
}

func NewBulletMesh(store track.TrackClientInterface, bucketId int32, meshName string, upwardSeparator string, downwardSeparator string) (*BulletMesh, error) {
	//VX:TODO check meshName and upward and downward are all valid wrt eachother
	return &BulletMesh{
		TrackStore:        store,
		BucketId:          bucketId,
		MeshName:          meshName,
		UpwardSeparator:   upwardSeparator,
		DownwardSeparator: downwardSeparator,
	}, nil
}
