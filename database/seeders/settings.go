package seeders

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Settings ===================================================
type Setting struct {
	ID         int    `toml:"id"`
	Key        string `toml:"key"`
	Value      string `toml:"value"`
	IsDisabled bool   `toml:"is_disabled"`
	IsReadOnly bool   `toml:"is_readonly"`
	FieldType  string `toml:"field_type"`
	DataType   string `toml:"data_type"`
}

type Settings struct {
	Elements []Setting `toml:"settings"`
}

func (Settings) File() string {
	return "settings.toml"
}

func (Settings) Table() string {
	return "settings"
}

func (s *Settings) Seed(db *sqlx.DB, qb *squirrel.StatementBuilderType) error {
	cols := []string{
		"id",
		"key",
		"value",
		"is_disabled",
		"is_readonly",
		"field_type",
		"data_type",
	}

	for _, v := range s.Elements {
		values := []any{
			v.ID,
			v.Key,
			v.Value,
			v.IsDisabled,
			v.IsReadOnly,
			v.FieldType,
			v.DataType,
		}
		err := genericSeeder(db, qb, s.Table(), cols, values)
		if err != nil {
			return err
		}
	}

	if _, err := db.Exec(
		`SELECT setval('settings_id_seq', (SELECT MAX(id) FROM settings));`,
	); err != nil {
		return errors.Errorf("error executing sql sequence update settings: %s", err.Error())
	}
	RunningSeedTable.Append(len(s.Elements), s.Table())
	return nil
}

func SettingsLoadFixtures() (*Settings, error) {
	var settings Settings

	err := loadFixtures(&settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}
