package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SetFCMToken sets the token value and type for profile, defaults type to
// fcm_token_customer.
func (m *Queries) SetFCMToken(
	id, tokenType, tokenValue *string,
	tx *sqlx.Tx,
) error {
	query := `
        UPDATE tokens 
        SET token_value=$3 
        WHERE user_id=$1 
              AND token_type=$2
    `
	result, err := tx.ExecContext(
		context.Background(),
		query,
		*id,
		*tokenType,
		*tokenValue,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		// No row exist so create one
		query := `
            INSERT INTO tokens 
                (user_id, token_type, token_value) 
            VALUES 
                ($1, $2, $3)
        `
		if _, err := tx.ExecContext(
			context.Background(),
			query,
			*id,
			*tokenType,
			*tokenValue,
		); err != nil {
			return err
		}
	}

	return nil
}

func (m *Queries) GetFCMToken(
	id *uuid.UUID,
	tokenType string,
) (*string, error) {
	var token string

	query := `
        SELECT token_value 
        FROM tokens 
        WHERE user_id=$1 
              AND token_type=$2
    `
	if err := m.DB.GetContext(
		context.Background(),
		&token,
		query,
		id,
		tokenType,
	); err != nil {
		return nil, err
	}
	return &token, nil
}

func (m *Queries) GetFCMTokens(notifiableUsersID *[]string) (*[]string, error) {
	var tokens []string
	results := m.QB.
		Select("token_value").
		From("tokens").
		LeftJoin("users ON tokens.user_id = users.id").
		Where(squirrel.Eq{"users.is_notifiable": true}).
		Where(squirrel.Eq{"user_id": *notifiableUsersID})
	query, args, err := results.ToSql()
	if err != nil {
		return nil, err
	}
	if err := m.DB.SelectContext(
		context.Background(),
		&tokens,
		query,
		args...,
	); err != nil {
		return nil, err
	}
	return &tokens, nil
}

func (m *Queries) GetRoleFCMTokens(
	roleName, tokenType string,
) (*[]string, error) {
	var tokens []string
	query := `
          SELECT
              token_value
          FROM
              tokens
          WHERE
              token_type = $2
              AND user_id IN (
                  SELECT
                      user_id
                  FROM
                      user_roles
                      JOIN users ON user_roles.user_id = users.id
                      JOIN roles ON user_roles.role_id = roles.id
                  WHERE
                      roles.name = $1
                      AND users.is_disabled IS FALSE
                      AND users.is_notifiable IS TRUE);
	`
	if err := m.DB.SelectContext(
		context.Background(),
		&tokens,
		query,
		roleName,
		tokenType,
	); err != nil {
		return nil, err
	}
	return &tokens, nil
}

func (m *Queries) ClearTokens(userID *uuid.UUID, tx *sqlx.Tx) error {
	query := `DELETE FROM tokens WHERE user_id=$1`
	if tx != nil {
		if _, err := tx.ExecContext(
			context.Background(),
			query,
			userID,
		); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
	} else {
		if _, err := m.DB.ExecContext(
			context.Background(),
			query,
			userID,
		); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
	}
	return nil
}
