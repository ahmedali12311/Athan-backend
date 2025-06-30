package role

import (
	"context"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	"github.com/google/uuid"
)

func (m *Queries) CreateSuperAdmin() (int, error) {
	var id int
	query := `
        INSERT INTO roles (name) 
        VALUES ('superadmin') 
        ON CONFLICT (name) 
        DO UPDATE SET name='superadmin'
        RETURNING id;
    `
	if err := m.DB.GetContext(
		context.Background(),
		&id,
		query,
	); err != nil {
		return 0, err
	}

	if _, err := m.DB.ExecContext(
		context.Background(),
		`SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));`,
	); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *Queries) CreateBasic() (int, error) {
	var id int
	query := `
        INSERT INTO roles (name) 
        VALUES ('basic') 
        ON CONFLICT (name) 
        DO UPDATE SET name='basic'
        RETURNING id;
    `
	if err := m.DB.GetContext(
		context.Background(),
		&id,
		query,
	); err != nil {
		return 0, err
	}

	if _, err := m.DB.ExecContext(
		context.Background(),
		`SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));`,
	); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *Queries) GrantBasic(userID *uuid.UUID, conn finder.Connection) error {
	query := `
        INSERT INTO user_roles (user_id, role_id) 
        VALUES ($1, (SELECT id FROM roles WHERE name='basic')) 
        ON CONFLICT (user_id, role_id) 
        DO NOTHING;
    `
	if _, err := conn.ExecContext(
		context.Background(),
		query,
		userID,
	); err != nil {
		return err
	}
	return nil
}
