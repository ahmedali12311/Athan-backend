package seeders

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Users ===================================================
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
	hash, _ := bcrypt.GenerateFromPassword([]byte(c.Password), 12)
	return hash
}

func Users(db *sqlx.DB, qb *squirrel.StatementBuilderType) error {
	cols := []string{
		"id",
		"ref",
		"name",
		"phone",
		"email",
		"password_hash",
	}

	for _, v := range users {
		values := []any{
			v.ID,
			"",
			v.Name,
			v.Phone,
			v.Email,
			v.GeneratePasswordHash(),
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
			v.ID,
			v.Role,
		); err != nil {
			return err
		}
	}

	RunningSeedTable.Append(len(users), "users")
	return nil
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
