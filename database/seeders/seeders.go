package seeders

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func SeedError(table string, err error) error {
	return errors.Errorf("[%s] couldn't seed data, %s", table, err.Error())
}

func SeedFixtures(db *sqlx.DB, qb *squirrel.StatementBuilderType) error {

	// settings
	settings, err := SettingsLoadFixtures()
	if err != nil {
		return SeedError(settings.Table(), err)
	}
	err = settings.Seed(db, qb)
	if err != nil {
		return SeedError(settings.Table(), err)
	}

	// roles
	roles, err := RolesLoadFixtures()
	if err != nil {
		return SeedError(roles.Table(), err)
	}
	err = roles.Seed(db, qb)
	if err != nil {
		return SeedError(roles.Table(), err)
	}

	// categories
	categories, err := CategoriesLoadFixtures()
	if err != nil {
		return SeedError(categories.Table(), err)
	}
	err = categories.Seed(db, qb)
	if err != nil {
		return SeedError(categories.Table(), err)
	}

	// users
	users, err := UsersLoadFixtures()
	if err != nil {
		return SeedError(users.Table(), err)
	}
	err = users.Seed(db, qb)
	if err != nil {
		return SeedError(users.Table(), err)
	}

	return nil
}
