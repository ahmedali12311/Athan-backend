package setting

import (
	"slices"

	"bitbucket.org/sadeemTechnology/backend-finder"
	"bitbucket.org/sadeemTechnology/backend-model"

	"github.com/labstack/echo/v4"
)

type Queries struct {
	*model.Dependencies
}

func New(d *model.Dependencies) *Queries {
	return &Queries{d}
}

var (
	selects = &[]string{
		"settings.*",
	}

	inserts = &[]string{
		"key",
		"value",
		"is_readonly",
		"is_disabled",
		"field_type",
		"data_type",
	}
)

func buildInput(setting *Model) (*[]any, error) {
	input := &[]any{
		setting.Key,
		setting.Value,
		setting.IsReadOnly,
		setting.IsDisabled,
		setting.FieldType,
		setting.DataType,
	}
	if len(*input) != len(*inserts) {
		return nil, finder.ErrInputLengthMismatch(input, inserts)
	}
	return input, nil
}

func (m *Queries) GetAll(
	ctx echo.Context,
	isPublic bool,
) (*finder.IndexResponse[*Model], error) {
	config := &finder.ConfigIndex{
		DB:       m.DB,
		QB:       m.QB,
		PGInfo:   m.PGInfo,
		IsPublic: isPublic,
		Selects:  selects,
	}
	return finder.IndexBuilder[*Model](ctx.QueryParams(), config)
}

func (m *Queries) GetOne(shown *Model, isPublic bool) error {
	c := &finder.ConfigShow{
		DB:       m.DB,
		QB:       m.QB,
		IsPublic: isPublic,
		Selects:  selects,
	}
	return finder.ShowOne(shown, c)
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
		Inserts: inserts,
		Selects: selects,
	}
	return finder.CreateOne(created, c)
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
		OptimisticLock: &finder.OptimisticLock{
			Name:  "updated_at",
			Value: updated.UpdatedAt,
		},
	}
	if slices.Contains(CoreKeys, updated.Key) {
		c.Input = &[]any{updated.Value}
		c.Inserts = &[]string{"value"}
	}
	return finder.UpdateOne(updated, c)
}

func (m *Queries) DeleteOne(deleted *Model, conn finder.Connection) error {
	c := &finder.ConfigDelete{
		DB:      conn,
		QB:      m.QB,
		Selects: selects,
	}
	return finder.DeleteOne(deleted, c)
}
