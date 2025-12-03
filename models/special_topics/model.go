package special_topics

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
	ID          uuid.UUID             `db:"id" json:"id"`
	Topic       string                `db:"topic" json:"topic"`
	Content     string                `db:"content" json:"content"`
	CategoryID  uuid.UUID             `db:"category_id" json:"-"`
	Img         *string               `db:"img"           json:"img"`
	Thumb       *string               `db:"thumb"         json:"thumb"`
	CreatedAt   time.Time             `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time             `db:"updated_at" json:"updated_at"`
	Category    category.MinimalModel `db:"category" json:"category"`
	CreatedByID uuid.UUID             `db:"created_by_id" json:"-"`
	CreatedBy   user.MinimalModel     `db:"created_by" json:"created_by"`
}

type MinimalModel struct {
	ID        uuid.UUID `db:"id" json:"id"`
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
	return "special_topic"
}

func (m *Model) TableName() string {
	return "special_topics"
}

func (m *Model) DefaultSearch() string {
	return "topic"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{"topic", "content"}
}
func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{}
}

func (m *Model) Initialize(v url.Values, conn finder.Connection) bool {
	isInsert := m.CreatedAt.Equal(time.Time{})
	if isInsert || m.ID == uuid.Nil {
		model.InputOrNewUUID(&m.ID, v)
	}
	return isInsert
}

func (m *Model) MergeAndValidate(v *validator.Validator) bool {
	_ = m.Initialize(v.Data.Values, v.Conn)

	if err := v.AssignImage("img", m, false); err != nil {
		v.Check(false, "img", err.Error())
	}

	v.UnmarshalInto("category", m.Category)
	if m.Category.ID != nil {
		v.CategoryValidator(m.Category.ID, "categories.id", consts.CategorySpecialTopicID)
		m.CategoryID = *m.Category.ID
	} else {
		m.CategoryID = uuid.Nil
	}

	v.AssignString("topic", &m.Topic, 1, 255)
	v.AssignString("content", &m.Content, 1, 255)
	return v.Valid()
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
