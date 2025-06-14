package user_notification

import (
	"net/url"
	"time"

	"bitbucket.org/sadeemTechnology/backend-finder"
	"bitbucket.org/sadeemTechnology/backend-model"
	"bitbucket.org/sadeemTechnology/backend-pgtypes"

	"app/models/user"

	"github.com/google/uuid"
)

const (
	ScopeAdmin    = "admin"
	ScopeCustomer = "customer"
)

type Model struct {
	ID         uuid.UUID         `db:"id"          json:"id"`
	UserID     *uuid.UUID        `db:"user_id"     json:"-"`
	IsRead     bool              `db:"is_read"     json:"is_read"`
	IsNotified bool              `db:"is_notified" json:"is_notified"`
	Title      string            `db:"title"       json:"title"`
	Body       string            `db:"body"        json:"body"`
	Response   *string           `db:"response"    json:"response"`
	Data       pgtypes.JSONBStr  `db:"data"        json:"data"`
	CreatedAt  time.Time         `db:"created_at"  json:"created_at"`
	UpdatedAt  time.Time         `db:"updated_at"  json:"updated_at"`
	User       user.MinimalModel `db:"user"        json:"user"`
}

// ---------------------------------------
// Implement the Model interface
// ---------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) ModelName() string {
	return "user_notification"
}

func (m *Model) TableName() string {
	return "user_notifications"
}

func (m *Model) DefaultSearch() string {
	return "title"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		"title",
		"body",
	}
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{}
} // Model interface end

func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	isInsert := m.CreatedAt.Equal(time.Time{})
	if isInsert && m.ID == uuid.Nil {
		model.InputOrNewUUID(&m.ID, v)
	}
	return isInsert
}
