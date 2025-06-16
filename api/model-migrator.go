package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	config "bitbucket.org/sadeemTechnology/backend-config"
	category "bitbucket.org/sadeemTechnology/backend-model-category"
	"github.com/rs/zerolog"
)

type migration struct {
	name string
	up   []byte
	down []byte

	alreadyApplied bool
	upFileName     string
	downFileName   string
}

func ModelMigrator(
	logger *zerolog.Logger,
	cfg *config.Settings,
) error {
	dir := config.GetRootPath(cfg.MigrationsRoot)

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	appModels := []migration{
		{
			name: category.MigrationName,
			up:   category.MigrationUp,
			down: category.MigrationDown,
		},
	}

	var maxMigrationInt int

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		parts := strings.SplitN(f.Name(), "_", 2)
		if len(parts) != 2 {
			continue
		}
		indexStr := parts[0]
		val, err := strconv.Atoi(indexStr)
		if err != nil {
			logger.Warn().
				Str("file", f.Name()).
				Msg("skipping file with invalid migration index")
			continue
		}
		if val > maxMigrationInt {
			maxMigrationInt = val
		}
		migrationName := strings.Split(parts[1], ".")[0]
		for i := range appModels {
			if appModels[i].name == migrationName {
				appModels[i].alreadyApplied = true
			}
		}

	}
	for _, m := range appModels {
		if !m.alreadyApplied {
			logger.Info().Msgf("applying migration for: %s", m.name)
			maxMigrationInt += 1
			next := fmt.Sprintf("%06d", maxMigrationInt)

			m.upFileName = filepath.Join(
				dir,
				fmt.Sprintf("%s_%s.up.sql", next, m.name),
			)
			m.downFileName = filepath.Join(
				dir,
				fmt.Sprintf("%s_%s.down.sql", next, m.name),
			)

			if err := os.WriteFile(
				m.upFileName,
				m.up,
				os.ModePerm,
			); err != nil {
				return err
			}
			if err := os.WriteFile(
				m.downFileName,
				m.down,
				os.ModePerm,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
