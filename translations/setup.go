package translations

import (
	"bitbucket.org/sadeemTechnology/backend-config"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
)

func Setup(logger *zerolog.Logger) *i18n.Bundle {
	bundle := i18n.NewBundle(language.Arabic)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	arFilePath := config.GetRootPath("active.ar.toml")
	if _, err := bundle.LoadMessageFile(arFilePath); err != nil {
		logger.Fatal().
			Err(err).
			Msgf("i18n: initializing arabic translations: %s", err.Error())
	}
	enFilePath := config.GetRootPath("active.en.toml")
	if _, err := bundle.LoadMessageFile(enFilePath); err != nil {
		logger.Fatal().
			Err(err).
			Msgf("i18n: initializing english translations: %s", err.Error())
	}
	return bundle
}
