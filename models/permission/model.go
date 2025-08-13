package permission

import (
	"fmt"
	"net/url"

	"bitbucket.org/sadeemTechnology/backend-finder"
	validator "bitbucket.org/sadeemTechnology/backend-validator"
)

type Model struct {
	ID         int    `db:"id"          json:"id"`
	Method     string `db:"method"      json:"method"`
	Path       string `db:"path"        json:"path"`
	Action     string `db:"action"      json:"action"`
	Model      string `db:"model"       json:"model"`
	Scope      string `db:"scope"       json:"scope"`
	IsLoggable bool   `db:"is_loggable" json:"is_loggable"`
	IsVisible  bool   `db:"is_visible"  json:"is_visible"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return fmt.Sprintf("%d", m.ID)
}

func (m *Model) ModelName() string {
	return "permission"
}

func (m *Model) TableName() string {
	return "permissions"
}

func (m *Model) DefaultSearch() string {
	return "path"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		m.DefaultSearch(),
		"action",
		"model",
	}
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{}
}

// Initialize permission panics, the entire model is runtime/startup handled
func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	panic("shouldn't Initialize permissions")
}

// Utilities ------------------------------------------------------------------

func BuildMap(perms []Model) map[string][]string {
	m := make(map[string][]string)
	if len(perms) == 0 {
		return m
	}

	for i := range perms {
		p := perms[i]
		if _, found := m[p.Model]; !found {
			m[p.Model] = []string{p.Action + ":" + p.Scope}
		} else {
			m[p.Model] = append(m[p.Model], p.Action+":"+p.Scope)
		}
	}
	return m
}

func (m *Model) MergeAndValidate(v *validator.Validator) bool {
	v.AssignBool("is_loggable", &m.IsLoggable)
	v.AssignBool("is_visible", &m.IsVisible)

	v.ValidateModelSchema(m, v.Schema)
	return v.Valid()
}
