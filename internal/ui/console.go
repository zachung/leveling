package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"leveling/internal/utils"
)

type Console struct {
	writer *tview.TextView
}

func NewConsole(writer *tview.TextView) *Console {
	return &Console{writer}
}

func (c *Console) Info(msg string) {
	fmt.Fprintln(c.writer, msg)
}

func (c *Console) BattleReport(msg string) {
	fmt.Fprintf(c.writer, "[%.9f] %s\n", utils.NowNanoSeconds(), msg)
}
