package utils

import (
	"time"

	"github.com/jedib0t/go-pretty/progress"
)

func NewProgress() progress.Writer {
	pw := progress.NewWriter()
	pw.SetStyle(progress.StyleCircle)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Options.PercentFormat = "%4.1f%%"
	return pw
}
