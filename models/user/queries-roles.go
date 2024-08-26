package user

import (
	"context"
	"errors"

	"app/models/role"

	"github.com/Masterminds/squirrel"
	"github.com/ahmedalkabir/finder"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (m *Queries) GetRoles(user *Model) error {
	user.Roles = &[]int{}
	query := `
          SELECT
              roles.id
          FROM
              roles
              INNER JOIN user_roles ON user_roles.role_id = roles.id
              INNER JOIN users ON user_roles.user_id = users.id
          WHERE
              users.id = $1
          ORDER BY
              roles.id
          
	`
	return m.DB.SelectContext(
		context.Background(),
		user.Roles,
		query,
		user.ID,
	)
}

func (m *Queries) AssignRoles(user *Model, conn finder.Connection) error {
	if _, err := conn.ExecContext(
		context.Background(),
		`DELETE FROM user_roles WHERE user_id=$1`,
		user.ID,
	); err != nil {
		return err
	}
	results := m.QB.
		Select("id", "name").
		From("roles").
		Where(squirrel.Eq{"id": user.Roles})
	query, args, err := results.ToSql()
	if err != nil {
		return err
	}
	var roles []role.Model
	if err := conn.Select(&roles, query, args...); err != nil {
		return err
	} else {
		if len(roles) == 0 {
			return nil
		}
		insertUserRoles := m.QB.
			Insert("user_roles").
			Columns("user_id", "role_id")

		for _, value := range roles {
			// user.Roles = append(user.Roles, role.Value(value.ID))
			insertUserRoles = insertUserRoles.Values(user.ID, value.ID)
		}
		query, args, err := insertUserRoles.ToSql()
		if err != nil {
			return err
		}
		if _, err := conn.ExecContext(
			context.Background(),
			query,
			args...,
		); err != nil {
			return err
		}
		return nil
	}
}

func (m *Queries) GrantRole(userID *uuid.UUID, roleID *int, tx *sqlx.Tx) error {
	var exists bool
	query := `
        SELECT EXISTS (
          SELECT 1
          FROM user_roles
          WHERE user_id = $1
                AND role_id = $2
        ) AS exists
	`
	initialError := errors.New("user already has role")
	if err := m.DB.GetContext(
		context.Background(),
		&exists,
		query,
		userID,
		*roleID,
	); err != nil {
		exists = false
		initialError = err
	}
	if !exists {
		query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`
		if result, err := tx.ExecContext(
			context.Background(),
			query,
			userID,
			roleID,
		); err != nil {
			if result == nil {
				return initialError
			}
			return err
		}
		return nil
	}
	return initialError
}

func (m *Queries) RevokeRole(
	userID *uuid.UUID,
	roleID *int,
	tx *sqlx.Tx,
) error {
	var exists bool
	query := `
          SELECT EXISTS (
              SELECT 1
              FROM user_roles
              WHERE user_id = $1
                AND role_id = $2
          ) AS exists
	`
	if err := m.DB.GetContext(
		context.Background(),
		&exists,
		query,
		userID,
		*roleID,
	); err != nil {
		exists = false
	}
	if exists {
		query2 := `DELETE FROM user_roles WHERE user_id=$1 AND role_id=$2`
		if _, err := tx.ExecContext(
			context.Background(),
			query2,
			userID,
			roleID,
		); err != nil {
			return err
		}
		return nil
	}
	return errors.New("user does not have role")
}
