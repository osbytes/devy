package infra

import (
	"bot/pkg/env"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger = log.Logger

func InitLogging(logLevel string) error {
	level, err := zerolog.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(level)

	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	Logger = Logger.Output(os.Stdout)
	if env.IsLocal() {
		Logger = Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return nil
}
