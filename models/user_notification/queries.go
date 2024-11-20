package user_notification

import (
	"app/model"
	"context"
	"errors"

	"github.com/m-row/finder"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

var (
	selects = &[]string{
		"user_notifications.*",

		"users.id as \"user.id\"",
		"users.name as \"user.name\"",
		"users.phone as \"user.phone\"",
		"users.email as \"user.email\"",
	}

	joins = &[]string{
		"users ON user_notifications.user_id = users.id",
		"reservations r ON (user_notifications.data->>'reservation')::uuid = r.id",
	}

	inserts = &[]string{
		"user_id",
		"is_read",
		"is_notified",
		"title",
		"body",
		"response",
		"data",
	}
)

type WhereScope struct {
	UserID   *uuid.UUID
	IsPublic bool
	Method   string
}

func wheres(ws *WhereScope) *[]squirrel.Sqlizer {
	w := &[]squirrel.Sqlizer{}

	*w = append(*w, squirrel.Or{
		squirrel.Eq{"user_notifications.data": nil},
		squirrel.Eq{"user_notifications.data->>'business'": nil},
	})

	if ws.IsPublic && ws.Method == "GET" {
		if ws.UserID != nil {
			*w = append(*w, squirrel.Or{
				squirrel.Eq{"user_notifications.user_id": ws.UserID},
				squirrel.Eq{"user_notifications.user_id": nil},
			})
		} else {
			*w = append(*w, squirrel.Or{
				squirrel.Eq{"user_notifications.user_id": nil},
			})
		}
	}

	if ws.Method != "GET" {
		*w = append(*w,
			squirrel.Eq{"user_notifications.user_id": ws.UserID},
		)
	}

	return w
}

func buildInput(m *Model) (*[]any, error) {
	data, err := m.Data.Value()
	if err != nil {
		return nil, err
	}
	input := &[]any{
		m.User.ID,
		m.IsRead,
		m.IsNotified,
		m.Title,
		m.Body,
		m.Response,
		data,
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
	ws *WhereScope,
) (*finder.IndexResponse[*Model], error) {
	config := &finder.ConfigIndex{
		DB:      m.DB,
		QB:      m.QB,
		Selects: selects,
		Joins:   joins,
		Wheres:  wheres(ws),
		PGInfo:  m.PGInfo,
		GroupBys: &[]string{
			"user_notifications.id",
			"users.id",
		},
	}
	return finder.IndexBuilder[*Model](ctx.QueryParams(), config)
}

func (m *Queries) GetOne(
	shown *Model,
	ws *WhereScope,
) error {

	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Selects: selects,
		Joins:   joins,
		Wheres:  wheres(ws),
	}
	return finder.ShowOne(shown, c)
}

func (m *Queries) Create(created *Model) error {
	input, err := buildInput(created)
	if err != nil {
		return err
	}
	c := &finder.ConfigStore{
		DB:      m.DB,
		QB:      m.QB,
		Input:   input,
		Inserts: inserts,
		Selects: selects,
		Joins:   joins,
	}
	return finder.CreateOne(created, c)
}

func (m *Queries) Update(
	updated *Model,
	ws *WhereScope,
	tx *sqlx.Tx,
) error {
	input, err := buildInput(updated)
	if err != nil {
		return err
	}
	c := &finder.ConfigUpdate{
		DB:      tx,
		QB:      m.QB,
		Input:   input,
		Inserts: inserts,
		Selects: selects,
		Joins:   joins,
		Wheres:  wheres(ws),
		OptimisticLock: &finder.OptimisticLock{
			Name:  "updated_at",
			Value: updated.UpdatedAt,
		},
	}
	return finder.UpdateOne(updated, c)
}

func (m *Queries) Delete(
	deleted *Model,
	ws *WhereScope,
) error {

	c := &finder.ConfigDelete{
		DB:      m.DB,
		QB:      m.QB,
		Selects: selects,
		Joins:   joins,
		Wheres:  wheres(ws),
	}
	return finder.DeleteOne(deleted, c)
}

func (m *Queries) BulkCreate(notifications []Model) error {
	inserts := m.QB.
		Insert("user_notifications").
		Columns(
			"user_id",
			"is_read",
			"is_notified",
			"title",
			"body",
			"response",
			"data",
		)
	for _, notification := range notifications {
		data, err := notification.Data.Value()
		if err != nil {
			return err
		}
		inserts = inserts.Values(
			notification.User.ID,
			notification.IsRead,
			notification.IsNotified,
			notification.Title,
			notification.Body,
			notification.Response,
			data,
		)
	}
	query, args, err := inserts.ToSql()
	if err != nil {
		return err
	}
	if _, err := m.DB.ExecContext(context.Background(), query, args...); err != nil {
		return err
	}
	return nil
}

func (m *Queries) ToggleRead(toggled *Model, ws *WhereScope) error {
	if ws.UserID == nil {
		return errors.New("user_id is required")
	}
	subquery := m.QB.
		Update("user_notifications").
		Set("is_read", squirrel.Expr("NOT is_read")).
		Where("id = ?", toggled.ID).
		Where("user_id = ?", ws.UserID).
		Suffix("RETURNING *")
	with := subquery.Prefix("WITH user_notifications AS (").Suffix(")")

	query, args, err := m.QB.
		Select(*selects...).
		PrefixExpr(with).
		From("user_notifications").
		LeftJoin("users ON user_notifications.user_id = users.id").
		ToSql()
	if err != nil {
		return err
	}
	if err := m.DB.GetContext(context.Background(), toggled, query, args...); err != nil {
		return err
	}
	return nil
}

func (m *Queries) MarkAllRead(ws *WhereScope) (int64, error) {
	if ws.UserID == nil {
		return 0, errors.New("user_id is required")
	}
	query, args, err := m.QB.
		Update("user_notifications").
		Set("is_read", true).
		Where("user_id = ?", ws.UserID).
		ToSql()
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

func (m *Queries) MarkAllUnread(ws *WhereScope) (int64, error) {
	if ws.UserID == nil {
		return 0, errors.New("user_id is required")
	}
	query, args, err := m.QB.
		Update("user_notifications").
		Set("is_read", false).
		Where("user_id = ?", ws.UserID).
		ToSql()
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
