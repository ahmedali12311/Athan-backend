package city

import (
	finder "bitbucket.org/sadeemTechnology/backend-finder"
	model "bitbucket.org/sadeemTechnology/backend-model"
	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
)

var (
	selects = &[]string{
		"cities.id",
		"cities.name",
		"cities.is_disabled",
		"ST_AsGeoJSON(cities.location)::json as \"location\"",
		"cities.created_at",
		"cities.updated_at",
	}
	inserts = &[]string{
		"id",
		"name",
		"location",
		"is_disabled",
	}
	joins = &[]string{}
)

type WhereScope struct {
	IsPublic   bool
	IsMarketer bool
}

func wheres(ws *WhereScope) *[]squirrel.Sqlizer {
	w := []squirrel.Sqlizer{}

	// Always filter disabled cities for public users
	if ws.IsPublic {
		w = append(w, squirrel.Expr("cities.is_disabled = false"))
	}

	return &w
}

func buildInput(m *Model) (*[]any, error) {
	input := &[]any{
		m.ID,
		m.Name,
		squirrel.Expr("ST_GeomFromGeoJSON(?::json)", m.Location),
		m.IsDisabled,
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
	cfg := &finder.ConfigIndex{
		DB:       m.DB,
		QB:       m.QB,
		PGInfo:   m.PGInfo,
		Joins:    joins,
		Selects:  selects,
		Wheres:   wheres(ws),
		GroupBys: &[]string{},
	}

	indexResponse, err := finder.IndexBuilder[*Model](ctx.QueryParams(), cfg)
	if err != nil {
		return nil, err
	}

	return indexResponse, nil
}

func (m *Queries) GetOne(shown *Model, ws *WhereScope) error {
	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Joins:   joins,
		Wheres:  wheres(ws),
		Selects: selects,
	}
	if err := finder.ShowOne(shown, c); err != nil {
		return err
	}

	return nil
}

func (m *Queries) CreateOne(created *Model, conn finder.Connection) error {
	input, err := buildInput(created)
	if err != nil {
		return err
	}
	c := &finder.ConfigStore{
		DB:      conn,
		QB:      m.QB,
		Input:   input,
		Joins:   joins,
		Inserts: inserts,
		Selects: selects,
	}
	if err := finder.CreateOne(created, c); err != nil {
		return err
	}

	return nil
}

func (m *Queries) UpdateOne(
	updated *Model,
	ws *WhereScope,
	conn finder.Connection,
) error {
	input, err := buildInput(updated)
	if err != nil {
		return err
	}
	c := &finder.ConfigUpdate{
		DB:      conn,
		QB:      m.QB,
		Joins:   joins,
		Input:   input,
		Inserts: inserts,
		Selects: selects,
	}

	if err := finder.UpdateOne(updated, c); err != nil {
		return err
	}
	return nil
}

func (m *Queries) DeleteOne(
	deleted *Model,
	ws *WhereScope,
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
