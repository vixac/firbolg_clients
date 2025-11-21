package test_suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
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

func TestOneToManyNamesThatEclipse(t *testing.T) {
	mesh, err := bullet_stl.NewBulletForwardMesh(BuildTestClient(), 42, "test_mesh_prefixes", ">")
	assert.NoError(t, err)

	a := bullet_stl.ListSubject{Value: "a"}
	ab := bullet_stl.ListSubject{Value: "a:b"}
	abc := bullet_stl.ListSubject{Value: "a:b:c"}

	abce := bullet_stl.ListSubject{Value: "a:b:c:e"}

	var pairs []bullet_stl.ManyToManyPair

	//a->z
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: a,
		Object:  bullet_stl.ListObject{Value: "z"},
	})
	//ab->x
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: ab,
		Object:  bullet_stl.ListObject{Value: "x"},
	})

	//abc -> d
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: abc,
		Object:  bullet_stl.ListObject{Value: "d"},
	})
	//abc -> e
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: abc,
		Object:  bullet_stl.ListObject{Value: "e"},
	})

	//abc -> f
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: abc,
		Object:  bullet_stl.ListObject{Value: "f"},
	})

	//abce -> g
	pairs = append(pairs, bullet_stl.ManyToManyPair{
		Subject: abce,
		Object:  bullet_stl.ListObject{Value: "g"},
	})

	err = mesh.AppendPairs(pairs)
	assert.NoError(t, err)
	//fetch a
	foundObjects, err := mesh.AllPairsForSubject(a)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 1)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "z")

	//fetch abc

	foundObjects, err = mesh.AllPairsForSubject(abc)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 3)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "d")
	assert.Equal(t, foundObjects.Pairs[1].Object.Value, "e")
	assert.Equal(t, foundObjects.Pairs[2].Object.Value, "f")

	//now lets fetch with the prefix

	foundObjects, err = mesh.AllPairsForPrefixSubject(ab)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 5)
	//these are sorted by the objects no the subjects
	assert.Equal(t, foundObjects.Pairs[0].Subject.Value, "a:b")
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "x")

	assert.Equal(t, foundObjects.Pairs[1].Subject.Value, "a:b:c")
	assert.Equal(t, foundObjects.Pairs[1].Object.Value, "d")

	assert.Equal(t, foundObjects.Pairs[2].Subject.Value, "a:b:c")
	assert.Equal(t, foundObjects.Pairs[2].Object.Value, "e")

	assert.Equal(t, foundObjects.Pairs[3].Subject.Value, "a:b:c")
	assert.Equal(t, foundObjects.Pairs[3].Object.Value, "f")

	assert.Equal(t, foundObjects.Pairs[4].Subject.Value, "a:b:c:e")
	assert.Equal(t, foundObjects.Pairs[4].Object.Value, "g")

}
