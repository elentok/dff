package table

import (
	"fmt"

	"github.com/elentok/dff/color"
)

type Cell struct {
	Color      color.Color
	Text       string
	AlignRight bool
}

func PrintTable(cells [][]Cell) {
	cellWidths := calcCellWidths(cells)

	for _, row := range cells {
		for i, cell := range row {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(cell.stringify(cellWidths[i]))
		}
		fmt.Println()
	}
}

func (c Cell) stringify(cellWidth int) string {
	var format string
	if c.AlignRight {
		format = fmt.Sprintf("%%%ds", cellWidth)
	} else {
		format = fmt.Sprintf("%%-%ds", cellWidth)
	}

	formatted := fmt.Sprintf(format, c.Text)
	return c.Color.Paint(formatted)
}

func calcCellWidths(cells [][]Cell) []int {
	if len(cells) == 0 || len(cells[0]) == 0 {
		return []int{}
	}

	cellWidths := make([]int, len(cells[0]))

	for _, row := range cells {
		for i, cell := range row {
			cellWidths[i] = max(cellWidths[i], len(cell.Text))
		}
	}

	return cellWidths
}
