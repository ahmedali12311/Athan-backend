package dailyprayertimes

import (
	"time"

	"app/models/city"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	model "bitbucket.org/sadeemTechnology/backend-model"
	validator "bitbucket.org/sadeemTechnology/backend-validator"

	"net/url"

	"github.com/google/uuid"
)

const (
	ScopeAdmin = "admin"
)

type Model struct {
	ID     uuid.UUID         `db:"id" json:"id"`
	CityID *uuid.UUID        `db:"city_id" json:"-"`
	City   city.MinimalModel `db:"city" json:"city"`
	Day    int               `db:"day" json:"day"`
	Month  int               `db:"month" json:"month"`

	FajrFirstTime  time.Time `db:"fajr_first_time" json:"fajr_first_time"`
	FajrSecondTime time.Time `db:"fajr_second_time" json:"fajr_second_time"`
	SunriseTime    time.Time `db:"sunrise_time" json:"sunrise_time"`
	DhuhrTime      time.Time `db:"dhuhr_time" json:"dhuhr_time"`
	AsrTime        time.Time `db:"asr_time" json:"asr_time"`
	MaghribTime    time.Time `db:"maghrib_time" json:"maghrib_time"`
	IshaTime       time.Time `db:"isha_time" json:"isha_time"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted bool      `db:"is_deleted" json:"is_deleted"`
}

type MinimalModel struct {
	ID    *uuid.UUID        `db:"id" json:"id"`
	City  city.MinimalModel `db:"city" json:"city"`
	Day   int               `db:"day" json:"day"`
	Month int               `db:"month" json:"month"`
}

// --- Model Interface Methods ---

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) ModelName() string {
	return "daily_prayer_times"
}

func (m *Model) TableName() string {
	return "daily_prayer_times"
}

func (m *Model) DefaultSearch() string {
	return "city_id"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		m.DefaultSearch(),
	}
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{
		{
			Table: "cities",
			Join: &finder.Join{
				From: "daily_prayer_times.city_id",
				To:   "cities.id",
			},
		},
	}
}

func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	isInsert := m.CreatedAt.Equal(time.Time{})
	if isInsert || m.ID == uuid.Nil {
		model.InputOrNewUUID(&m.ID, v)
	}
	return isInsert
}

// --- Utilities ---

func (m *Model) MergeAndValidate(v *validator.Validator) bool {
	_ = m.Initialize(v.Data.Values, v.Conn)

	v.AssignInt("day", &m.Day)
	v.AssignInt("month", &m.Month)

	v.Check(m.Day >= 1 && m.Day <= 31, "day", "Day must be between 1 and 31")
	v.Check(m.Month >= 1 && m.Month <= 12, "month", "Month must be between 1 and 12")

	v.AssignClock("fajr_first_time", &m.FajrFirstTime)
	v.AssignClock("fajr_second_time", &m.FajrSecondTime)
	v.AssignClock("sunrise_time", &m.SunriseTime)
	v.AssignClock("dhuhr_time", &m.DhuhrTime)
	v.AssignClock("asr_time", &m.AsrTime)
	v.AssignClock("maghrib_time", &m.MaghribTime)
	v.AssignClock("isha_time", &m.IshaTime)

	// Handle city relation
	v.UnmarshalInto("city", &m.City)
	v.UUIDExistsInDB(m.City.ID, "city.id", "id", "cities", true)
	if m.City.ID != nil {
		m.CityID = m.City.ID
	}

	v.ValidateModelSchema(m, v.Schema)
	return v.Valid()
}
