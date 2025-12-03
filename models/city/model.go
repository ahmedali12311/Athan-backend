//nolint:lll
package city

import (
	"net/url"
	"time"

	pgtypes "bitbucket.org/sadeemTechnology/backend-pgtypes"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	model "bitbucket.org/sadeemTechnology/backend-model"
	validator "bitbucket.org/sadeemTechnology/backend-validator"
	"github.com/google/uuid"
)

const (
	ScopeAdmin = "admin"
)

type Model struct {
	ID         uuid.UUID     `db:"id"                        json:"id"                        csv:"-"`
	Name       string        `db:"name"                      json:"name"                      csv:"name"`
	IsDisabled bool          `db:"is_disabled"               json:"is_disabled"               csv:"is_disabled"`
	Location   pgtypes.Point `db:"location"                  json:"location"                  csv:"-"`
	CreatedAt  time.Time     `db:"created_at"                json:"created_at"                csv:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"                json:"updated_at"                csv:"updated_at"`
}

type MinimalModel struct {
	ID   *uuid.UUID `db:"id"   json:"id"   csv:"-"`
	Name *string    `db:"name" json:"name" csv:"name"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) ModelName() string {
	return "city"
}

func (m *Model) TableName() string {
	return "cities"
}

func (m *Model) DefaultSearch() string {
	return "name"
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
	return &[]finder.RelationField{}
}

func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	isInsert := m.CreatedAt.Equal(time.Time{})
	if isInsert && m.ID == uuid.Nil {
		model.InputOrNewUUID(&m.ID, v)
	}
	return isInsert
}

// Utilities ------------------------------------------------------------------

func (m *Model) MergeAndValidate(v *validator.Validator) bool {
	m.Initialize(v.Data.Values, v.Conn)
	v.AssignString("name", &m.Name, 1, 50)
	v.AssignBool("is_disabled", &m.IsDisabled)

	v.UnmarshalInto("location", &m.Location)

	v.ValidateModelSchema(m, v.Schema)
	return v.Valid()
}
