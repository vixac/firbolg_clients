package bullet

import (
	depot "github.com/vixac/firbolg_clients/bullet/depot"
	track "github.com/vixac/firbolg_clients/bullet/track"
	wayfinder "github.com/vixac/firbolg_clients/bullet/wayfinder"
)

type BulletClientInterface interface {
	track.TrackClientInterface
	depot.DepotClientInterface
	wayfinder.WayFinderClientInterface
}
