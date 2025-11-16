package bullet_stl

import (
	bullet_client "github.com/vixac/firbolg_clients/bullet/bullet_interface"
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
	GetOSubjectViaObject(o ListObject) (*ListSubject, error)
}

// VX:TODO This is pure implementation. Perhaps does not need to implement TwoWayList interface at all. Can just be a type with methods.
// The bullet client implementation of the OnewayList
type TwoWayListImpl struct {
	ForwardList  OneWayList
	BackwardList OneWayList
}

// just a convenience, becuase TwoWay is bullet agnostic.
func NewBulletTwoWayList(store bullet_client.TrackClientInterface, bucketId int32, listName string, forwardSeparator string, backwardSeparator string) (*TwoWayListImpl, error) {
	forwardList, err := NewBulletOneWayList(store, bucketId, listName, forwardSeparator)
	if err != nil {
		return nil, err
	}
	backwardList, err := NewBulletOneWayList(store, bucketId, listName, backwardSeparator)
	if err != nil {
		return nil, err
	}
	return &TwoWayListImpl{
		ForwardList:  forwardList,
		BackwardList: backwardList,
	}, nil
}

func (l *TwoWayListImpl) Upsert(s ListSubject, o ListObject) error {
	err := l.ForwardList.Upsert(s, o)
	if err != nil {
		return err
	}
	//now we backwards subject and object
	return l.BackwardList.Upsert(o.Invert(), s.Invert())
}

// VX:TODO test
func (l *TwoWayListImpl) DeleteViaSub(s ListSubject) error {
	//VX:TODO WRITE FOR TWOLIST
	o, err := l.ForwardList.GetObject(s)
	if err != nil {
		return err
	}
	err = l.ForwardList.DeleteViaSub(s)
	if err != nil {
		return err
	}
	return l.BackwardList.DeleteViaSub(o.Invert())
}

// VX:TODO test
func (l *TwoWayListImpl) GetObjectViaSubject(s ListSubject) (*ListObject, error) {
	return l.ForwardList.GetObject(s)
}

// VX:TODO test
func (l *TwoWayListImpl) GetOSubjectViaObject(o ListObject) (*ListSubject, error) {
	res, err := l.BackwardList.GetObject(o.Invert())
	if err != nil {
		return nil, err
	}
	inverted := res.Invert()
	return &inverted, nil
}
