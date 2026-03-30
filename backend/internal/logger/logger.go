package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init(debug bool) {
	var output io.Writer = os.Stdout

	if debug {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	Log = zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("service", "zenwoot").
		Logger()

	if debug {
		Log = Log.Level(zerolog.DebugLevel)
	}
}

func Info() *zerolog.Event {
	return Log.Info()
}

func Debug() *zerolog.Event {
	return Log.Debug()
}

func Warn() *zerolog.Event {
	return Log.Warn()
}

func Error() *zerolog.Event {
	return Log.Error()
}

func Fatal() *zerolog.Event {
	return Log.Fatal()
}

func With() zerolog.Context {
	return Log.With()
}
