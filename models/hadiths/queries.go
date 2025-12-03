package hadiths

import (
	"net/url"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	model "bitbucket.org/sadeemTechnology/backend-model"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

var (
	selects = &[]string{
		"hadithses.*",
		"c.id as \"category.id\"",
		"c.name as \"category.name\"",
	}

	inserts = &[]string{
		"text",
		"source",
		"topic",
		"category_id",
	}
	baseJoins = &[]string{
		"categories ON hadithses.category_id = categories.id",
	}
)

func buildInput(m *Model) (*[]any, error) {
	input := &[]any{
		m.Text,
		m.Source,
		m.Topic,
		m.CategoryID,
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
	IsAdmin     bool
	IsPublic    bool
	UserID      *uuid.UUID
	QueryParams url.Values
}

func getJoins(ws *WhereScope) *[]string {
	return baseJoins
}

func wheres(ws *WhereScope) *[]squirrel.Sqlizer {
	w := []squirrel.Sqlizer{}
	if ws.IsAdmin {
		return &w
	}

	if ws.UserID != nil {
	}

	if !ws.IsAdmin {
	}

	return &w
}

func (m *Queries) GetAll(
	ctx echo.Context,
	ws *WhereScope,
) (*finder.IndexResponse[*Model], error) {
	c := &finder.ConfigIndex{
		DB:      m.DB,
		QB:      m.QB,
		PGInfo:  m.PGInfo,
		Wheres:  wheres(ws),
		Selects: selects,
		Joins:   getJoins(ws),
	}
	return finder.IndexBuilder[*Model](ctx.QueryParams(), c)
}

func (m *Queries) GetOne(shown *Model, ws *WhereScope) error {
	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Wheres:  wheres(ws),
		Selects: selects,
		Joins:   getJoins(ws),
	}
	return finder.ShowOne(shown, c)
}

func (m *Queries) CreateOne(created *Model, tx *sqlx.Tx) error {
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
	}
	return finder.CreateOne(created, c)
}

func (m *Queries) UpdateOne(updated *Model, ws *WhereScope, tx *sqlx.Tx) error {
	input, err := buildInput(updated)
	if err != nil {
		return err
	}
	c := &finder.ConfigUpdate{
		DB:      m.DB,
		QB:      m.QB,
		Input:   input,
		Wheres:  wheres(ws),
		Inserts: inserts,
		Selects: selects,
		OptimisticLock: &finder.OptimisticLock{
			Name:  "updated_at",
			Value: updated.UpdatedAt,
		},
	}
	return finder.UpdateOne(updated, c)
}

func (m *Queries) DeleteOne(deleted *Model, ws *WhereScope, tx *sqlx.Tx) error {
	c := &finder.ConfigDelete{
		DB:      m.DB,
		QB:      m.QB,
		Wheres:  wheres(ws),
		Selects: selects,
	}
	return finder.DeleteOne(deleted, c)
}
