package utils

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

func NewTable() table.Writer {
	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	tbl.SetStyle(table.StyleColoredMagentaWhiteOnBlack)
	tbl.SetRowPainter(func(row table.Row) text.Colors {
		if row[len(row)-1] == SyncStatusPossiblySynced {
			return text.Colors{text.BgHiMagenta, text.FgHiBlack}
		}

		if row[len(row)-1] == SyncStatusSynced {
			return text.Colors{text.BgHiGreen, text.FgHiWhite}
		}

		if row[len(row)-1] == SyncStatusGoodToGo {
			return text.Colors{text.BgBlack, text.FgHiWhite}
		}

		return text.Colors{text.BgRed, text.FgHiWhite}
	})
	return tbl
}
