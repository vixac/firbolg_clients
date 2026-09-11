package bullet_stl

type ListSubject struct {
	Value string
}

func (s ListSubject) Invert() ListObject {
	return ListObject{Value: s.Value}
}
func (o ListObject) Invert() ListSubject {
	return ListSubject{Value: o.Value}
}

type ListObject struct {
	Value string
}

type ManyToManyPair struct {
	Subject ListSubject //the item above
	Object  ListObject  //the item below
	Rank    int32       //just metadata for the user
}
