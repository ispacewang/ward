package common

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Str("svc", "docgen").Logger().Level(zerolog.InfoLevel)
}

func SleepMillis(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
