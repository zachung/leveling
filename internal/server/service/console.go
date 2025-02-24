package service

import (
	log "github.com/sirupsen/logrus"
	"leveling/internal/server/contract"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

type Console struct {
}

func NewConsole() contract.Console {
	return new(Console)
}

func (c Console) Info(msg string, args ...any) {
	log.Infof(msg, args...)
}

func (c Console) Debug(msg string, args ...any) {
	log.Debugf(msg, args...)
}
