package wallet

import (
	"net/url"
	"time"

	"bitbucket.org/sadeemTechnology/backend-finder"

	"github.com/google/uuid"
)

const (
	ScopeAdmin    = "admin"
	ScopeCustomer = "customer"
)

type Model struct {
	ID             uuid.UUID `db:"id"               json:"id"`
	Credit         float64   `db:"credit"           json:"credit"`
	TrxCountCredit int       `db:"trx_count_credit" json:"trx_count_credit"`
	TrxCountDebit  int       `db:"trx_count_debit"  json:"trx_count_debit"`
	TrxTotalCredit float64   `db:"trx_total_credit" json:"trx_total_credit"`
	TrxTotalDebit  float64   `db:"trx_total_debit"  json:"trx_total_debit"`
	CreatedAt      time.Time `db:"created_at"       json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"       json:"updated_at"`
}

type MinimalModel struct {
	Credit         float64 `db:"credit"           json:"credit"`
	TrxCountCredit int     `db:"trx_count_credit" json:"trx_count_credit"`
	TrxCountDebit  int     `db:"trx_count_debit"  json:"trx_count_debit"`
	TrxTotalCredit float64 `db:"trx_total_credit" json:"trx_total_credit"`
	TrxTotalDebit  float64 `db:"trx_total_debit"  json:"trx_total_debit"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) ModelName() string {
	return "wallet"
}

func (m *Model) TableName() string {
	return "wallets"
}

func (m *Model) DefaultSearch() string {
	return "id"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{"id"}
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{}
}

func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	panic("shouldn't initialize wallet")
}
