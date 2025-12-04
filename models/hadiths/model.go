package hadiths

import (
	"app/models/consts"
	"app/models/user"
	"net/url"
	"time"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	model "bitbucket.org/sadeemTechnology/backend-model"
	category "bitbucket.org/sadeemTechnology/backend-model-category"
	validator "bitbucket.org/sadeemTechnology/backend-validator"
	"github.com/google/uuid"
)

var (
	ScopeAdmin = "admin"
	ScopeOwn   = "own"
)

type Model struct {
	ID          uuid.UUID             `db:"id"         json:"id"`
	Text        string                `db:"text" json:"text"`
	Source      string                `db:"source" json:"source"`
	Topic       string                `db:"topic" json:"topic"`
	CategoryID  uuid.UUID             `db:"category_id" json:"-"`
	CreatedAt   time.Time             `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time             `db:"updated_at" json:"updated_at"`
	CreatedByID uuid.UUID             `db:"created_by_id" json:"-"`
	CreatedBy   user.MinimalModel     `db:"created_by" json:"created_by"`
	Category    category.MinimalModel `db:"category" json:"category"`
}

type MinimalModel struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Model methods --------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) ModelName() string {
	return "hadith"
}

func (m *Model) TableName() string {
	return "hadiths"
}

func (m *Model) DefaultSearch() string {
	return "text"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{m.DefaultSearch(), "Topic", "source"}
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
	_ = m.Initialize(v.Data.Values, v.Conn)

	v.UnmarshalInto("category", &m.Category)
	if m.Category.ID != nil {
		v.CategoryValidator(m.Category.ID, "categories.id", consts.CategoryHadithID)
		m.CategoryID = *m.Category.ID
	} else {
		m.CategoryID = uuid.Nil
	}

	v.AssignString("text", &m.Text, 1, 255)
	v.AssignString("source", &m.Source, 1, 255)
	v.AssignString("topic", &m.Topic, 1, 255)
	return v.Valid()
}
