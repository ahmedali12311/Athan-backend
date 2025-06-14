package fcm_notification

import (
	"bitbucket.org/sadeemTechnology/backend-finder"
	"bitbucket.org/sadeemTechnology/backend-model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	selects = &[]string{
		"fcm_notifications.*",

		"senders.id as \"sender.id\"",
		"senders.name as \"sender.name\"",
		"senders.phone as \"sender.phone\"",
		"senders.email as \"sender.email\"",
	}
	joins = &[]string{
		"users as senders ON fcm_notifications.sender_id = senders.id",
	}
	inserts = &[]string{
		"title",
		"body",
		"topic",
		"is_sent",
		"send_at",
		"sender_id",
		"response",
		"data",
	}
)

func buildInput(m *Model) (*[]any, error) {
	input := &[]any{
		m.Title,
		m.Body,
		m.Topic,
		m.IsSent,
		m.SendAt,
		m.SenderID,
		m.Response,
		m.Data,
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

type WhereScope struct {
	SenderID *uuid.UUID
}

func (m *Queries) GetAll(
	ctx echo.Context,
	_ *WhereScope,
) (*finder.IndexResponse[*Model], error) {
	c := &finder.ConfigIndex{
		DB:      m.DB,
		QB:      m.QB,
		PGInfo:  m.PGInfo,
		Joins:   joins,
		Selects: selects,
		GroupBys: &[]string{
			"fcm_notifications.id",
			"senders.id",
		},
	}
	return finder.IndexBuilder[*Model](ctx.QueryParams(), c)
}

func (m *Queries) GetOne(shown *Model, _ *WhereScope) error {
	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Joins:   joins,
		Selects: selects,
	}
	return finder.ShowOne(shown, c)
}

func (m *Queries) CreateOne(created *Model) error {
	input, err := buildInput(created)
	if err != nil {
		return err
	}
	c := &finder.ConfigStore{
		DB:      m.DB,
		QB:      m.QB,
		Input:   input,
		Joins:   joins,
		Inserts: inserts,
		Selects: selects,
	}
	return finder.CreateOne(created, c)
}

func (m *Queries) UpdateOne(
	updated *Model,
	_ *WhereScope,
	conn finder.Connection,
) error {
	input, err := buildInput(updated)
	if err != nil {
		return err
	}
	c := &finder.ConfigUpdate{
		DB:      conn,
		QB:      m.QB,
		Input:   input,
		Joins:   joins,
		Inserts: inserts,
		Selects: selects,
		OptimisticLock: &finder.OptimisticLock{
			Name:  "updated_at",
			Value: updated.UpdatedAt,
		},
	}
	return finder.UpdateOne(updated, c)
}

func (m *Queries) DeleteOne(
	deleted *Model,
	_ *WhereScope,
	conn finder.Connection,
) error {
	c := &finder.ConfigDelete{
		DB:      conn,
		QB:      m.QB,
		Joins:   joins,
		Selects: selects,
	}
	return finder.DeleteOne(deleted, c)
}
