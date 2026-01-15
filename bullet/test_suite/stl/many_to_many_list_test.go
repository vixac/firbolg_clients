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
	foundObjects, err := mesh.AllPairsForSubject(english)

	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 2)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "churchill")
	assert.Equal(t, foundObjects.Pairs[1].Object.Value, "newton")

	//fetch churchills languages
	churchillsLanguages, err := mesh.AllPairsForObject(bullet_stl.ListObject{Value: "churchill"})
	assert.NoError(t, err)
	assert.Equal(t, len(churchillsLanguages.Pairs), 2)
	assert.Equal(t, churchillsLanguages.Pairs[0].Subject.Value, "english")
	assert.Equal(t, churchillsLanguages.Pairs[0].Object.Value, "churchill")
	assert.Equal(t, churchillsLanguages.Pairs[1].Subject.Value, "french")
	assert.Equal(t, churchillsLanguages.Pairs[1].Object.Value, "churchill")

	//fetch france
	foundObjects, err = mesh.AllPairsForSubject(french)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 2)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "churchill")
	assert.Equal(t, foundObjects.Pairs[1].Object.Value, "napoleon")

	//delete churchill

	err = mesh.RemoveObject(bullet_stl.ListObject{Value: "churchill"})
	assert.NoError(t, err)

	//find english, only newton
	foundObjects, err = mesh.AllPairsForSubject(english)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 1)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "newton")

	//find french, only napoleon
	foundObjects, err = mesh.AllPairsForSubject(french)
	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 1)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "napoleon")

	err = mesh.RemoveSubject(italian)
	assert.NoError(t, err)
	foundObjects, err = mesh.AllPairsForSubject(italian)
	assert.NoError(t, err)
	assert.True(t, foundObjects == nil)

	//remove english
	err = mesh.RemoveSubject(english)
	assert.NoError(t, err)

	foundObjects, err = mesh.AllPairsForSubject(english)
	assert.NoError(t, err)
	assert.True(t, foundObjects == nil)
}

func TestManyToManyNamesThatEclipse(t *testing.T) {
	mesh, err := bullet_stl.NewBulletMesh(BuildTestClient(), 42, "test_many_mesh", ">", "<")
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

func TestAllPairsForManySubjects_Basic(t *testing.T) {
	mesh, err := bullet_stl.NewBulletMesh(BuildTestClient(), 42, "test_many_subjects", ">", "<")
	assert.NoError(t, err)

	english := bullet_stl.ListSubject{Value: "english"}
	french := bullet_stl.ListSubject{Value: "french"}
	italian := bullet_stl.ListSubject{Value: "italian"}

	pairs := []bullet_stl.ManyToManyPair{
		{Subject: english, Object: bullet_stl.ListObject{Value: "newton"}},
		{Subject: english, Object: bullet_stl.ListObject{Value: "churchill"}},
		{Subject: french, Object: bullet_stl.ListObject{Value: "napoleon"}},
		{Subject: italian, Object: bullet_stl.ListObject{Value: "galileo"}},
	}

	err = mesh.AppendPairs(pairs)
	assert.NoError(t, err)

	foundObjects, err := mesh.AllPairsForSubject(french)

	assert.NoError(t, err)
	assert.Equal(t, len(foundObjects.Pairs), 1)
	assert.Equal(t, foundObjects.Pairs[0].Object.Value, "napoleon")

	res, err := mesh.AllPairsForManySubjects([]bullet_stl.ListSubject{
		english,
		italian,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 3, len(res.Pairs))

	// Sorted by subject, then object
	assert.Equal(t, "english", res.Pairs[0].Subject.Value)
	assert.Equal(t, "churchill", res.Pairs[0].Object.Value)

	assert.Equal(t, "english", res.Pairs[1].Subject.Value)
	assert.Equal(t, "newton", res.Pairs[1].Object.Value)

	assert.Equal(t, "italian", res.Pairs[2].Subject.Value)
	assert.Equal(t, "galileo", res.Pairs[2].Object.Value)
}

func TestAllPairsForManyObjects_Basic(t *testing.T) {
	mesh, err := bullet_stl.NewBulletMesh(BuildTestClient(), 42, "test_many_objects", ">", "<")
	assert.NoError(t, err)

	english := bullet_stl.ListSubject{Value: "english"}
	french := bullet_stl.ListSubject{Value: "french"}
	italian := bullet_stl.ListSubject{Value: "italian"}

	churchill := bullet_stl.ListObject{Value: "churchill"}
	napoleon := bullet_stl.ListObject{Value: "napoleon"}
	galileo := bullet_stl.ListObject{Value: "galileo"}

	pairs := []bullet_stl.ManyToManyPair{
		{Subject: english, Object: churchill},
		{Subject: french, Object: churchill},
		{Subject: french, Object: napoleon},
		{Subject: italian, Object: galileo},
	}

	err = mesh.AppendPairs(pairs)
	assert.NoError(t, err)

	// Fetch pairs for multiple objects
	res, err := mesh.AllPairsForManyObjects([]bullet_stl.ListObject{
		churchill,
		galileo,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 3, len(res.Pairs))

	// Sorted by subject, then object
	assert.Equal(t, "english", res.Pairs[0].Subject.Value)
	assert.Equal(t, "churchill", res.Pairs[0].Object.Value)

	assert.Equal(t, "french", res.Pairs[1].Subject.Value)
	assert.Equal(t, "churchill", res.Pairs[1].Object.Value)

	assert.Equal(t, "italian", res.Pairs[2].Subject.Value)
	assert.Equal(t, "galileo", res.Pairs[2].Object.Value)
}

func TestAllPairsForManyObjects_EmptyResult(t *testing.T) {
	mesh, err := bullet_stl.NewBulletMesh(BuildTestClient(), 42, "test_many_objects_empty", ">", "<")
	assert.NoError(t, err)

	// Query for objects that don't exist
	res, err := mesh.AllPairsForManyObjects([]bullet_stl.ListObject{
		{Value: "nonexistent1"},
		{Value: "nonexistent2"},
	})
	assert.NoError(t, err)
	assert.Nil(t, res)
}
