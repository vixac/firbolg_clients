package bullet_stl

type BulletIdSpace struct {
	//This represents the length of each bullet id in its aasci form.
	AasciIdSize int //the ID space this maps to equals 36^AasciIdSize
}

type BulletId struct {
	intValue   *int64
	aasciValue *string
}

func (b *BulletId) IntValue() (int64, error) {
	if b.intValue != nil {
		return *b.intValue, nil
	}
	mustExistAasci := *b.aasciValue
	intVal, err := AasciBulletIdToInt(mustExistAasci)
	if err != nil {
		return 0, err
	}
	b.intValue = &intVal
	return intVal, nil
}
