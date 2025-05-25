package user

import (
	"fmt"
	"hash/crc32"
	"net/url"
	"time"

	"app/config"
	"app/model"
	"github.com/m-row/validator"

	"github.com/google/uuid"
	"github.com/m-row/finder"
	"github.com/m-row/pgtypes"
)

const (
	ScopeOwn   = "own"
	ScopeAdmin = "admin"

	SuperAdminID = "322f3e97-4e7e-4c2e-a765-1c0ce517f2f8"
)

type Model struct {
	ID           uuid.UUID      `db:"id"            json:"id"`
	Ref          string         `db:"ref"           json:"ref"`
	Name         *string        `db:"name"          json:"name"`
	Phone        *string        `db:"phone"         json:"phone"`
	Email        *string        `db:"email"         json:"email"`
	Password     password       `db:"-"             json:"-"`
	PasswordHash *[]byte        `db:"password_hash" json:"-"`
	Img          *string        `db:"img"           json:"img"`
	Thumb        *string        `db:"thumb"         json:"thumb"`
	Gender       *GenderValue   `db:"gender"        json:"gender"`
	Details      *string        `db:"details"       json:"details"`
	Birthdate    *string        `db:"birthdate"     json:"birthdate"`
	Location     *pgtypes.Point `db:"location"      json:"location"`
	IsAnon       bool           `db:"is_anon"       json:"is_anon"`
	IsNotifiable bool           `db:"is_notifiable" json:"is_notifiable"`
	IsDisabled   bool           `db:"is_disabled"   json:"is_disabled"`
	IsConfirmed  bool           `db:"is_confirmed"  json:"is_confirmed"`
	IsDeleted    bool           `db:"is_deleted"    json:"is_deleted"`
	IsVerified   bool           `db:"is_verified"   json:"is_verified"`
	LastRef      int            `db:"last_ref"      json:"last_ref"`

	Pin       *string    `db:"pin"        json:"-"`
	PinExpiry *time.Time `db:"pin_expiry" json:"-"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	Roles       *[]int `db:"roles"       json:"roles,omitempty"`
	Permissions *[]int `db:"permissions" json:"permissions,omitempty"`
}

type MinimalModel struct {
	ID         *uuid.UUID `db:"id"          json:"id"`
	Name       *string    `db:"name"        json:"name"`
	Phone      *string    `db:"phone"       json:"phone"`
	Email      *string    `db:"email"       json:"email"`
	DocumentID *uuid.UUID `db:"document_id" json:"-"`
}

type UserRole struct {
	UserID uuid.UUID `db:"user_id"`
	RoleID int       `db:"role_id"`
}

type MinimalModelLessonComment struct {
	ID    *uuid.UUID `db:"id"    json:"id"`
	Name  *string    `db:"name"  json:"name"`
	Phone *string    `db:"phone" json:"phone"`
	Email *string    `db:"email" json:"email"`
	Img   *string    `db:"img"   json:"img"`
	Thumb *string    `db:"thumb" json:"thumb"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) ModelName() string {
	return "user"
}

func (m *Model) TableName() string {
	return "users"
}

func (m *Model) DefaultSearch() string {
	return "name"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		m.DefaultSearch(),
		"email",
		"phone",
	}
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{
		{
			Table: "roles",
			Join: &finder.Join{
				From: "users.id",
				To:   "roles.id",
			},
			Through: &finder.Through{
				Table: "user_roles",
				Join: &finder.Join{
					From: "user_roles.user_id",
					To:   "user_roles.role_id",
				},
			},
		},
	}
}

// Initialize generates a uuid and crc32 checksum hash for the user based on
// that uuid, must be called for user model
func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	isInsert := m.CreatedAt.Equal(time.Time{})
	if isInsert || m.ID == uuid.Nil {
		// 1
		model.InputOrNewUUID(&m.ID, v)

		// 2 // TODO: preload this
		crc32q := crc32.MakeTable(config.CRC32Poly)
		m.Ref = fmt.Sprintf(
			"%08x",
			crc32.Checksum([]byte(m.ID.String()), crc32q),
		)
	}
	return isInsert
}

// Has Image ------------------------------------------------------------------

func (m *Model) GetImg() *string {
	return m.Img
}

func (m *Model) SetImg(name *string) {
	m.Img = name
}

func (m *Model) GetThumb() *string {
	return m.Thumb
}

func (m *Model) SetThumb(name *string) {
	m.Thumb = name
}

// Utilities ------------------------------------------------------------------

func (m *Model) MergeAndValidate(v *validator.Validator) bool {
	isInsert := m.Initialize(v.Data.Values, v.Conn)

	if err := v.AssignImage("img", m, false); err != nil {
		v.Check(false, "img", err.Error())
	}

	m.Name = v.AssignString("name", m.Name, 0, 100)
	m.Email = v.AssignString("email", m.Email, 0, 100)
	m.Details = v.AssignString("details", m.Details, 0, 3000)
	m.Birthdate = v.AssignDate("birthdate", m.Birthdate)

	v.AssignBool("is_disabled", &m.IsDisabled)
	v.AssignBool("is_notifiable", &m.IsNotifiable)
	v.UnmarshalInto("location", &m.Location)

	// Enums ------------------------------------------------------------------
	m.Gender = validator.AssignENUM(v, "gender", m.Gender)

	// Special merge logic ----------------------------------------------------
	m.MergePhone(v)
	m.MergeEmailPassword(v, isInsert, false)

	// sets the user as anonymous if no email or phone is entered -------------
	if isInsert && m.Phone == nil && m.Email == nil {
		m.IsAnon = true
	}
	// only admin allowed to modify -------------------------------------------
	// protects superadmin from suicide
	if m.ID.String() != SuperAdminID {
		v.UnmarshalInto("roles", &m.Roles, "admin")
		v.UnmarshalInto("permissions", &m.Permissions, "admin")
		if m.Roles == nil && !m.IsAnon {
			v.Check(false, "roles", "this should never happen")
			return false
		}
		if !m.IsAnon {
			v.Check(len(*m.Roles) != 0, "roles", v.T.ValidateNotEmptyRoles())
		}
	}

	v.ValidateModelSchema(m, m.TableName(), v.Schema)
	return v.Valid()
}

func (m *Model) MergeOTPCreate(
	v *validator.Validator,
	pin *string,
	exp *time.Time,
) {
	bd := "2001-01-01"
	g := Male

	m.Initialize(v.Data.Values, v.Conn)
	m.MergePhone(v)
	m.Pin = pin
	m.PinExpiry = exp

	m.Name = nil
	m.Email = nil
	m.Gender = &g
	m.Birthdate = &bd
	m.IsVerified = false
	m.Location = nil // TODO: test this
	// TODO: add lowest common public,own scopes
}
