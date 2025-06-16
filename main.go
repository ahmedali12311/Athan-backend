package main

import (
	"app/api"
	"app/schemas"
	"app/translations"
	"bitbucket.org/sadeemTechnology/backend-config"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	cfg := config.GetSettings()

	logger := api.GetLogger(cfg.Env, cfg.AppCode)

	bundle := translations.Setup(logger)
	jschemas := schemas.BuildSchemas(logger)
	app := api.NewAPI(cfg, bundle, jschemas, logger, false)

	if err := app.Serve(e); err != nil {
		logger.Fatal().Msgf("app.Serve error: %s", err.Error())
	}
}
