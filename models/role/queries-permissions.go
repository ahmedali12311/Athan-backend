package role

import (
	"context"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	"github.com/Masterminds/squirrel"
)

func (m *Queries) GrantAllPermissions(roleID int) (int64, error) {
	perms := []int{}

	if err := m.DB.SelectContext(
		context.Background(),
		&perms,
		`select id from permissions`,
	); err != nil {
		return 0, err
	}

	inserts := m.QB.
		Insert("role_permissions").
		Columns("role_id", "permission_id")
	for _, v := range perms {
		inserts = inserts.Values(roleID, v)
	}
	inserts = inserts.Suffix(`ON CONFLICT DO NOTHING`)
	query, args, err := inserts.ToSql()
	if err != nil {
		return 0, err
	}
	result, err := m.DB.ExecContext(context.Background(), query, args...)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}
func (m *Queries) CreateAdmin() (int, error) {
	var id int
	query := `
        INSERT INTO roles (name) 
        VALUES ('admin') 
        ON CONFLICT (name) 
        DO UPDATE SET name='admin'
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

func (m *Queries) GrantByScope(roleID int, scope string) (int64, error) {
	perms := []int{}

	if err := m.DB.SelectContext(
		context.Background(),
		&perms,
		`
            SELECT id FROM permissions WHERE scope IN ($1, 'own')
        `,
		scope,
	); err != nil {
		return 0, err
	}

	if len(perms) == 0 {
		return 0, nil
	}

	inserts := m.QB.
		Insert("role_permissions").
		Columns("role_id", "permission_id")
	for _, v := range perms {
		inserts = inserts.Values(roleID, v)
	}
	inserts = inserts.Suffix(`ON CONFLICT DO NOTHING`)

	query, args, err := inserts.ToSql()
	if err != nil {
		return 0, err
	}

	result, err := m.DB.ExecContext(context.Background(), query, args...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}
func (m *Queries) GetPermissions(role *Model) error {
	role.Permissions = []int{}

	query := `
        SELECT
            permissions.id
        FROM
            permissions
            INNER JOIN role_permissions r ON r.permission_id = permissions.id 
            INNER JOIN roles ON r.role_id = roles.id 
        WHERE
            roles.id = $1
        ORDER BY
            id
	`

	return m.DB.SelectContext(
		context.Background(),
		&role.Permissions,
		query,
		role.ID,
	)
}

func (m *Queries) SyncPermissions(role *Model, conn finder.Connection) error {
	if _, err := conn.ExecContext(
		context.Background(),
		`DELETE FROM role_permissions WHERE role_id = $1`,
		role.ID,
	); err != nil {
		return err
	}
	query, args, err := m.QB.
		Select("id").
		From("permissions").
		Where(squirrel.Eq{"id": role.Permissions}).
		ToSql()
	if err != nil {
		return err
	}
	var permissions []int
	if err := conn.Select(&permissions, query, args...); err != nil {
		return err
	}
	if len(permissions) == 0 {
		return nil
	}
	rolePerms := m.QB.
		Insert("role_permissions").
		Columns("role_id", "permission_id").
		Suffix("RETURNING permission_id as id")

	for _, value := range permissions {
		rolePerms = rolePerms.Values(role.ID, value)
	}

	query2, args2, err := rolePerms.ToSql()
	if err != nil {
		return err
	}
	if err := conn.SelectContext(
		context.Background(),
		&role.Permissions,
		query2,
		args2...,
	); err != nil {
		return err
	}
	return nil
}

func (m *Queries) RevokeAllPermissions(roleID int) (int64, error) {
	result, err := m.DB.ExecContext(
		context.Background(),
		`
            DELETE FROM role_permissions WHERE role_id = $1
        `,
		roleID,
	)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}
