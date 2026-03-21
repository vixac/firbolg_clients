package test_suite

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	bullet_stl "github.com/vixac/firbolg_clients/bullet/bullet_stl/containers"
)

func TestCollectionCreateAndFetchAll(t *testing.T) {
	client := BuildTestClient()
	collection := bullet_stl.NewBulletCollection(42, client, client)

	id1, err := collection.CreateItemUnder("fruit:apple", "red and tasty", nil)
	assert.NoError(t, err)
	assert.NotNil(t, id1)
	assert.Equal(t, "fruit:apple", id1.Key)
	assert.Equal(t, int32(42), id1.Bucket)

	id2, err := collection.CreateItemUnder("fruit:banana", "yellow and sweet", nil)
	assert.NoError(t, err)
	assert.NotNil(t, id2)

	id3, err := collection.CreateItemUnder("veggie:carrot", "orange", nil)
	assert.NoError(t, err)
	assert.NotNil(t, id3)

	all, err := collection.AllItems()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(all))

	byKey := make(map[string]string)
	for k, v := range all {
		byKey[k.Key] = v
	}
	assert.Equal(t, "red and tasty", byKey["fruit:apple"])
	assert.Equal(t, "yellow and sweet", byKey["fruit:banana"])
	assert.Equal(t, "orange", byKey["veggie:carrot"])
}

func TestCollectionFetchByPrefix(t *testing.T) {
	client := BuildTestClient()
	collection := bullet_stl.NewBulletCollection(42, client, client)

	_, err := collection.CreateItemUnder("animal:dog", "woof", nil)
	assert.NoError(t, err)
	_, err = collection.CreateItemUnder("animal:cat", "meow", nil)
	assert.NoError(t, err)
	_, err = collection.CreateItemUnder("vehicle:car", "vroom", nil)
	assert.NoError(t, err)

	animals, err := collection.AllItemsUnderPrefix("animal")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(animals))

	byKey := make(map[string]string)
	for k, v := range animals {
		byKey[k.Key] = v.Payload
	}
	assert.Equal(t, "woof", byKey["animal:dog"])
	assert.Equal(t, "meow", byKey["animal:cat"])

	vehicles, err := collection.AllItemsUnderPrefix("vehicle")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(vehicles))
}

func TestCollectionDeleteItems(t *testing.T) {
	client := BuildTestClient()
	collection := bullet_stl.NewBulletCollection(42, client, client)

	id1, err := collection.CreateItemUnder("item:one", "payload one", nil)
	assert.NoError(t, err)
	assert.NotNil(t, id1)

	id2, err := collection.CreateItemUnder("item:two", "payload two", nil)
	assert.NoError(t, err)
	assert.NotNil(t, id2)

	id3, err := collection.CreateItemUnder("item:three", "payload three", nil)
	assert.NoError(t, err)
	assert.NotNil(t, id3)

	all, err := collection.AllItems()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(all))

	err = collection.DeleteItems([]bullet_stl.CollectionId{*id1, *id2})
	assert.NoError(t, err)

	remaining, err := collection.AllItems()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(remaining))
	for k, v := range remaining {
		assert.Equal(t, "item:three", k.Key)
		assert.Equal(t, "payload three", v)
	}
}

func TestCollectionItemsForKeys(t *testing.T) {
	client := BuildTestClient()
	collection := bullet_stl.NewBulletCollection(42, client, client)

	_, err := collection.CreateItemUnder("animal:dog", "woof", nil)
	assert.NoError(t, err)
	_, err = collection.CreateItemUnder("animal:dog:puppy", "yip", nil)
	assert.NoError(t, err)
	_, err = collection.CreateItemUnder("animal:cat", "meow", nil)
	assert.NoError(t, err)
	_, err = collection.CreateItemUnder("vehicle:car", "vroom", nil)
	assert.NoError(t, err)

	// Exact key lookup — should NOT return "animal:dog:puppy" even though it shares a prefix
	res, err := collection.ItemsForKeys([]string{"animal:dog", "vehicle:car"})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))

	byKey := make(map[string]string)
	for k, v := range res {
		byKey[k.Key] = v.Payload
	}
	assert.Equal(t, "woof", byKey["animal:dog"])
	assert.Equal(t, "vroom", byKey["vehicle:car"])
	assert.Empty(t, byKey["animal:dog:puppy"])
	assert.Empty(t, byKey["animal:cat"])
}

func TestCollectionItemsForKeys_Missing(t *testing.T) {
	client := BuildTestClient()
	collection := bullet_stl.NewBulletCollection(42, client, client)

	_, err := collection.CreateItemUnder("a", "alpha", nil)
	assert.NoError(t, err)

	// One key exists, one does not — should return only the existing one
	res, err := collection.ItemsForKeys([]string{"a", "nonexistent"})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))

	byKey := make(map[string]string)
	for k, v := range res {
		byKey[k.Key] = v.Payload
	}
	assert.Equal(t, "alpha", byKey["a"])
}

func TestCollectionEditPayload(t *testing.T) {
	client := BuildTestClient()
	collection := bullet_stl.NewBulletCollection(42, client, client)

	id, err := collection.CreateItemUnder("config:theme", "light", nil)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	all, err := collection.AllItems()
	assert.NoError(t, err)
	assert.Equal(t, "light", all[*id])

	err = collection.EditPayload(*id, "dark", nil)
	assert.NoError(t, err)

	all, err = collection.AllItems()
	assert.NoError(t, err)
	assert.Equal(t, "dark", all[*id])
}

func TestCollectionEmptyFetch(t *testing.T) {
	client := BuildTestClient()
	collection := bullet_stl.NewBulletCollection(42, client, client)

	all, err := collection.AllItems()
	assert.NoError(t, err)
	assert.Nil(t, all)
}

func TestCollectionUpdatedTime(t *testing.T) {
	client := BuildTestClient()
	collection := bullet_stl.NewBulletCollection(42, client, client)

	// Create without a timestamp — Updated should be zero
	id, err := collection.CreateItemUnder("ts:item", "v1", nil)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	items, err := collection.AllItemsUnderPrefix("ts")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))
	for _, item := range items {
		assert.Equal(t, "v1", item.Payload)
		assert.True(t, item.Updated.IsZero(), "Updated should be zero when no time is provided")
	}

	// Edit with a specific timestamp — Updated should reflect it
	editTime := time.Unix(1700001000, 0)
	err = collection.EditPayload(*id, "v2", &editTime)
	assert.NoError(t, err)

	items, err = collection.AllItemsUnderPrefix("ts")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))
	for _, item := range items {
		assert.Equal(t, "v2", item.Payload)
		assert.Equal(t, editTime.Unix(), item.Updated.Unix())
	}
}
