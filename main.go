package main

import (
	"fmt"

	bullet "github.com/vixac/firbolg_clients/bullet/wayfinder"
	util "github.com/vixac/firbolg_clients/util"
)

func main() {
	//compilation check for conformance
	var wayfinderClient bullet.WayfinderClientInterface
	fg := util.NewFirbolgClient("url", 213)
	wayfinderClient = &bullet.WayFinderClient{Client: fg}

	fmt.Printf("wayfinder works, %+v", wayfinderClient)

}
