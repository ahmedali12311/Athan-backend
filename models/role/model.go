package role

import (
	"fmt"
	"net/url"

	"github.com/m-row/model"
	"github.com/m-row/validator"

	"github.com/google/uuid"
	"github.com/m-row/finder"
)

const (
	ScopeAdmin  = "admin"
	ScopePublic = "public"

	AdminRole    = 2
	MerchantRole = 3
	CustomerRole = 4
)

type Model struct {
	ID          int    `db:"id"          json:"id"`
	Name        string `db:"name"        json:"name"`
	Permissions []int  `db:"permissions" json:"permissions,omitempty"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return fmt.Sprintf("%d", m.ID)
}

func (m *Model) ModelName() string {
	return "role"
}

func (m *Model) TableName() string {
	return "roles"
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
	isInsert := m.ID == 0
	if isInsert {
		model.SelectSeqID(&m.ID, m.TableName(), conn)
	}
	return isInsert
}

// utilties -------------------------------------------------------------------

func (m *Model) MergeAndValidate(v *validator.Validator) bool {
	m.Initialize(v.Data.Values, v.Conn)

	v.AssignString("name", &m.Name, 1, 50)
	v.UnmarshalInto("permissions", &m.Permissions)

	v.ValidateModelSchema(m, m.TableName(), v.Schema)
	return v.Valid()
}

func (m *Model) ValidateUserRole(
	v *validator.Validator,
	userID *uuid.UUID,
	roleID *int,
) bool {
	v.AssignInt("role_id", roleID)
	v.IDExistsInDB(roleID, "role_id", "id", "roles", true)

	v.AssignUUID("user_id", "id", "users", userID, true)
	v.UUIDExistsInDB(userID, "user_id", "id", "users", true)

	return !v.Valid()
}
