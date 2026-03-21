package bullet_stl

import (
	"math"
)

/*
* An int aasci pairing that supports incrementing, and can be instantiated with either the aasci form or the int64 form.
 */
type BulletId struct {
	IntValue   int64
	AasciValue string
}

func (b BulletId) Next() BulletId {
	id, _ := NewBulletIdFromInt(b.IntValue + 1)
	return *id
}

func NewBulletIdFromInt(val int64) (*BulletId, error) {
	str, err := BulletIdIntToAasci(val)
	if err != nil {
		return nil, err
	}
	id := NewBulletIdComplete(str, val)
	return &id, nil

}

// this is basically a wrapper for BulletId
func NewBulletIdFromString(aasci string) (*BulletId, error) {
	intVal, err := AasciBulletIdToInt(aasci)

	if err != nil {
		return nil, err
	}
	return &BulletId{
		AasciValue: aasci,
		IntValue:   intVal,
	}, nil
}

func NewBulletIdComplete(aasci string, intValue int64) BulletId {
	return BulletId{
		AasciValue: aasci,
		IntValue:   intValue,
	}
}

func FitsInInt32(v int64) bool {
	return v >= math.MinInt32 && v <= math.MaxInt32
}
