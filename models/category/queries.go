package category

import (
	"context"

	"app/config"
	"app/model"
	"app/pkg/sorter"

	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder"
)

var inserts = &[]string{
	"id",
	"name",
	"parent_id",
	"super_parent_id",
	"sort",
	"depth",
	"is_disabled",
	"is_featured",
	"img",
	"thumb",
}

func buildInput(category *Model) (*[]any, error) {
	input := &[]any{
		category.ID,
		category.Name,
		category.Parent.ID,
		category.SuperParent.ID,
		category.Sort,
		category.Depth,
		category.IsDisabled,
		category.IsFeatured,
		category.Img,
		category.Thumb,
	}

	if len(*input) != len(*inserts) {
		return nil, finder.ErrInputLengthMismatch(input, inserts)
	}

	return input, nil
}

func joins(alias string) *[]string {
	return &[]string{
		"categories as p ON " + alias + ".parent_id = p.id",
		"categories as sp ON " + alias + ".super_parent_id = sp.id",
	}
}

func selects(alias string) *[]string {
	return &[]string{
		alias + ".*",
		config.SQLSelectURLPath(alias, "img", "\"img\""),
		config.SQLSelectURLPath(alias, "thumb", "\"thumb\""),

		"p.id as \"parent.id\"",
		"p.name as \"parent.name\"",

		"sp.id as \"super_parent.id\"",
		"sp.name as \"super_parent.name\"",
	}
}

type WhereScope struct {
	IsPublic,
	ShowArchived bool
	SortBeforeUpdate int
}

func wheres(_ string, _ *WhereScope) *[]squirrel.Sqlizer {
	w := []squirrel.Sqlizer{}
	// if !ws.ShowArchived {
	// 	w = append(w, squirrel.Expr(alias+".is_archived=false"))
	// }
	return &w
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
		IsPublic: ws.IsPublic,
		Joins:    joins("categories"),
		Wheres:   wheres("categories", ws),
		Selects:  selects("categories"),
		GroupBys: &[]string{
			"categories.id",
			"p.id",
			"sp.id",
		},
	}
	indexResponse, err := finder.IndexBuilder[*Model](ctx.QueryParams(), cfg)
	if err != nil {
		return nil, err
	}
	return indexResponse, nil
}

func (m *Queries) GetOne(
	shown *Model,
	ws *WhereScope,
) error {
	c := &finder.ConfigShow{
		DB:       m.DB,
		QB:       m.QB,
		IsPublic: ws.IsPublic,
		Joins:    joins("categories"),
		Wheres:   wheres("categories", ws),
		Selects:  selects("categories"),
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
		DB:         conn,
		QB:         m.QB,
		Input:      input,
		Inserts:    inserts,
		TableAlias: "c1",
		Joins:      joins("c1"),
		Selects:    selects("c1"),
	}
	if err := sorter.AdjustSort(
		sorter.Create,
		created,
		conn,
		m.QB,
	); err != nil {
		return err
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
		DB:         conn,
		QB:         m.QB,
		Input:      input,
		Inserts:    inserts,
		TableAlias: "c1",
		Joins:      joins("c1"),
		Wheres:     wheres("c1", ws),
		Selects:    selects("c1"),
		OptimisticLock: &finder.OptimisticLock{
			Name:  "updated_at",
			Value: updated.UpdatedAt,
		},
	}
	if err := sorter.AdjustSort(
		sorter.Update,
		updated,
		conn,
		m.QB,
	); err != nil {
		return err
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
		DB:         m.DB,
		QB:         m.QB,
		TableAlias: "c1",
		Joins:      joins("c1"),
		Wheres:     wheres("c1", ws),
		Selects:    selects("c1"),
	}
	if err := sorter.AdjustSort(
		sorter.Delete,
		deleted,
		conn,
		m.QB,
	); err != nil {
		return err
	}
	if err := finder.DeleteOne(deleted, c); err != nil {
		return err
	}
	return nil
}

func (m *Queries) HasChildren(category *Model) (bool, error) {
	count := 0
	query, args, err := m.QB.
		Select("COUNT(*)").
		From("categories").
		Where("parent_id = ?", category.ID).
		ToSql()
	if err != nil {
		return false, err
	}

	if err := m.DB.GetContext(
		context.Background(),
		&count,
		query,
		args...,
	); err != nil {
		return count > 0, err
	}
	return count > 0, nil
}
