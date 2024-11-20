package seeders

import (
	"log"

	"app/models/role"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func Roles(db *sqlx.DB, qb *squirrel.StatementBuilderType) {
	cols := []string{
		"id",
		"name",
	}
	for i := range roles {
		values := []any{
			roles[i].ID,
			roles[i].Name,
		}
		genericSeeder(db, qb, "roles", cols, values)
	}
	if _, err := db.Exec(
		`SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));`,
	); err != nil {
		log.Panicf(
			"error executing sql sequence update roles: %s",
			err.Error(),
		)
	}
	RunningSeedTable.Append(len(roles), "roles")
}

var roles = []role.Model{
	{
		ID:   1,
		Name: "superadmin",
	},
	{
		ID:   1,
		Name: "admin",
	},
	{
		ID:   1,
		Name: "customer",
	},
}
