package api

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/nwidger/jsoncolor"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func GetLogger(env, appCode string) *zerolog.Logger {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.TimestampFieldName = "time"

	var logger zerolog.Logger

	if env != "production" {
		zerolog.InterfaceMarshalFunc = func(v any) ([]byte, error) {
			f := jsoncolor.NewFormatter()

			// set custom colors
			f.SpaceColor = color.New(color.FgRed, color.Bold)
			f.CommaColor = color.New(color.FgWhite, color.Bold)
			f.ColonColor = color.New(color.FgYellow, color.Bold)
			f.ObjectColor = color.New(color.FgGreen, color.Bold)
			f.ArrayColor = color.New(color.FgHiRed)
			f.FieldColor = color.New(color.FgCyan)
			f.StringColor = color.New(color.FgHiYellow)
			f.TrueColor = color.New(color.FgCyan, color.Bold)
			f.FalseColor = color.New(color.FgHiRed)
			f.NumberColor = color.New(color.FgHiMagenta)
			f.NullColor = color.New(color.FgWhite, color.Bold)
			f.StringQuoteColor = color.New(color.FgBlue, color.Bold)

			return jsoncolor.MarshalIndentWithFormatter(v, "", "  ", f)
		}
		logger = zerolog.
			New(zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}).
			With().
			Timestamp().
			Str("app", appCode).
			Caller().
			Logger()
	} else {
		logger = zerolog.
			New(os.Stdout).
			With().
			Timestamp().
			Str("app", appCode).
			Caller().
			Logger()
	}

	log.Logger = logger
	return &logger
}
