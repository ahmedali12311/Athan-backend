package daily_prayer_times

import (
	"app/models/city"
	"net/url"
	"time"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	model "bitbucket.org/sadeemTechnology/backend-model"
	validator "bitbucket.org/sadeemTechnology/backend-validator"
	"github.com/google/uuid"
)

var (
	ScopeAdmin = "admin"
	ScopeOwn   = "own"
)

type Model struct {
	ID             uuid.UUID          `db:"id" json:"id"`
	CityID         uuid.UUID          `db:"city_id" json:"-"`
	Day            int                `db:"day" json:"day"`
	Month          int                `db:"month" json:"month"`
	FajrFirstTime  string             `db:"fajr_first_time" json:"fajr_first_time"`
	FajrSecondTime string             `db:"fajr_second_time" json:"fajr_second_time"`
	SunriseTime    string             `db:"sunrise_time" json:"sunrise_time"`
	DhuhrTime      string             `db:"dhuhr_time" json:"dhuhr_time"`
	AsrTime        string             `db:"asr_time" json:"asr_time"`
	MaghribTime    string             `db:"maghrib_time" json:"maghrib_time"`
	IshaTime       string             `db:"isha_time" json:"isha_time"`
	CreatedAt      time.Time          `db:"created_at" json:"created_at"`
	City           *city.MinimalModel `db:"city" json:"city"`
}

type MinimalModel struct {
	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Model methods --------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) ModelName() string {
	return "daily_prayer_times"
}

func (m *Model) TableName() string {
	return "daily_prayer_timeses"
}

func (m *Model) DefaultSearch() string {
	return "name"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{"name", "description"}
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{
		{
			Table: "citieses",
			Join: &finder.Join{
				From: "daily_prayer_timeses.city_id",
				To:   "citieses.id",
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

func (m *Model) MergeAndValidate(v *validator.Validator) bool {
	_ = m.Initialize(v.Data.Values, v.Conn)

	v.UnmarshalInto("city", m.City)
	if m.City != nil && m.City.ID != nil {
		v.UUIDExistsInDB(m.City.ID, "city_id", "id", "cities", true)
		m.CityID = *m.City.ID
	} else {
		v.Check(false, "city", "must enter a city!")
	}

	v.AssignInt("day", &m.Day)
	v.AssignInt("month", &m.Month)
	return v.Valid()
}
