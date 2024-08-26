package seeders

import (
	"context"
	"io"
	"os"
	"strings"

	"app/config"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Users ===================================================
type User struct {
	ID           uuid.UUID `toml:"id"`
	Name         string    `toml:"name"`
	Phone        string    `toml:"phone"`
	Email        string    `toml:"email"`
	Password     string    `toml:"password"`
	PasswordHash []byte
	Role         int `toml:"role"`
}

func (c *User) GeneratePasswordHash() []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(c.Password), 12)
	return hash
}

type Users struct {
	Elements []User `toml:"users"`
}

func (c Users) File() string {
	return "users.toml"
}

func (Users) Table() string {
	return "users"
}

func (s *Users) Seed(db *sqlx.DB, qb *squirrel.StatementBuilderType) error {
	cols := []string{
		"id",
		"ref",
		"name",
		"phone",
		"email",
		"password_hash",
	}

	for _, v := range s.Elements {
		values := []any{
			v.ID,
			"",
			v.Name,
			v.Phone,
			v.Email,
			v.GeneratePasswordHash(),
		}
		err := genericSeeder(db, qb, s.Table(), cols, values)
		if err != nil {
			return err
		}
		query2 := `
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
			query2,
			v.ID,
			v.Role); err != nil {
			return err
		}
	}

	RunningSeedTable.Append(len(s.Elements), s.Table())
	return nil
}

func UsersLoadFixtures() (*Users, error) {
	var users Users

	err := loadFixtures(&users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func loadFixtures(data Fixtures) error {
	current_dir, err := os.Getwd()
	if err != nil {
		return err
	}

	config.GetRootPath(config.SeedersRoot)
	filename := strings.Join([]string{
		current_dir,
		"database", "seeders", "fixtures_data", data.File(),
	},
		"/")

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(b, data)
	if err != nil {
		return err
	}

	return nil
}
