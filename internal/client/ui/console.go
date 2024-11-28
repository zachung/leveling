package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"leveling/internal/client/constract"
	"leveling/internal/utils"
)

type Console struct {
	writer *tview.TextView
}

func NewConsole(writer *tview.TextView) *constract.Console {
	console := constract.Console(&Console{writer})

	return &console
}

func (c Console) Info(msg string, args ...any) {
	fmt.Fprintf(c.writer, msg, args...)
}

func (c Console) BattleReport(msg string) {
	fmt.Fprintf(c.writer, "[%.9f] %s\n", utils.NowNanoSeconds(), msg)
}
