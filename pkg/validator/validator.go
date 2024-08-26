package validator

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"app/config"
	"app/pkg/generics"
	"app/translations"

	"github.com/ahmedalkabir/finder"
	"github.com/labstack/echo/v4"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	js "github.com/santhosh-tekuri/jsonschema/v5"
)

// var (
//	TimeFormats = []string{
//		time.RFC3339,
//		"2006-01-02",
//		"2006-01-02T15",
//		"2006-01-02T15:04",
//		"2006-01-02 15",
//		"2006-01-02 15:04",
//	}
//
// LibyanPhoneRX = regexp.MustCompile(`^218.?([0-9]){9}$`)
// TimeRX = regexp.MustCompile("^([01]?[0-9]|2[0-3]):[0-5][0-9]$")
// InSquareBracketsRX = regexp.MustCompile(".*\\[(.*?)].*")
// ArabicRX = regexp.MustCompile("[\\u0600-\\u06FF]")
// )

var ErrNotSupportedLocation = errors.New("not supported location")

type Errors map[string][]string

// Validator type Define a new Validator type which contains a map of
// validation errors.
type Validator struct {
	T        *translations.Translations
	DB       finder.Connection
	Schema   *js.Schema
	Scopes   []string
	Data     *Data
	Error    *js.ValidationError
	Logger   *zerolog.Logger
	newFile  string
	newImg   string
	newThumb string
	oldFile  *string
	oldImg   *string
	oldThumb *string
}

// NewValidator is a helper which creates a new Validator instance with an
// empty errors map.
func NewValidator(
	ctx echo.Context,
	db finder.Connection,
	logger *zerolog.Logger,
	schemas map[string]*js.Schema,
	schemaName string,
) (*Validator, error) {
	t, ok := ctx.Get("t").(*translations.Translations)
	if !ok {
		return nil, errors.New("no \"t\" key found in validator context")
	}
	scopes, ok := ctx.Get("scopes").([]string)
	if !ok {
		return nil, errors.New("no \"scopes\" key found in validator context")
	}
	var sch *js.Schema
	if schemas != nil {
		loc := config.GetSchemaURL(schemaName)
		foundSchema, found := schemas[loc]
		if !found {
			return nil, errors.New("selected schema name not found in map")
		}
		sch = foundSchema
	}
	v := &Validator{
		T:      t,
		DB:     db,
		Schema: sch,
		Scopes: scopes,
		Logger: logger,
		Error: &js.ValidationError{
			KeywordLocation:         "",
			AbsoluteKeywordLocation: "",
			InstanceLocation:        "",
			Message:                 "",
			Causes:                  []*js.ValidationError{},
		},
	}
	if err := v.Parse(ctx.Request()); err != nil {
		return nil, err
	}
	return v, nil
}

func AssignENUM[T ~string](
	v *Validator,
	key string,
	property *T,
	allowedScopes ...string,
) {
	if v.Data.KeyExists(key) {
		v.Permit(key, allowedScopes)
		if val := v.Data.Get(key); val != "" {
			if property != nil {
				*property = T(val)
			}
		}
	}
}

func AssignNullableENUM[T ~string](
	v *Validator,
	key string,
	property *T,
	allowedScopes ...string,
) *T {
	if v.Data.KeyExists(key) {
		v.Permit(key, allowedScopes)
		if val := v.Data.Get(key); val != "" {
			enum := T(val)
			property = &enum
		}
	}
	return property
}

// Permit checks the scopes array of request context
//
// it always allows admin to modify
//
// if now allowed array is provided it assumes anyone is allowed to modify
//
// otherwise it checks each allowed against context scopes
// and errors if no match is found
func (v *Validator) Permit(key string, allowed []string) {
	if allowed == nil {
		return
	}
	if len(allowed) == 0 {
		return
	}
	for i := range v.Scopes {
		if v.Scopes[i] == "admin" {
			return
		}
		if slices.Contains(allowed, v.Scopes[i]) {
			return
		}
	}
	v.Check(false, key, v.T.NotPermitted(v.Scopes, allowed))
}

// ValidateModelSchema schema.Validate accepts only interface{} type
// so marshaling Model into []byte and then to interface{} is required.
func (v *Validator) ValidateModelSchema(model finder.Model, schema *js.Schema) {
	data, err := json.Marshal(model)
	if err != nil {
		v.AddModelSchemaError(model, errors.New("couldn't marshal input"))
		return
	}
	var body any
	if err := json.Unmarshal(data, &body); err != nil {
		v.AddModelSchemaError(model, errors.New("couldn't unmarshal input"))
		return
	}
	if schema == nil {
		v.AddModelSchemaError(model, errors.New("null validator schema"))
		return
	}
	if err := schema.Validate(body); err != nil {
		switch err := err.(type) { //nolint:errorlint // not of comparable err
		case *js.ValidationError:
			v.AddCause(err)
			return
		default:
			v.AddModelSchemaError(model, err)
			return
		}
	}
}

// ValidateInputModelSchema schema.Validate accepts only interface{} type
// so marshaling Model into []byte and then to interface{} is required.
type IInputModel interface {
	TableName() string
}

func (v *Validator) ValidateInputModelSchema(model IInputModel, schema *js.Schema) {
	data, err := json.Marshal(model)
	if err != nil {
		v.AddInputModelSchemaError(model, errors.New("couldn't marshal input"))
		return
	}
	var body any
	if err := json.Unmarshal(data, &body); err != nil {
		v.AddInputModelSchemaError(model, errors.New("couldn't unmarshal input"))
		return
	}
	if schema == nil {
		v.AddInputModelSchemaError(model, errors.New("null validator schema"))
		return
	}
	if err := schema.Validate(body); err != nil {
		switch err := err.(type) { //nolint:errorlint // not of comparable err
		case *js.ValidationError:
			v.AddCause(err)
			return
		default:
			v.AddInputModelSchemaError(model, err)
			return
		}
	}
}

func (v *Validator) ValidatePropertySchema(key string) {
	if !v.Data.KeyExists(key) {
		v.Check(false, key, v.T.ValidateRequired())
		return
	}
	v.validateSchema(v.Data.GetBytes(key), key)
}

func (v *Validator) ValidateInterfaceSchema(i any) {
	data, err := json.Marshal(i)
	if err != nil {
		v.Check(false, "input", "couldn't unmarshal input")
		return
	}
	v.validateSchema(data, "input")
}

func (v *Validator) validateSchema(data []byte, key string) {
	var body any
	if err := json.Unmarshal(data, &body); err != nil {
		v.Check(false, key, "couldn't unmarshal input")
		return
	}
	if err := v.Schema.Validate(body); err != nil {
		switch err := err.(type) { //nolint:errorlint // not of comparable err
		case *js.ValidationError:
			v.AddCause(err)
		default:
			v.Check(false, key, err.Error())
		}
	}
}

func (v *Validator) AddModelSchemaError(model finder.Model, err error) {
	cause := &js.ValidationError{
		KeywordLocation:         "model/bind",
		AbsoluteKeywordLocation: "schema/validation",
		InstanceLocation:        model.TableName() + "/model/bind",
		Message:                 err.Error(),
		Causes:                  []*js.ValidationError{},
	}
	v.Error.Causes = append(v.Error.Causes, cause)
}

func (v *Validator) AddInputModelSchemaError(model IInputModel, err error) {
	cause := &js.ValidationError{
		KeywordLocation:         "model/bind",
		AbsoluteKeywordLocation: "schema/validation",
		InstanceLocation:        model.TableName() + "/model/bind",
		Message:                 err.Error(),
		Causes:                  []*js.ValidationError{},
	}
	v.Error.Causes = append(v.Error.Causes, cause)
}

func (v *Validator) Valid() bool {
	return v.Error.Message == "" && len(v.Error.Causes) == 0
}

func (v *Validator) AddCause(cause *js.ValidationError) {
	v.Error.Causes = append(v.Error.Causes, cause)
}

// Check the boolean condition, if !ok an error will be added to the causes
func (v *Validator) Check(ok bool, instanceLocation, message string) {
	if !ok {
		cause := &js.ValidationError{
			InstanceLocation: instanceLocation,
			Message:          message,
			Causes:           []*js.ValidationError{},
		}
		v.Error.Causes = append(v.Error.Causes, cause)
	}
}

func (v *Validator) GetErrorMap() Errors {
	errMap := make(Errors)
	errMap = v.loopCauses(errMap, v.Error.Causes)
	return errMap
}

func (v *Validator) loopCauses(
	errMap Errors,
	causes []*js.ValidationError,
) Errors {
	if len(causes) > 0 {
		for _, cause := range causes {
			key := cause.InstanceLocation
			key = strings.Replace(key, "/", "", 1)
			key = strings.ReplaceAll(key, "/", ".")
			message := cause.Message
			if len(cause.Causes) != 0 {
				errMap = v.loopCauses(errMap, cause.Causes)
			} else {
				if _, exists := errMap[key]; !exists {
					errMap[key] = []string{message}
				} else {
					errMap[key] = append(errMap[key], message)
				}
			}
		}
	}
	return errMap
}

// AssignString to allow only admin to modify attribute:
//
//	v.AssignString("name", &m.Name, "admin")
//
// to allow anyone to modify:
//
//	v.AssignString("name", &m.Name)
//
// to pass specific scopes, admin will still be permitted:
//
//	v.AssignString("name", &m.Name, "vendor", "driver")
func (v *Validator) AssignString(
	key string,
	property *string,
	allowedScopes ...string,
) *string {
	if v.Data.KeyExists(key) {
		v.Permit(key, allowedScopes)
		if val := v.Data.Values.Get(key); val != "" {
			if property == nil {
				property = new(string)
			}
			*property = val
		}
	}
	return property
}

func (v *Validator) AssignBool(
	key string,
	property *bool,
	allowedScopes ...string,
) {
	if v.Data.KeyExists(key) {
		v.Permit(key, allowedScopes)
		if property == nil {
			property = generics.Ptr(false)
		}
		*property = v.Data.GetBool(key)
	}
}

func (v *Validator) AssignInt(
	key string,
	property *int,
	allowedScopes ...string,
) {
	if v.Data.KeyExists(key) {
		v.Permit(key, allowedScopes)
		if property == nil {
			property = generics.Ptr(0)
		}
		if value, err := strconv.ParseInt(v.Data.Get(key), 10, 0); err != nil {
			v.Check(false, key, v.T.ValidateInt())
		} else {
			*property = int(value)
		}
	}
}

func (v *Validator) AssignFloat(
	key string,
	property *float64,
	allowedScopes ...string,
) {
	if v.Data.KeyExists(key) {
		v.Permit(key, allowedScopes)
		if property == nil {
			property = generics.Ptr(0.0)
		}
		if value, err := strconv.ParseFloat(v.Data.Get(key), 64); err != nil {
			v.Check(false, key, v.T.ValidateRequiredFloat())
		} else {
			*property = value
		}
	}
}

func (v *Validator) AssignDate(key string, property *string) {
	if v.Data.KeyExists(key) {
		if val := v.Data.Get(key); val != "" {
			if t, err := time.Parse(time.DateOnly, val); err != nil {
				v.Check(false, key, err.Error())
			} else {
				s := t.Format("2006-01-02")
				if property == nil {
					property = generics.Ptr("")
				}
				*property = s
			}
		}
	}
}

func (v *Validator) AssignNullableDate(key string, property *string) *string {
	if v.Data.KeyExists(key) {
		if val := v.Data.Get(key); val != "" {
			if t, err := time.Parse(time.DateOnly, val); err != nil {
				v.Check(false, key, err.Error())
			} else {
				s := t.Format("2006-01-02")
				if property == nil {
					property = generics.Ptr("")
				}
				*property = s
			}
		}
	}
	return property
}

func (v *Validator) AssignTimestamp(
	key string,
	property *time.Time,
	allowedScopes ...string,
) {
	if v.Data.KeyExists(key) {
		if val := v.Data.Get(key); val != "" {
			v.Permit(key, allowedScopes)
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				v.Check(false, key, err.Error())
				return
			}
			if property != nil {
				*property = t
			}
		}
	}
}

func (v *Validator) AssignClock(
	key string,
	property *string,
	allowedScopes ...string,
) *string {
	if v.Data.KeyExists(key) {
		if val := v.Data.Get(key); val != "" {
			v.Permit(key, allowedScopes)
			t, err := time.Parse("15:04:05", val)
			if err != nil {
				v.Check(false, key, err.Error())
				return nil
			}
			if property == nil {
				property = new(string)
			}
			*property = t.Format("15:04:05")
		}
	}
	return property
}

func (v *Validator) AssignUUID(
	key, tableName string,
	property *uuid.UUID,
	required bool,
	allowedScopes ...string,
) {
	keyUUID := v.Data.GetUUID(key)
	if keyUUID != nil {
		v.Permit(key, allowedScopes)
		if property == nil {
			property = generics.Ptr(uuid.Nil)
		}
		*property = *keyUUID
		v.Exists(property, key, tableName, required)
	}
}

func (v *Validator) UnmarshalInto(
	key string,
	property any,
	allowedScopes ...string,
) {
	if v.Data.KeyExists(key) {
		v.Permit(key, allowedScopes)
		if err := v.Data.GetAndUnmarshalJSON(key, property); err != nil {
			v.Check(false, key, err.Error())
		}
	}
}

// Exists check if an id exists in a table row.
func (v *Validator) Exists(
	id any,
	key, tableName string,
	required bool,
) {
	var exists bool
	query := fmt.Sprintf(
		`SELECT EXISTS(select 1 from %s where id=$1)`,
		tableName,
	)
	if err := v.DB.GetContext(
		context.Background(),
		&exists,
		query,
		id,
	); err != nil {
		exists = false
	}
	if required {
		v.Check(exists, key, v.T.ValidateExistsInDB())
	}
}

// IDExistsInDB checks if the field value of an int id exists in database
func (v *Validator) IDExistsInDB(
	id *int,
	key, fieldName, table string,
	required bool,
) {
	if id != nil {
		if required {
			v.Check(false, fieldName, v.T.ValidateRequired())
		}
	}
	// only allows the check if the value in the model is not equal to the input
	v.Exists(id, key, table, required)
}

// UUIDExistsInDB checks if the field value of a uuid exists in database.
func (v *Validator) UUIDExistsInDB(
	id *uuid.UUID,
	key, fieldName, table string,
	required bool,
) {
	if id == nil && required {
		v.Check(false, fieldName, v.T.ValidateRequired())
	}
	v.Exists(id, key, table, required)
}

// UserIDHasRole checks if the user id has role name associated with it
func (v *Validator) UserIDHasRole(
	fieldName string,
	userID *uuid.UUID,
	roleName string,
) {
	var exists bool

	query := `
		SELECT EXISTS(
		    SELECT 1
		    FROM users
                LEFT JOIN user_roles on users.id = user_roles."user_id"
                LEFT JOIN roles on roles.id = user_roles."role_id"
		    WHERE roles.name = $2
                AND users.id = $1
		) AS exists;
	`

	if err := v.DB.GetContext(
		context.Background(),
		&exists,
		query,
		userID,
		roleName,
	); err != nil {
		exists = false
	}
	v.Check(exists, fieldName, v.T.ValidateMustHaveRole(roleName))
}
