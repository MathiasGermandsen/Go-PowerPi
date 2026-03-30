package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	parsed, err := zerolog.ParseLevel(level)
	if err != nil {
		parsed = zerolog.InfoLevel
	}

	Log = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"},
	).Level(parsed).With().Timestamp().Caller().Logger()
}
