package token

import (
	"app/model"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder"
)

var selects = &[]string{
	"tokens.id",
	"tokens.user_id",
	"tokens.token_type",
	"tokens.token_value",
	"tokens.data::json",
	"tokens.created_at",
	"tokens.updated_at",
}

type Queries struct {
	*model.Dependencies
}

func New(d *model.Dependencies) *Queries {
	return &Queries{d}
}

func wheres(userID *uuid.UUID) *[]squirrel.Sqlizer {
	w := &[]squirrel.Sqlizer{}
	if userID != nil {
		*w = append(*w, squirrel.Expr("tokens.user_id=?", userID))
	}
	return w
}

func (m *Queries) GetAll(
	ctx echo.Context,
	userID *uuid.UUID,
) (*finder.IndexResponse[*Model], error) {
	config := &finder.ConfigIndex{
		DB:      m.DB,
		QB:      m.QB,
		PGInfo:  m.PGInfo,
		Wheres:  wheres(userID),
		Selects: selects,
		GroupBys: &[]string{
			"tokens.id",
		},
	}
	return finder.IndexBuilder[*Model](ctx.QueryParams(), config)
}

func (m *Queries) GetOne(shown *Model, userID *uuid.UUID) error {
	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Selects: selects,
		Wheres:  wheres(userID),
	}
	return finder.ShowOne(shown, c)
}
