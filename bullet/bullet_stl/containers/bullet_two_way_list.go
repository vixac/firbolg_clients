package bullet_stl

import (
	bullet_client "github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

/*
*A bullet agnostic data structure which allows insertions of subject object (key, value) pairs. It's 1->1, so both sides are considered primary keys.
 */
type TwoWayList interface {
	Upsert(s ListSubject, o ListObject) error
	DeleteViaSub(s ListSubject) error
	DeleteViaObj(o ListObject) error
	GetObjectViaSubject(s ListSubject) (*ListObject, error)
	GetSubjectViaObject(o ListObject) (*ListSubject, error)
	GetSubjectsViaObjectForMany(objects []ListObject) (map[ListObject]*ListSubject, error)
}

// VX:TODO This is pure implementation. Perhaps does not need to implement TwoWayList interface at all. Can just be a type with methods.
// The bullet client implementation of the OnewayList
type TwoWayListImpl struct {
	forwardList  OneWayList
	backwardList OneWayList
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
		forwardList:  forwardList,
		backwardList: backwardList,
	}, nil
}

func (l *TwoWayListImpl) Upsert(s ListSubject, o ListObject) error {
	//this replaces any references to subject or object.
	//these lists dont enforce that the object is unique, so we need to do it manually

	existingForwardSubject, err := l.backwardList.GetObject(o.Invert())
	if err != nil {
		return err
	}
	if existingForwardSubject != nil {
		//the forward list needs to have this object explicitly removed.
		l.forwardList.DeleteBySub(existingForwardSubject.Invert())
	}

	err = l.forwardList.Upsert(s, o)
	if err != nil {
		return err
	}
	//now we backwards subject and object
	return l.backwardList.Upsert(o.Invert(), s.Invert())
}

// VX:TODO test
func (l *TwoWayListImpl) DeleteViaSub(s ListSubject) error {
	//VX:TODO WRITE FOR TWOLIST
	o, err := l.forwardList.GetObject(s)
	if err != nil {
		return err
	}
	err = l.forwardList.DeleteBySub(s)
	if err != nil {
		return err
	}
	if o == nil {
		return nil
	}
	return l.backwardList.DeleteBySub(o.Invert())
}

func (l *TwoWayListImpl) GetObjectViaSubject(s ListSubject) (*ListObject, error) {
	return l.forwardList.GetObject(s)
}

func (l *TwoWayListImpl) GetSubjectsViaObjectForMany(objects []ListObject) (map[ListObject]*ListSubject, error) {
	var inverted []ListSubject
	for _, v := range objects {
		inverted = append(inverted, v.Invert())
	}

	res, err := l.backwardList.GetObjectForMany(inverted)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	resultMap := make(map[ListObject]*ListSubject)
	for key, value := range res {
		if value == nil {
			continue
		}
		newObject := key.Invert()
		newSubject := value.Invert()
		resultMap[newObject] = &newSubject
	}
	return resultMap, nil

}

// VX:TODO test
func (l *TwoWayListImpl) GetSubjectViaObject(o ListObject) (*ListSubject, error) {
	res, err := l.backwardList.GetObject(o.Invert())
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	inverted := res.Invert()
	return &inverted, nil
}
