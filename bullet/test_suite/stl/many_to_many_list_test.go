package test_suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
)

func TestManyToManyInsertAndDelete(t *testing.T) {
	mesh, err := bullet_stl.NewBulletMesh(BuildTestClient(), 42, "test_many_mesh", ">", "<")
	assert.NoError(t, err)

	english := bullet_stl.ListSubject{Value: "english"}
	french := bullet_stl.ListSubject{Value: "french"}
	italian := bullet_stl.ListSubject{Value: "italian"}
	var pairs []bullet_stl.ManyToManyPair

	//english speakers: newton, churchill
	//french speakers: churchill, napoleon
	//italian speakers: galileo
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: english,
		Object:  bullet_stl.ListObject{Value: "churchill"},
	})

	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: french,
		Object:  bullet_stl.ListObject{Value: "churchill"},
	})
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: english,
		Object:  bullet_stl.ListObject{Value: "newton"},
	})

	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: french,
		Object:  bullet_stl.ListObject{Value: "napoleon"},
	})

	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: italian,
		Object:  bullet_stl.ListObject{Value: "galileo"},
	})
	err = mesh.AppendPairs(pairs)
	assert.NoError(t, err)

	//fetch uk
	foundObjects, err := mesh.AllPairsForSubject(english, false)

	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 2)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "churchill")
	assert.Equal(t, foundObjects.Pairs[1].Object.Value, "newton")

	//	churchillPair := foundObjects.Pairs[0]
	//fetch france
	foundObjects, err = mesh.AllPairsForSubject(french, false)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 2)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "churchill")
	assert.Equal(t, foundObjects.Pairs[1].Object.Value, "napoleon")

	//delete churchill

	err = mesh.RemoveObject(bullet_stl.ListObject{Value: "churchill"})
	assert.NoError(t, err)

	//find english, only newton
	foundObjects, err = mesh.AllPairsForSubject(english, false)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 1)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "newton")

	//find french, only napoleon
	foundObjects, err = mesh.AllPairsForSubject(french, false)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 1)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "napoleon")

	err = mesh.RemoveSubject(italian)
	assert.NoError(t, err)
	foundObjects, err = mesh.AllPairsForSubject(italian, false)
	assert.NoError(t, err)
	assert.True(t, foundObjects == nil)

	//remove english
	err = mesh.RemoveSubject(english)
	assert.NoError(t, err)

	foundObjects, err = mesh.AllPairsForSubject(english, false)
	assert.NoError(t, err)
	assert.True(t, foundObjects == nil)

}
