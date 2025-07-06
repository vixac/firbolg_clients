package main

import (
	"fmt"

	bullet "github.com/vixac/firbolg_clients/bullet/wayfinder"
	wayfinder "github.com/vixac/firbolg_clients/bullet/wayfinder"
)

func main() {
	//compilation check for conformance
	var wayfinderClient bullet.WayFinderClientInterface
	wayfinderClient = wayfinder.NewWayFinderClient(
		"baseUrl",
		12345,
	)

	fmt.Printf("wayfinder works, %+v", wayfinderClient)

}
