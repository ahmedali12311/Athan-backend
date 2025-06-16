package api

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	config "bitbucket.org/sadeemTechnology/backend-config"
	category "bitbucket.org/sadeemTechnology/backend-model-category"
	"github.com/rs/zerolog"
)

func ModelMigrator(
	logger *zerolog.Logger,
	cfg *config.Settings,
) {
	migrationsDir := config.GetRootPath(cfg.MigrationsRoot)

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		panic(err)
	}
	var current, next string

	type migration struct {
		name string
		up   []byte
		down []byte

		alreadyApplied bool
		upFileName     string
		downFileName   string
	}

	appModels := []migration{
		{
			name: category.MigrationName,
			up:   category.MigrationUp,
			down: category.MigrationDown,
		},
	}

	var currentInt int

	for _, f := range files {
		fname := strings.SplitN(f.Name(), "_", 2)
		if len(fname) > 0 {
			current = fname[0]
			var appliedFileName string
			migrationFileName := strings.Split(fname[1], ".")
			if len(migrationFileName) > 0 {
				appliedFileName = migrationFileName[0]
			}
			for i := range appModels {
				if appModels[i].name == appliedFileName {
					appModels[i].alreadyApplied = true
				}
			}
			val, err := strconv.Atoi(current)
			if err != nil {
				logger.Error().Err(err).Msg("error getting migration number")
			}
			currentInt = val
		}
	}
	for _, m := range appModels {
		if !m.alreadyApplied {
			logger.Info().Msgf("applying migration for model: %s", m)
			currentInt += 1
			next = fmt.Sprintf("%06d", currentInt)
			m.upFileName = fmt.Sprintf(
				"%s/%s_%s.up.sql",
				migrationsDir,
				next,
				m.name,
			)
			m.downFileName = fmt.Sprintf(
				"%s/%s_%s.down.sql",
				migrationsDir,
				next,
				m.name,
			)

			if err := os.WriteFile(
				m.upFileName,
				m.up,
				0o644,
			); err != nil {
				panic(err)
			}
			if err := os.WriteFile(
				m.downFileName,
				m.down,
				0o644,
			); err != nil {
				panic(err)
			}
		}
	}
}
