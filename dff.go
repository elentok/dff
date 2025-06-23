package main

import (
	"fmt"
	"regexp"

	"elentok.com/dff/color"
	"elentok.com/dff/disk"
	"elentok.com/dff/table"
)

var statusToColor = map[disk.Status]color.Color{
	disk.Good:    color.Green,
	disk.Warning: color.Yellow,
	disk.Error:   color.Red,
}

var excludePattern = regexp.MustCompile("^/(System/Volumes|snap)")

func main() {
	disks, err := disk.LoadDisks(*excludePattern)
	if err != nil {
		fmt.Println("Failed getting disks: ", err)
		return
	}

	data := make([][]table.Cell, len(disks)+1)
	data[0] = []table.Cell{
		{Text: "Mount", Color: color.Blue},
		{Text: "Usage", Color: color.Blue},
		{Text: "Free", Color: color.Blue},
		{Text: "Size", Color: color.Blue},
		{Text: "Device", Color: color.Blue},
	}
	for i, disk := range disks {
		color := statusToColor[disk.Status]
		data[i+1] = []table.Cell{
			{Text: disk.Mount, Color: color},
			{Text: fmt.Sprintf("%d%%", int(disk.UsedPercentage)), AlignRight: true, Color: color},
			{Text: formatKBs(disk.AvailableKB), AlignRight: true, Color: color},
			{Text: formatKBs(disk.Size), AlignRight: true, Color: color},
			{Text: disk.Device, Color: color},
		}
	}

	table.PrintTable(data)
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
