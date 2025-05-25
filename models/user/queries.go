package user

import (
	"context"
	"database/sql"
	"errors"

	"app/config"
	"github.com/m-row/model"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder"
)

var (
	selects = &[]string{
		"users.id",
		"users.ref",
		"users.name",
		"users.phone",
		"users.email",
		"users.password_hash",
		"users.gender",
		"users.is_anon",
		"users.is_notifiable",
		"users.is_disabled",
		"users.is_confirmed",
		"users.is_deleted",
		"users.is_verified",
		"users.last_ref",
		"users.pin",
		"users.pin_expiry",

		"users.created_at",
		"users.updated_at",

		"TO_CHAR(users.birthdate, 'YYYY-MM-DD') as birthdate",
		"ST_AsGeoJSON(users.location)::json as \"location\"",
		config.SQLSelectURLPath("users", "img", "img"),
		config.SQLSelectURLPath("users", "thumb", "thumb"),
	}
	joins   = &[]string{}
	inserts = &[]string{
		"id",
		"ref",
		"email",
		"password_hash",
		"name",
		"phone",
		"gender",
		"birthdate",
		"is_anon",
		"is_notifiable",
		"is_disabled",
		"is_confirmed",
		"is_verified",
		"last_ref",
		"img",
		"thumb",
		"location",
		"pin",
		"pin_expiry",
	}
)

func buildInput(m *Model) (*[]any, error) {
	hash := ""
	if m.Password.Hash != nil {
		// this sets the password hash on create/ or when provided in update
		hash = string(*m.Password.Hash)
	} else if m.PasswordHash != nil {
		// this handle updates when password is not updated
		hash = string(*m.PasswordHash)
	}
	input := &[]any{
		m.ID,
		squirrel.Expr("upper(?)", m.Ref),
		squirrel.Expr("lower(?)", m.Email),
		hash,
		m.Name,
		m.Phone,
		m.Gender,
		m.Birthdate,
		m.IsAnon,
		m.IsNotifiable,
		m.IsDisabled,
		m.IsConfirmed,
		m.IsVerified,
		m.LastRef,
		m.Img,
		m.Thumb,
		squirrel.Expr("ST_GeomFromGeoJSON(?::json)", m.Location),
		m.Pin,
		m.PinExpiry,
	}
	if len(*input) != len(*inserts) {
		return nil, finder.ErrInputLengthMismatch(input, inserts)
	}
	return input, nil
}

type Queries struct {
	*model.Dependencies
}

func New(d *model.Dependencies) *Queries {
	return &Queries{d}
}

func (m *Queries) GetAll(
	ctx echo.Context,
) (*finder.IndexResponse[*Model], error) {
	cfg := &finder.ConfigIndex{
		DB:      m.DB,
		QB:      m.QB,
		PGInfo:  m.PGInfo,
		Joins:   joins,
		Selects: selects,
		GroupBys: &[]string{
			"users.id",
		},
	}
	indexResponse, err := finder.IndexBuilder[*Model](ctx.QueryParams(), cfg)
	if err != nil {
		return nil, err
	}
	if err := m.EagerLoad(indexResponse.Data); err != nil {
		return nil, err
	}
	return indexResponse, nil
}

func (m *Queries) GetOne(shown *Model) error {
	if shown.ID == uuid.Nil && shown.Phone == nil && shown.Email == nil {
		return nil
	}
	wheres := &[]squirrel.Sqlizer{}
	if shown.ID != uuid.Nil {
		*wheres = append(
			*wheres,
			squirrel.Expr("users.id=?", shown.ID.String()),
		)
	}
	if shown.Phone != nil {
		if *shown.Phone != "" {
			expr := squirrel.Expr("users.phone = ?", *shown.Phone)
			*wheres = append(*wheres, expr)
		}
	}
	if shown.Pin != nil {
		if *shown.Pin != "" {
			expr := squirrel.Expr("users.pin = ?", *shown.Pin)
			*wheres = append(*wheres, expr)
		}
	}
	if shown.Email != nil {
		if *shown.Email != "" {
			expr := squirrel.Expr("users.email = LOWER(?)", *shown.Email)
			*wheres = append(*wheres, expr)
		}
	}
	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Joins:   joins,
		Wheres:  wheres,
		Selects: selects,
	}
	if err := finder.ShowOne(shown, c); err != nil {
		return err
	}
	if err := m.GetRoles(shown); err != nil {
		return err
	}

	return m.GetPermissions(shown, c.DB)
}

// CreateOne inserts a user with roles,
//
// requires a transaction.
func (m *Queries) CreateOne(created *Model, conn *sqlx.Tx) error {
	input, err := buildInput(created)
	if err != nil {
		return err
	}
	c := &finder.ConfigStore{
		DB:      conn,
		QB:      m.QB,
		Input:   input,
		Inserts: inserts,
		Selects: selects,
		Joins:   joins,
	}
	if err := finder.CreateOne(created, c); err != nil {
		return err
	}

	return m.AssignRoles(created, conn)
}

func (m *Queries) UpdateOne(updated *Model, conn finder.Connection) error {
	input, err := buildInput(updated)
	if err != nil {
		return err
	}
	c := &finder.ConfigUpdate{
		DB:      conn,
		QB:      m.QB,
		Input:   input,
		Inserts: inserts,
		Selects: selects,
		Joins:   joins,
		OptimisticLock: &finder.OptimisticLock{
			Name:  "updated_at",
			Value: updated.UpdatedAt,
		},
	}
	if err := finder.UpdateOne(updated, c); err != nil {
		return err
	}
	if updated.Roles != nil {
		if len(*updated.Roles) != 0 {
			if err := m.AssignRoles(updated, conn); err != nil {
				return err
			}
			if err := m.GetPermissions(updated, conn); err != nil {
				return err
			}
		}
	}
	return nil
}

// ClearOne Clear user account data.
func (m *Queries) ClearOne(userID *uuid.UUID, conn finder.Connection) error {
	// clear user data
	query1 := `
       UPDATE
           users
       SET
           phone = NULL,
           email = NULL,
           password_hash = NULL,
           name = NULL,
           gender = 'male',
           birthdate = NULL,
           is_notifiable = FALSE,
           is_disabled = TRUE,
           is_confirmed = FALSE,
           is_deleted = TRUE,
           img = NULL,
           thumb = NULL
       WHERE
           id = $1
    `
	if _, err := conn.ExecContext(
		context.Background(),
		query1,
		userID,
	); err != nil {
		return err
	}
	// clear user tokens
	if _, err := conn.ExecContext(
		context.Background(),
		`DELETE FROM tokens WHERE user_id=$1`,
		*userID,
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	// clear user roles
	if _, err := conn.ExecContext(
		context.Background(),
		`DELETE FROM user_roles WHERE user_id=$1`,
		*userID,
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	// clear user permissions
	if _, err := conn.ExecContext(
		context.Background(),
		`DELETE FROM user_permissions WHERE user_id=$1`,
		*userID,
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	return nil
}

func (m *Queries) Verify(userID *uuid.UUID) error {
	if _, err := m.DB.ExecContext(
		context.Background(),
		`
           UPDATE users
           SET is_verified = TRUE
           WHERE id = $1
        `,
		*userID,
	); err != nil {
		return err
	}
	return nil
}
