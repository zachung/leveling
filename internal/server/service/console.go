package service

import (
	log "github.com/sirupsen/logrus"
	"leveling/internal/server/contract"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type Console struct {
}

func NewConsole() *contract.Console {
	console := contract.Console(&Console{})

	return &console
}

func (c Console) Info(msg string, args ...any) {
	log.Infof(msg, args...)
}

func (c Console) Debug(msg string, args ...any) {
	log.Debugf(msg, args...)
}
