package main

import (
	"fmt"

	"elentok.com/dff/disk"
)

func main() {
	disks, err := disk.LoadDisks()
	if err != nil {
		fmt.Println("Failed getting disks: ", err)
		return
	}

	for _, disk := range disks {
		fmt.Printf("Device: %s\t%s free\n", disk.Mount, formatKBs(disk.AvailableKB))
	}
}

func formatKBs(kbs float64) string {
	units := [...]string{"KB", "MB", "GB", "TB"}

	value := kbs

	for _, unit := range units {
		if value < 1024 {
			return fmt.Sprintf("%.2f%s", value, unit)
		}
		value = value / 1024
	}

	return fmt.Sprintf("%.2f%s", value, units[len(units)-1])
}
