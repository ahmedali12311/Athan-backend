package main

import (
	"io"
	"path/filepath"
	"strconv"
	"time"

	"app/api"
	"app/config"
	"app/schemas"
	"app/translations"

	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var log_writer io.Writer

func main() {
	cfg := config.GetSettings()

	e := echo.New()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.TimestampFieldName = "time"

	logger := zerolog.
		New(log_writer).
		With().
		Timestamp().
		Str("app", cfg.AppCode).
		Caller().
		Logger()

	bundle := translations.Setup(&logger)
	jschemas := schemas.BuildSchemas(&logger)
	app := api.NewAPI(cfg, bundle, jschemas, &logger, false)

	if err := app.Serve(e); err != nil {
		logger.Fatal().Msgf("app.Serve error: %s", err.Error())
	}
}
