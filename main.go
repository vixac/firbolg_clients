package main

import (
	"fmt"

	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	rest_bullet "github.com/vixac/firbolg_clients/bullet/rest_bullet/wayfinder"
)

func main() {
	//compilation check for conformance
	var wayfinderClient bullet_interface.WayFinderClientInterface
	wayfinderClient = rest_bullet.NewWayFinderClient(
		"baseUrl",
		12345,
	)

	fmt.Printf("wayfinder works, %+v", wayfinderClient)

}
