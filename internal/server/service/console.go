package service

import (
	"fmt"
	"leveling/internal/server/constract"
)

type Console struct {
}

func NewConsole() *constract.Console {
	console := constract.Console(&Console{})

	return &console
}

func (c Console) Info(msg string, args ...any) {
	fmt.Printf(msg, args...)
}
