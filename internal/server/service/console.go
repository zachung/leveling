package service

import (
	"leveling/internal/server/contract"
	"log"
)

type Console struct {
}

func NewConsole() *contract.Console {
	console := contract.Console(&Console{})

	return &console
}

func (c Console) Info(msg string, args ...any) {
	log.Printf(msg, args...)
}
