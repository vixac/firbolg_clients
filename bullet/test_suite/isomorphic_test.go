package test_suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ram "github.com/vixac/bullet/store/ram"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	local_bullet "github.com/vixac/firbolg_clients/bullet/local_bullet"
)

// the goal here is to test that both clients behave in the same way.
// The problem is I don't have a complete rest client setup, as each
// rest client test just sets up the 1 endpoint being tested.
// this can be added later.

func buildClients() []bullet_interface.BulletClientInterface {
	store := ram.NewRamStore()
	localClient := &local_bullet.LocalBullet{
		Store: store,
		AppId: 12,
	}
	var clients []bullet_interface.BulletClientInterface
	clients = append(clients, localClient)
	return clients
	//VX:TODO add rest client in here, and make this a map
}
func TestSomething(t *testing.T) {
	clients := buildClients()
	for _, c := range clients {
		err := c.TrackInsertOne(1, "testKey", int64(1234), nil, nil)
		assert.NoError(t, err)
	}

}
