package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"leveling/internal/client/contract"
	"leveling/internal/utils"
)

type Console struct {
	writer *tview.TextView
}

func NewConsole(writer *tview.TextView) *contract.Console {
	console := contract.Console(&Console{writer})

	return &console
}

func (c Console) Info(msg string, args ...any) {
	fmt.Fprintf(c.writer, msg, args...)
}

func (c Console) BattleReport(msg string) {
	fmt.Fprintf(c.writer, "[%.9f] %s\n", utils.NowNanoSeconds(), msg)
}
