package wallet_transaction

import (
	"net/url"
	"time"

	"github.com/m-row/model"
	"github.com/m-row/pgtypes"
	"github.com/m-row/validator"

	"github.com/m-row/finder"

	"github.com/google/uuid"
)

const (
	ScopeAdmin = "admin"
	ScopeOwn   = "own"
)

type Model struct {
	ID uuid.UUID `db:"id" json:"id"`

	// WalletID references users.id foreign field
	WalletID      uuid.UUID  `db:"wallet_id"       json:"wallet_id"`
	RechargedByID *uuid.UUID `db:"recharged_by_id" json:"-"`
	Type          TypeValue  `db:"type"            json:"type"`
	Amount        float64    `db:"amount"          json:"amount"`
	PaymentMethod *string    `db:"payment_method"  json:"payment_method"`

	PaymentReference *string `db:"payment_reference" json:"payment_reference"`

	Notes         *string       `db:"notes"          json:"notes"`
	IsConfirmed   bool          `db:"is_confirmed"   json:"is_confirmed"`
	TLyncURL      *string       `db:"tlync_url"      json:"tlync_url"`
	TLyncResponse pgtypes.JSONB `db:"tlync_response" json:"tlync_response"`
	CreatedAt     time.Time     `db:"created_at"     json:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at"     json:"updated_at"`

	SerialNumber int     `db:"-" json:"-"`
	Phone        *string `db:"-" json:"-"`

	User        WalletUser `db:"user"         json:"user"`
	RechargedBy WalletUser `db:"recharged_by" json:"recharged_by"`
}

type WalletUser struct {
	ID   *uuid.UUID `db:"id"   json:"id"`
	Name *string    `db:"name" json:"name"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) ModelName() string {
	return "wallet_transaction"
}

func (m *Model) TableName() string {
	return "wallet_transactions"
}

func (m *Model) DefaultSearch() string {
	return "notes"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		m.DefaultSearch(),
		"type",
		"payment_method",
		"payment_reference",
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
	m.PaymentMethod = v.
		AssignString("payment_method", m.PaymentMethod, 0, 500)
	m.PaymentReference = v.
		AssignString("payment_reference", m.PaymentReference, 0, 500)
	m.Notes = v.AssignString("notes", m.Notes, 0, 500)

	validator.AssignENUM(v, "type", &m.Type)

	v.AssignFloat("amount", &m.Amount)
	v.UnmarshalInto("user", &m.User, ScopeAdmin)
	v.AssignBool("is_confirmed", &m.IsConfirmed, ScopeAdmin)

	v.ValidateModelSchema(m, m.TableName(), v.Schema)
	return v.Valid()
}

func (m *Model) MergeTransfer(v *validator.Validator) bool {
	m.Initialize(v.Data.Values, v.Conn)
	v.AssignFloat("amount", &m.Amount)
	v.Check(m.Amount > 0, "amount", v.T.ValidateMustBeGtZero())

	if len(v.Data.Get("phone")) == 6 {
		v.AssignInt("phone", &m.SerialNumber)
		if m.SerialNumber < 100000 || m.SerialNumber > 999999 {
			v.Check(false, "serial_number", "Invalid serial number")
		}
	} else {
		m.MergePhone(v)
	}

	return v.Valid()
}
