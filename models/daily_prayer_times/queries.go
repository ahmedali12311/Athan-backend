package daily_prayer_times

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		"daily_prayer_timeses.*",
		"c.id as \"city.id\"",
		"c.name as \"city.name\"",
	}

	inserts = &[]string{
		"city_id",
		"day",
		"month",
		"fajr_first_time",
		"fajr_second_time",
		"sunrise_time",
		"dhuhr_time",
		"asr_time",
		"maghrib_time",
		"isha_time",
	}
	baseJoins = &[]string{
		"cities ON daily_prayer_timeses.city_id = cities.id",
	}
)

type CityDue struct {
	Name string `db:"city_name"`
	Time string `db:"prayer_time"`
}

func buildInput(m *Model) (*[]any, error) {
	input := &[]any{
		m.CityID,
		m.Day,
		m.Month,
		m.FajrFirstTime,
		m.FajrSecondTime,
		m.SunriseTime,
		m.DhuhrTime,
		m.AsrTime,
		m.MaghribTime,
		m.IshaTime,
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

func (m *Queries) GetDueCities(
	ctx context.Context,
	prayerTimeColumn string,
	day, month int,
	targetTime string,
) ([]CityDue, error) {

	query := fmt.Sprintf(`
        SELECT c.name AS city_name, dpt.%s AS prayer_time
        FROM daily_prayer_times dpt
        JOIN cities c ON dpt.city_id = c.id
        WHERE dpt.day = $1 
        AND dpt.month = $2
        AND dpt.%s = $3
    `, prayerTimeColumn, prayerTimeColumn)
	var results []CityDue

	err := m.DB.SelectContext(ctx, &results, query, day, month, targetTime)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return []CityDue{}, nil
	}
	return results, err
}
