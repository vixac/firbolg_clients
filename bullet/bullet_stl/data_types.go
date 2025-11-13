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
