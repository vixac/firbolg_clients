package test_suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vixac/bullet/store/ram"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"github.com/vixac/firbolg_clients/bullet/bullet_stl"
	"github.com/vixac/firbolg_clients/bullet/local_bullet"
)

func buildClient() bullet_interface.BulletClientInterface {
	store := ram.NewRamStore()
	localClient := &local_bullet.LocalBullet{
		Store: store,
		AppId: 12,
	}
	return localClient
}

func TestInsertAndDelete(t *testing.T) {
	list, err := bullet_stl.NewBulletOneWayList(buildClient(), 42, "test_one_way_list", ":")
	assert.NoError(t, err)

	subject := bullet_stl.ListSubject{Value: "subject"}
	err = list.Upsert(subject, bullet_stl.ListObject{Value: "object"})
	assert.NoError(t, err)

	foundObject, err := list.GetObject(subject)
	assert.NoError(t, err)
	assert.Equal(t, foundObject.Value, "object")

	err = list.DeleteBySub(subject)
	assert.NoError(t, err)

	foundObject, err = list.GetObject(subject)
	assert.NoError(t, err)
	assert.Nil(t, foundObject)

}
func TestOneWayListPersonAge(t *testing.T) {
	ageList, err := bullet_stl.NewBulletOneWayList(buildClient(), 42, "test_one_way_list", ":")
	assert.NoError(t, err)

	alice := bullet_stl.ListSubject{Value: "alice"}
	bob := bullet_stl.ListSubject{Value: "bob"}
	carol := bullet_stl.ListSubject{Value: "carol"}

	//no age yet, so nothing to return
	bobsAge, err := ageList.GetObject(bob)
	assert.NoError(t, err)
	assert.Nil(t, bobsAge)

	//insert age 20
	err = ageList.Upsert(bob, bullet_stl.ListObject{Value: "20"})
	assert.NoError(t, err)

	//fetch age, expect 20
	bobsAge, err = ageList.GetObject(bob)
	assert.NoError(t, err)
	assert.Equal(t, bobsAge.Value, "20")

	//now add alice and carol
	err = ageList.Upsert(alice, bullet_stl.ListObject{Value: "30"})
	assert.NoError(t, err)
	err = ageList.Upsert(carol, bullet_stl.ListObject{Value: "40"})
	assert.NoError(t, err)

	//upsert age to 21

	err = ageList.Upsert(bob, bullet_stl.ListObject{Value: "21"})
	assert.NoError(t, err)

	//add alice and carol
	aliceAge, err := ageList.GetObject(alice)
	assert.NoError(t, err)
	assert.Equal(t, aliceAge.Value, "30")

	carolAge, err := ageList.GetObject(carol)
	assert.NoError(t, err)
	assert.Equal(t, carolAge.Value, "40")

	//delete carol
	err = ageList.DeleteBySub(carol)
	assert.NoError(t, err)
	carolAge, err = ageList.GetObject(carol)
	assert.NoError(t, err)
	assert.True(t, carolAge == nil)

	//fetch age, expect 21
	bobsAge, err = ageList.GetObject(bob)
	assert.NoError(t, err)
	assert.Equal(t, bobsAge.Value, "21")

}
