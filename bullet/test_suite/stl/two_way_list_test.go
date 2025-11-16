package test_suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vixac/firbolg_clients/bullet/bullet_stl"
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
