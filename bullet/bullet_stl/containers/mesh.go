package bullet_stl

// / a way for many to many relationships. It's still using subject -> object notation
// / so the subjet namespaces and object namespaces are considered separate, but they
// can contain the same ids for example a->b, a->c, and a->a are fine,
// they have complimentary keys b<-a, c<-a, and a<-a. Mesh takes twice the storage of ForwardMesh for that reason.
type Mesh interface {
	AppendPairs(pairs []ManyToManyPair) error
	RemovePairs(pairs []ManyToManyPair) error
	AllPairsForSubject(subject ListSubject) (*PairFetchResponse, error)
	AllPairsForManySubjects(subject []ListSubject) (*PairFetchResponse, error)
	AllPairsForPrefixSubject(subject ListSubject) (*PairFetchResponse, error)
	AllPairsForObject(object ListObject) (*PairFetchResponse, error)
	AllPairsForManyObjects(objects []ListObject) (*PairFetchResponse, error)
}
