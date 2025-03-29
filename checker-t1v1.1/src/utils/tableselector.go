package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TableSelector(itemCount int, table *tview.Table) func(event *tcell.EventKey) *tcell.EventKey {
	// Used for closure
	selectedRow := 0
	selectedCol := 0
	return func(event *tcell.EventKey) *tcell.EventKey {

		sCell := table.GetCell(selectedRow, selectedCol)

		_, backColor, _ := sCell.Style.Decompose()

		// Reverse colors
		sCell.SetTextColor(backColor)
		sCell.SetBackgroundColor(tcell.ColorDefault)

		// Get max column for current column
		fillDiff := table.GetColumnCount()*table.GetRowCount() - itemCount

		maxCol := table.GetColumnCount()
		if selectedRow+1 > table.GetRowCount()-fillDiff {
			maxCol = table.GetColumnCount() - 1
		}
		maxRow := table.GetRowCount()
		if itemCount-(selectedCol+1)*table.GetRowCount() < 0 {
			maxRow = itemCount % table.GetRowCount()
		}

		switch event.Key() {
		case tcell.KeyRight:
			selectedCol++
			if selectedCol == maxCol {
				selectedCol = 0
			}
		case tcell.KeyLeft:
			selectedCol--
			if selectedCol < 0 {
				selectedCol = maxCol - 1
			}
		case tcell.KeyUp:
			selectedRow--
			if selectedRow < 0 {
				selectedRow = maxRow - 1
			}
		case tcell.KeyDown:
			selectedRow++
			if selectedRow == maxRow {
				selectedRow = 0
			}
		case tcell.KeyEnter:
			sCell := table.GetCell(selectedRow, selectedCol)
			if sCell != nil {
				sCell.Clicked()
			}
		default:
			return event

		}

		sCell = table.GetCell(selectedRow, selectedCol)

		textColor, _, _ := sCell.Style.Decompose()

		// Reverse colors
		sCell.SetBackgroundColor(textColor)
		sCell.SetTextColor(tcell.ColorWhite)

		return nil
	}
}
