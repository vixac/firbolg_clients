package test_suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
)

func TestTwoWayInsertAndDelete(t *testing.T) {
	footballerNumbers, err := bullet_stl.NewBulletTwoWayList(BuildTestClient(), 42, "test_two_way_list", ">", "<")
	assert.NoError(t, err)

	alice := bullet_stl.ListSubject{Value: "alice"}
	bob := bullet_stl.ListSubject{Value: "bob"}
	carol := bullet_stl.ListSubject{Value: "carol"}

	//insert alice as 9
	footballerNumbers.Upsert(alice, bullet_stl.ListObject{Value: "9"})
	aliceNumber, err := footballerNumbers.GetObjectViaSubject(alice)
	assert.NoError(t, err)
	assert.Equal(t, aliceNumber.Value, "9")

	//insert bob as 10
	footballerNumbers.Upsert(bob, bullet_stl.ListObject{Value: "10"})
	bobNumber, err := footballerNumbers.GetObjectViaSubject(bob)
	assert.NoError(t, err)
	assert.Equal(t, bobNumber.Value, "10")

	//insert carol as 9, which will remove alice entirely

	footballerNumbers.Upsert(carol, bullet_stl.ListObject{Value: "9"})
	carolNumber, err := footballerNumbers.GetObjectViaSubject(carol)
	assert.NoError(t, err)
	assert.Equal(t, carolNumber.Value, "9")

	//alice was removed by carols insertion as #9
	aliceNumber, err = footballerNumbers.GetObjectViaSubject(alice)
	assert.NoError(t, err)
	assert.True(t, aliceNumber == nil)

	carolSubject, err := footballerNumbers.GetOSubjectViaObject(bullet_stl.ListObject{Value: "9"})
	assert.NoError(t, err)
	assert.Equal(t, carolSubject.Value, "carol")

	//reinsert bob as 9, and confirm that carol is deleted.
	footballerNumbers.Upsert(bob, bullet_stl.ListObject{Value: "9"})
	bobNumber, err = footballerNumbers.GetObjectViaSubject(bob)
	assert.NoError(t, err)
	assert.Equal(t, bobNumber.Value, "9")
	bobSubject, err := footballerNumbers.GetOSubjectViaObject(bullet_stl.ListObject{Value: "9"})
	assert.NoError(t, err)
	assert.Equal(t, bobSubject.Value, "bob")

	carolNumber, err = footballerNumbers.GetObjectViaSubject(carol)
	assert.NoError(t, err)
	assert.True(t, carolNumber == nil)
}

func TestTwoWayNamesThatEclipse(t *testing.T) {
	names, err := bullet_stl.NewBulletTwoWayList(BuildTestClient(), 42, "test_two_way", ">", "<")
	assert.NoError(t, err)

	a := bullet_stl.ListSubject{Value: "a"}
	ab := bullet_stl.ListSubject{Value: "a:b"}
	abc := bullet_stl.ListSubject{Value: "a:b:c"}
	ad := bullet_stl.ListSubject{Value: "a:d"}
	e := bullet_stl.ListSubject{Value: "e"}

	var pairs []bullet_stl.ManyToManyPair
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: a,
		Object:  bullet_stl.ListObject{Value: "ant"},
	})
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: ab,
		Object:  bullet_stl.ListObject{Value: "abi"},
	})

	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: abc,
		Object:  bullet_stl.ListObject{Value: "abacus"},
	})

	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: ad,
		Object:  bullet_stl.ListObject{Value: "advark"},
	})
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: e,
		Object:  bullet_stl.ListObject{Value: "elephant"},
	})

	//doesnt need to be pairs actually but it works
	for _, p := range pairs {
		err := names.Upsert(p.Subject, p.Object)
		assert.NoError(t, err)
	}

	//fetch a
	foundObject, err := names.GetObjectViaSubject(a)

	assert.NoError(t, err)
	assert.Equal(t, foundObject.Value, "ant")

	foundSubject, err := names.GetOSubjectViaObject(bullet_stl.ListObject{Value: "ant"})
	assert.NoError(t, err)
	assert.Equal(t, foundSubject.Value, "a")

	//fetch a:b
	foundObject, err = names.GetObjectViaSubject(ab)

	assert.NoError(t, err)
	assert.Equal(t, foundObject.Value, "abi")

	foundSubject, err = names.GetOSubjectViaObject(bullet_stl.ListObject{Value: "abi"})
	assert.NoError(t, err)
	assert.Equal(t, foundSubject.Value, "a:b")

}
