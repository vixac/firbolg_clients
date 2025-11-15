package localbullet

import (
	store "github.com/vixac/bullet/store/store_interface"
)

/*
Local bullet is an implemtnation of BulletClientInterface that uses a local implemetnation of store, so no network calls required.
*/
type LocalBullet struct {
	store store.Store
	appId int32
}
