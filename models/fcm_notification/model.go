package fcm_notification

import (
	"net/url"
	"time"

	"app/model"
	"app/models/user"
	"github.com/google/uuid"
	"github.com/m-row/finder"
	"github.com/m-row/pgtypes"
	"github.com/m-row/validator"
)

const (
	TokenTypeStandard = "fcm_token"
)

var AllTokenTypes = []string{
	TokenTypeStandard,
}

var (
	ScopeAdmin = "admin"
	ScopeUser  = "user"
)

type Model struct {
	ID        uuid.UUID        `db:"id"         json:"id"`
	Title     string           `db:"title"      json:"title"`
	Body      string           `db:"body"       json:"body"`
	Topic     *string          `db:"topic"      json:"topic"`
	IsSent    bool             `db:"is_sent"    json:"is_sent"`
	SendAt    *time.Time       `db:"send_at"    json:"send_at"`
	Response  *string          `db:"response"   json:"response"`
	Data      pgtypes.JSONBStr `db:"data"       json:"data"`
	SenderID  *uuid.UUID       `db:"sender_id"  json:"-"`
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt time.Time        `db:"updated_at" json:"updated_at"`

	Sender *user.MinimalModel `db:"sender" json:"sender"`
}

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) ModelName() string {
	return "fcm_notification"
}

func (m *Model) TableName() string {
	return "fcm_notifications"
}

func (m *Model) DefaultSearch() string {
	return "title"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		"title",
	}
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

func (m *Model) MergeAndValidate(v *validator.Validator) bool {
	isInsert := m.Initialize(v.Data.Values, v.Conn)

	v.AssignString("title", &m.Title, 0, 500)
	v.AssignString("body", &m.Body, 0, 500)
	m.Topic = v.AssignString("topic", m.Topic, 0, 500)
	v.AssignBool("is_sent", &m.IsSent)

	v.AssignTimestamp("send_at", m.SendAt)

	if isInsert && m.SendAt == nil {
		m.IsSent = true
	}
	if m.Sender != nil {
		if m.Sender.ID != nil {
			m.SenderID = m.Sender.ID
		}
	}
	v.UnmarshalInto("data", &m.Data)

	v.ValidateModelSchema(m, m.TableName(), v.Schema)
	return v.Valid()
}
