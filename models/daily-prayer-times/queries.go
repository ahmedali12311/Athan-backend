package dailyprayertimes

import (
	"time"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	model "bitbucket.org/sadeemTechnology/backend-model"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	selects = &[]string{
		"daily_prayer_times.*",
		"cities.id as \"city.id\"",
		"cities.name as \"city.name\"",
		"daily_prayer_times.fajr_first_time",
		"daily_prayer_times.fajr_second_time",
		"daily_prayer_times.sunrise_time",
		"daily_prayer_times.dhuhr_time",
		"daily_prayer_times.asr_time",
		"daily_prayer_times.maghrib_time",
		"daily_prayer_times.isha_time",
	}

	inserts = &[]string{
		"id",
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
		"cities ON daily_prayer_times.city_id = cities.id",
	}

	GroupBys = &[]string{
		"daily_prayer_times.id",
		"cities.id",
	}
)

type WhereScope struct {
	IsAdmin bool
	CityID  *uuid.UUID
	Day     *int
	Month   *int
	Date    *time.Time
}

func getJoins(ws *WhereScope) *[]string {
	return baseJoins
}

func buildInput(m *Model) (*[]any, error) {
	input := &[]any{
		m.ID,
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

func getSelects() *[]string {
	return selects
}

func wheres(ws *WhereScope) *[]squirrel.Sqlizer {
	w := []squirrel.Sqlizer{}

	if !ws.IsAdmin {
		w = append(w, squirrel.Eq{"daily_prayer_times.is_deleted": false})
	}

	if ws.CityID != nil {
		w = append(w, squirrel.Eq{"daily_prayer_times.city_id": *ws.CityID})
	}

	if ws.Day != nil {
		w = append(w, squirrel.Eq{"daily_prayer_times.day": *ws.Day})
	}
	if ws.Month != nil {
		w = append(w, squirrel.Eq{"daily_prayer_times.month": *ws.Month})
	}
	if ws.Date != nil {
		w = append(w, squirrel.Eq{"daily_prayer_times.day": ws.Date.Day()})
		w = append(w, squirrel.Eq{"daily_prayer_times.month": ws.Date.Month()})
	}

	return &w
}

func getOrderBy(ws *WhereScope) string {
	return "daily_prayer_times.city_id ASC, daily_prayer_times.month ASC, daily_prayer_times.day ASC"
}

func (m *Queries) GetAll(
	ctx echo.Context,
	ws *WhereScope,
) (*finder.IndexResponse[*Model], error) {

	cfg := &finder.ConfigIndex{
		DB:           m.DB,
		QB:           m.QB,
		PGInfo:       m.PGInfo,
		Joins:        getJoins(ws),
		Selects:      getSelects(),
		Wheres:       wheres(ws),
		GroupBys:     GroupBys,
		OverrideSort: getOrderBy(ws),
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
		Joins:   getJoins(ws),
		Wheres:  wheres(ws),
		Selects: getSelects(),
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
		Inserts: inserts,
		Selects: selects,
		Joins:   baseJoins,
	}

	return finder.CreateOne(created, c)
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
	whereConditions := wheres(ws)

	*whereConditions = append(*whereConditions, squirrel.Eq{"daily_prayer_times.id": updated.ID})

	c := &finder.ConfigUpdate{
		DB:      conn,
		QB:      m.QB,
		Joins:   getJoins(ws),
		Wheres:  whereConditions,
		Input:   input,
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
	ws *WhereScope,
	conn finder.Connection,
) error {
	c := &finder.ConfigDelete{
		DB:      conn,
		QB:      m.QB,
		Joins:   getJoins(ws),
		Selects: getSelects(),
	}
	return finder.DeleteOne(deleted, c)
}
