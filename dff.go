package main

import (
	"fmt"
	"regexp"

	"github.com/elentok/dff/disk"
	"github.com/elentok/dff/render"
)

var excludePattern = regexp.MustCompile("^/(System/Volumes|snap)")

func main() {
	disks, err := disk.LoadDisks(*excludePattern)
	if err != nil {
		fmt.Println("Failed getting disks: ", err)
		return
	}

	render.RenderTable(disks)
}
