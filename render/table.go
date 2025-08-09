package render

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/elentok/dff/disk"
)

const (
	Yellow = lipgloss.Color("3")
	Green  = lipgloss.Color("2")
	Red    = lipgloss.Color("1")
	Blue   = lipgloss.Color("4")
)

var cellStyle = lipgloss.NewStyle().Padding(0, 1)

var goodStyle = cellStyle.Foreground(Green)
var warningStyle = cellStyle.Foreground(Yellow)
var errorStyle = cellStyle.Foreground(Red)
var headerStyle = cellStyle.Foreground(Blue)

var statusToStyle = map[disk.Status]lipgloss.Style{
	disk.Good:    goodStyle,
	disk.Warning: warningStyle,
	disk.Error:   errorStyle,
}

func RenderTable(disks []disk.Disk) {
	t := table.
		New().
		Border(lipgloss.RoundedBorder()).
		Headers("Mount", "Usage", "Free", "Size", "Device").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}

			style := statusToStyle[disks[row].Status]
			if col >= 1 && col <= 3 {
				return style.Align(lipgloss.Right)
			}

			return style
		})

	for _, d := range disks {
		usage := fmt.Sprintf("%d%%", int(d.UsedPercentage))
		free := formatKBs(d.AvailableKB)
		size := formatKBs(d.Size)

		t.Row(d.Mount, usage, free, size, d.Device)
	}

	fmt.Println(t)
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
