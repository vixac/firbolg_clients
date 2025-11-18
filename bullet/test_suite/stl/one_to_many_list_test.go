package test_suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vixac/firbolg_clients/bullet/bullet_stl"
)

func TestOneToManyInsertAndDelete(t *testing.T) {
	mesh, err := bullet_stl.NewBulletForwardMesh(BuildTestClient(), 42, "test_mesh", ":")
	assert.NoError(t, err)

	uk := bullet_stl.ListSubject{Value: "uk"}
	france := bullet_stl.ListSubject{Value: "france"}
	italy := bullet_stl.ListSubject{Value: "germany"}
	var pairs []bullet_stl.ManyToManyPair
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: uk,
		Object:  bullet_stl.ListObject{Value: "churchill"},
	})
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: uk,
		Object:  bullet_stl.ListObject{Value: "newton"},
	})

	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: france,
		Object:  bullet_stl.ListObject{Value: "napoleon"},
	})

	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: italy,
		Object:  bullet_stl.ListObject{Value: "galileo"},
	})
	err = mesh.AppendPairs(pairs)
	assert.NoError(t, err)

	//fetch uk
	foundObjects, err := mesh.AllPairsForSubject(uk)

	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 2)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "churchill")
	assert.Equal(t, foundObjects.Pairs[1].Object.Value, "newton")

	newtonPair := foundObjects.Pairs[1]
	//fetch france
	foundObjects, err = mesh.AllPairsForSubject(france)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 1)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "napoleon")

	//delete newton

	var deletePairs []bullet_stl.ManyToManyPair
	deletePairs = append(deletePairs, newtonPair)
	err = mesh.RemovePairs(deletePairs)
	assert.NoError(t, err)

	//find uk, only churchill
	foundObjects, err = mesh.AllPairsForSubject(uk)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 1)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "churchill")

	err = mesh.RemoveSubject(italy)
	assert.NoError(t, err)
	foundObjects, err = mesh.AllPairsForSubject(italy)
	assert.NoError(t, err)
	assert.True(t, foundObjects == nil)

}
