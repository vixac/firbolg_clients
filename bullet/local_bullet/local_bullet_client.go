package local_bullet

import (
	store "github.com/vixac/bullet/store/store_interface"
)

/*
Local bullet is an implemtnation of BulletClientInterface that uses a local implemetnation of store, so no network calls required.
WayFinder is no longer part of Store in v0.2.0; provide a WayFinderStore implementation separately if needed.
*/
type LocalBullet struct {
	Store      store.Store
	WayFinder  store.WayFinderStore
	Space      store.TenancySpace
}
