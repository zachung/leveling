package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"leveling/internal/client/contract"
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
