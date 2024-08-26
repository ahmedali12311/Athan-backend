package setting

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"time"

	"app/model"
	"app/pkg/validator"

	"github.com/ahmedalkabir/finder"
)

const (
	ScopeAdmin  = "admin"
	ScopePublic = "public"
)

type (
	Model struct {
		ID         int       `db:"id"          json:"id"`
		Key        string    `db:"key"         json:"key"`
		Value      string    `db:"value"       json:"value"`
		IsDisabled bool      `db:"is_disabled" json:"is_disabled"`
		IsReadOnly bool      `db:"is_readonly" json:"is_readonly"`
		FieldType  string    `db:"field_type"  json:"field_type"`
		DataType   string    `db:"data_type"   json:"data_type"`
		CreatedAt  time.Time `db:"created_at"  json:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"  json:"updated_at"`
	}
	MinimalModel struct {
		Key   string `db:"key"   json:"key"`
		Value string `db:"value" json:"value"`
	}
	CostingValues struct {
		OrderMinimum float64
		FreeDelivery float64
	}
	SASv4Values struct {
		Username string
		Password string
		Token    string
	}
)

// Model ----------------------------------------------------------------------

func (m *Model) GetID() string {
	return fmt.Sprintf("%d", m.ID)
}

func (m *Model) ModelName() string {
	return "setting"
}

func (m *Model) TableName() string {
	return "settings"
}

func (m *Model) DefaultSearch() string {
	return "key"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{
		m.DefaultSearch(),
		"value",
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

// utilties -------------------------------------------------------------------

func (m *Model) MergeAndValidate(v *validator.Validator) bool {
	data := v.Data
	var key string

	v.AssignString("key", &key)
	// this ensures that core keys are not modified other than value
	if !slices.Contains(CoreKeys, m.Key) {
		m.Initialize(v.Data.Values, v.DB)
		v.AssignString("key", &m.Key)
		v.AssignString("data_type", &m.DataType)
		v.AssignString("field_type", &m.FieldType)
		v.AssignBool("is_disabled", &m.IsDisabled)
		v.AssignBool("is_readonly", &m.IsReadOnly)
	}

	if data.KeyExists("value") {
		m.Value = data.Get("value")
		switch m.DataType {
		case "number":
			floatVal, err := strconv.ParseFloat(m.Value, 64)
			if err != nil {
				v.Check(false, "value", v.T.ValidateRequiredFloat())
				return false
			}
			m.Value = fmt.Sprintf("%.2f", floatVal)
		case "boolean":
			boolVal, err := strconv.ParseBool(m.Value)
			if err != nil {
				v.Check(false, "value", v.T.ValidateBool())
				return false
			}
			m.Value = fmt.Sprintf("%t", boolVal)
		default:
			v.Check(m.Value != "", "value", v.T.ValidateRequired())
		}
	}
	v.ValidateModelSchema(m, v.Schema)
	return v.Valid()
}
