package scheduler

import (
	"app/models"
	"app/translations"
	"app/utilities"
	"time"

	config "bitbucket.org/sadeemTechnology/backend-config"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog"
)

type Config struct {
	DB     *sqlx.DB
	FB     *firebase.App
	FBM    *messaging.Client
	Config *config.Settings
	Logger *zerolog.Logger
	Models *models.Models
	Utils  *utilities.Utils

	T          *translations.Translations
	Lang       string
	LangBundle *i18n.Bundle
}

func New(cfg *Config) *gocron.Scheduler {
	s := gocron.NewScheduler(time.UTC)

	localizer := i18n.NewLocalizer(cfg.LangBundle, cfg.Lang)
	t := &translations.Translations{Localizer: localizer}
	t.TranslateModels()
	cfg.T = t

	j1, err := s.Every(1).Minute().Do(ScheduledPrayerNotifications, cfg)
	if err != nil {
		cfg.Logger.Error().Err(err).Msg("job running ScheduledPrayerNotifications error")
	}
	cfg.Logger.Info().Msgf(
		"scheduled time: %s, %s",
		j1.GetName(),
		j1.ScheduledTime().String(),
	)

	return s
}
