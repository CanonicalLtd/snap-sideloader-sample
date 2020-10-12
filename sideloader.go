package main

import (
	"github.com/slimjim777/snap-sideloader/service"
	"github.com/slimjim777/snap-sideloader/snapd"
)

func main() {
	client := snapd.NewClient("/mnt")
	srv := service.NewWebService("127.0.0.1", "5000", client)
	srv.Start()
}
