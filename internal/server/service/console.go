package service

import (
	"fmt"
	"leveling/internal/server/contract"
)

type Console struct {
}

func NewConsole() *contract.Console {
	console := contract.Console(&Console{})

	return &console
}

func (c Console) Info(msg string, args ...any) {
	fmt.Printf(msg, args...)
}
