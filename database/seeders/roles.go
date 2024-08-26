package seeders

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Roles ===================================================
type Role struct {
	ID   int    `toml:"id"`
	Name string `toml:"name"`
}

type Roles struct {
	Elements []Role `toml:"roles"`
}

func (Roles) File() string {
	return "roles.toml"
}

func (Roles) Table() string {
	return "roles"
}

func (s *Roles) Seed(db *sqlx.DB, qb *squirrel.StatementBuilderType) error {
	cols := []string{
		"id",
		"name",
	}

	for _, v := range s.Elements {
		values := []any{
			v.ID,
			v.Name,
		}
		err := genericSeeder(db, qb, s.Table(), cols, values)
		if err != nil {
			return err
		}
	}

	if _, err := db.Exec(
		`SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));`,
	); err != nil {
		return errors.Errorf("error executing sql sequence update roles: %s", err.Error())
	}
	RunningSeedTable.Append(len(s.Elements), s.Table())
	return nil
}

func RolesLoadFixtures() (*Roles, error) {
	var roles Roles

	err := loadFixtures(&roles)
	if err != nil {
		return nil, err
	}
	return &roles, nil
}
