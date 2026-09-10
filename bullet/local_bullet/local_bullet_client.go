// Package local_bullet provides the legacy constructor name for Bullet's local client.
package local_bullet

import (
	bullet_local "github.com/vixac/bullet/client/local"
	"github.com/vixac/bullet/model"
	"github.com/vixac/bullet/store/store_interface"
)

func NewLocalClient(store store_interface.Store, space model.TenancySpace) *bullet_local.Client {
	return bullet_local.New(store, space)
}
