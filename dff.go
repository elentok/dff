package main

import (
	"fmt"
	"os"

	"elentok.com/dff/disk"
	"github.com/olekukonko/tablewriter"
)

func main() {
	disks, err := disk.LoadDisks()
	if err != nil {
		fmt.Println("Failed getting disks: ", err)
		return
	}

	var data [][]string
	for _, disk := range disks {
		row := []string{disk.Device, formatKBs(disk.AvailableKB)}
		data = append(data, row)
	}

	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"Device", "Free"}
	table.Header(header)
	table.Bulk(data)
	table.Render()
	//
	// for _, disk := range disks {
	// 	fmt.Printf("Device: %s\t%s free\n", disk.Mount, formatKBs(disk.AvailableKB))
	// }
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
