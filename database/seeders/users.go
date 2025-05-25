package seeders

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID
	Ref          string
	Name         string
	Phone        string
	Email        string
	Password     string
	PasswordHash []byte
	Role         int
}

func (c *User) GeneratePasswordHash() []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), 12)
	if err != nil {
		log.Panicf(
			"error executing GeneratePasswordHash: %s",
			err.Error(),
		)
	}
	return hash
}

func Users(db *sqlx.DB, qb *squirrel.StatementBuilderType) {
	cols := []string{
		"id",
		"ref",
		"name",
		"phone",
		"email",
		"password_hash",
	}

	for i := range users {
		values := []any{
			users[i].ID,
			users[i].Ref,
			users[i].Name,
			users[i].Phone,
			users[i].Email,
			users[i].GeneratePasswordHash(),
		}
		genericSeeder(db, qb, "users", cols, values)

		query := `
            INSERT INTO user_roles (
                user_id,
                role_id
            ) 
            VALUES (
                $1,
                $2
            ) 
            ON CONFLICT (user_id, role_id) 
            DO NOTHING;
        `
		if _, err := db.ExecContext(context.Background(),
			query,
			users[i].ID,
			users[i].Role,
		); err != nil {
			log.Panicf(
				"error executing sql for user_roles insert: %s",
				err.Error(),
			)
		}
	}

	RunningSeedTable.Append(len(users), "users")
}

var users = []User{
	{
		ID:       parseUUID("280a8eb8-8add-4287-a33b-c21a4ada2eaf"),
		Ref:      "BA32D47B",
		Name:     "admin",
		Phone:    "218920000001",
		Email:    "admin@example.com",
		Password: "password",
		Role:     2,
	},
	{
		ID:       parseUUID("f97721e4-8471-4513-9c51-d21326df4f18"),
		Ref:      "839331E3",
		Name:     "customer1",
		Phone:    "218920000002",
		Email:    "customer1@example.com",
		Password: "password",
		Role:     3,
	},
}
