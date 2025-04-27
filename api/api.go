package api

import (
	"context"
	"errors"

	"app/apierrors"
	"app/config"
	"app/controller"
	"app/controllers"
	"app/database"
	"app/database/seeders"
	"app/models"
	"app/models/permission"
	"app/models/user"
	"app/utilities"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog"
	js "github.com/santhosh-tekuri/jsonschema/v5"
	"google.golang.org/api/option"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

var (
	CommitCount    = "0"
	CommitDescribe = "dev"
	Version        = "1." + CommitCount + "." + CommitDescribe
)

type Application struct {
	DB          *sqlx.DB
	Lang        string
	LangBundle  *i18n.Bundle
	Config      *config.Settings
	Logger      *zerolog.Logger
	Utils       *utilities.Utils
	APIErrors   *apierrors.APIErrors
	Models      *models.Models
	Controllers *controllers.Controllers
	FB          *firebase.App
	FBM         *messaging.Client
	CtxUser     *user.Model
	Permissions []permission.Model
}

func NewAPI(
	cfg *config.Settings,
	bundle *i18n.Bundle,
	schemas map[string]*js.Schema,
	logger *zerolog.Logger,
	isTest bool,
) *Application {
	// init firebase app
	opt := option.WithCredentialsFile(config.GoogleServiceAccount)
	fb, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil && !isTest {
		logger.Fatal().
			Msgf("firebase: error initializing app: %s", err.Error())
	}
	fbm, err := fb.Messaging(context.Background())
	if err != nil && !isTest {
		logger.Fatal().
			Msgf("couldn't get firebase messaging client, %s", err.Error())
	}
	logger.Info().Msg("firebase app and messaging client initiated")

	utils := utilities.Get(
		&Version,
		&CommitDescribe,
		&CommitCount,
		cfg,
		logger,
		fb,
		fbm,
	)
	apiErr := apierrors.Get(utils, logger)

	db, err := database.OpenSQLX(cfg.DSN)
	if err != nil {
		logger.Fatal().Msgf("couldn't open db: %s", err.Error())
	}
	logger.Info().Msg("database connection pool established")
	logger.Info().Msg("file:///" + config.GetRootPath(config.MigrationsRoot))
	mig, err := migrate.New(
		"file:///"+config.GetRootPath(cfg.MigrationsRoot),
		cfg.DSN,
	)
	if err != nil {
		logger.Fatal().
			Msgf("couldn't create migration instance: %s", err.Error())
	}
	if err := mig.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			logger.Fatal().
				Msgf("couldn't run migration up: %s", err.Error())
		}
		logger.Info().Msgf("migrations: %s", err.Error())
	}
	migv, migdirty, err := mig.Version()
	if err != nil {
		logger.Fatal().Msgf("couldn't close migration err: %s", err.Error())
	}
	logger.Info().Msgf("database migration version: %d, %t", migv, migdirty)

	err1, err2 := mig.Close()
	if err1 != nil {
		logger.Fatal().Msgf("couldn't close migration err1: %s", err1.Error())
	}
	if err2 != nil {
		logger.Fatal().Msgf("couldn't close migration err2: %s", err2.Error())
	}

	tc := make(map[string][]string)
	if err := database.PGInfo(db, tc); err != nil {
		logger.Fatal().Msgf("couldn't get pgInfo: %s", err.Error())
	}
	logger.Info().Msg("database table and column info saved")

	m := models.Setup(db, fb, fbm, tc)

	// initial data seeders, only adds and avoids overwriting on conflict
	seeders.Settings(m.DB, m.QB)
	seeders.Categories(m.DB, m.QB)
	seeders.Roles(m.DB, m.QB)
	seeders.Users(m.DB, m.QB)
	seeders.PrintTable(seeders.RunningSeedTable)

	// sms, err := sms.NewSMS(
	// 	sms.SMSConfig{
	// 		Url: config.OTP_URL,
	// 		Jwt: config.OTP_JWT,
	// 		Key: config.OTP_KEY,
	// 	},
	// 	sms.WithLogger(logger))

	// if err != nil {
	// 	logger.Fatal().Msg(err.Error())
	// }

	deps := &controller.Dependencies{
		Schemas: schemas,
		Utils:   utils,
		APIErr:  apiErr,
		Models:  m,
		// SMS:     sms,
	}

	ctrls := controllers.Setup(deps)

	newApi := &Application{
		DB:          db,
		Lang:        "ar",
		LangBundle:  bundle,
		Config:      cfg,
		Logger:      logger,
		Utils:       utils,
		APIErrors:   apiErr,
		Models:      m,
		Controllers: ctrls,
		FB:          fb,
		FBM:         fbm,
		CtxUser:     nil,
	}
	return newApi
}
