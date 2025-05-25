package token

import (
	"fmt"
	"net/url"
	"time"

	"app/model"

	"github.com/google/uuid"
	"github.com/m-row/finder"
	"github.com/m-row/pgtypes"
)

const (
	ScopeAdmin = "admin"
	ScopeOwn   = "own"

	TypeFCM = "fcm_token"
)

type Model struct {
	ID         int           `db:"id"          json:"id"`
	UserID     uuid.UUID     `db:"user_id"     json:"user_id"`
	TokenType  *string       `db:"token_type"  json:"token_type"`
	TokenValue *string       `db:"token_value" json:"token_value"`
	Data       pgtypes.JSONB `db:"data"        json:"data"`
	CreatedAt  *time.Time    `db:"created_at"  json:"created_at"`
	UpdatedAt  *time.Time    `db:"updated_at"  json:"updated_at"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return fmt.Sprintf("%d", m.ID)
}

func (m *Model) ModelName() string {
	return "token"
}

func (m *Model) TableName() string {
	return "tokens"
}

func (m *Model) DefaultSearch() string {
	return "token_value"
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
	if isInsert && m.ID == 0 {
		model.SelectSeqID(&m.ID, m.TableName(), conn)
	}
	return isInsert
}
