package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type TableInfo struct {
	TableName       string `db:"table_name"`
	ColumnName      string `db:"column_name"`
	IsNullable      string `db:"is_nullable"`
	DataType        string `db:"data_type"`
	TableSchema     string `db:"table_schema"`
	OrdinalPosition string `db:"ordinal_position"`
}

type ForeignKeyInfo struct {
	TableName         string `db:"table_name"`
	ColumnName        string `db:"column_name"`
	ForeignTableName  string `db:"foreign_table_name"`
	ForeignColumnName string `db:"foreign_column_name"`
	ConstraintName    string `db:"constraint_name"`
}

func getForeignKeys(db finder.Connection, table string) ([]ForeignKeyInfo, error) {
	var fks []ForeignKeyInfo
	query := `
		SELECT
			tc.table_name,
			kcu.column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name,
			tc.constraint_name
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY' 
			AND tc.table_name = $1
			AND tc.table_schema = 'public'
	`
	err := db.Select(&fks, query, table)
	return fks, err
}

func generateModelFiles(table string, db finder.Connection) error {
	ti := []TableInfo{}

	query := `
        SELECT
            table_name,
            column_name,
            is_nullable,
            data_type,
            table_schema,
            ordinal_position
        FROM "information_schema"."columns"
        WHERE table_name = $1
              AND table_schema = 'public'
        ORDER BY ordinal_position
        `
	if err := db.Select(&ti, query, table); err != nil {
		return fmt.Errorf("error querying table %s: %w", table, err)
	}

	if len(ti) == 0 {
		return fmt.Errorf("no columns found for table: %s", table)
	}

	// Get foreign keys
	fks, err := getForeignKeys(db, table)
	if err != nil {
		log.Printf("Warning: Could not get foreign keys: %v", err)
	}

	log.Printf("Found %d columns and %d foreign keys for table %s", len(ti), len(fks), table)

	// Generate model.go
	if err := generateModelGo(table, ti, fks); err != nil {
		return err
	}

	// Generate queries.go
	if err := generateQueriesGo(table, ti, fks); err != nil {
		return err
	}

	return nil
}
func generateModelGo(table string, ti []TableInfo, fks []ForeignKeyInfo) error {
	f := filepath.Clean(
		path.Join("./", "models", table, "model.go"),
	)

	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
		return err
	}

	if _, err := os.Stat(f); err == nil {
		if err = os.Remove(f); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o755)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write package and imports
	fmt.Fprintf(
		file,
		`package %s

import (
	"net/url"
	"time"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	model "bitbucket.org/sadeemTechnology/backend-model"
	pgtypes "bitbucket.org/sadeemTechnology/backend-pgtypes"
	validator "bitbucket.org/sadeemTechnology/backend-validator"
	"github.com/google/uuid"
)

var (
	ScopeAdmin = "admin"
	ScopeOwn   = "own"
)

type Model struct {
`,
		table,
	)

	// Write struct fields
	for _, v := range ti {
		datatype := getGoDataType(v.DataType, v.IsNullable, v.ColumnName)
		fieldName := toPascalCase(v.ColumnName)

		jsonTag := v.ColumnName
		if strings.HasSuffix(v.ColumnName, "_id") && v.ColumnName != "id" {
			jsonTag = "-"
		}

		fmt.Fprintf(
			file,
			"\t%s %s `db:\"%s\" json:\"%s\"`\n",
			fieldName,
			datatype,
			v.ColumnName,
			jsonTag,
		)
	}
	addedRelations := make(map[string]bool)
	for _, fk := range fks {
		if fk.ColumnName != "id" {
			relationName := toPascalCase(strings.TrimSuffix(fk.ColumnName, "_id"))
			packageName := strings.ToLower(fk.ForeignTableName)

			if !addedRelations[relationName] {
				fmt.Fprintf(file, "\t%s *%s.MinimalModel `db:\"%s\" json:\"%s\"`\n",
					relationName, packageName, strings.ToLower(relationName), strings.ToLower(relationName))
				addedRelations[relationName] = true
			}
		}
	}

	// Write MinimalModel struct
	fmt.Fprint(file, "}\n\n")
	fmt.Fprint(file, "type MinimalModel struct {\n")

	// Add essential fields to MinimalModel
	for _, v := range ti {
		if isEssentialField(v.ColumnName) {
			datatype := getGoDataType(v.DataType, v.IsNullable, v.ColumnName)
			fieldName := toPascalCase(v.ColumnName)

			fmt.Fprintf(
				file,
				"\t%s %s `db:\"%s\" json:\"%s\"`\n",
				fieldName,
				datatype,
				v.ColumnName,
				v.ColumnName,
			)
		}
	}
	fmt.Fprint(file, "}\n\n")

	// Write model methods
	fmt.Fprint(file, `// Model methods --------------------------------------------------------------

func (m *Model) GetID() string {
	return m.ID.String()
}

func (m *Model) Columns(pgInfo map[string][]string) *[]string {
	return finder.GetColumns(m, pgInfo)
}

func (m *Model) ModelName() string {
	return "`+table+`"
}

func (m *Model) TableName() string {
	return "`+getTableName(table)+`"
}

func (m *Model) DefaultSearch() string {
	return "name"
}

func (m *Model) SearchFields() *[]string {
	return &[]string{"name", "description"}
}

func (m *Model) Relations() *[]finder.RelationField {
	return &[]finder.RelationField{
`)

	addedJoins := make(map[string]bool)
	for _, fk := range fks {
		if fk.ColumnName != "id" {
			tableName := getTableName(fk.ForeignTableName)
			joinKey := fmt.Sprintf("%s.%s", getTableName(table), fk.ColumnName)

			if !addedJoins[joinKey] {
				fmt.Fprintf(file, `		{
			Table: "%s",
			Join: &finder.Join{
				From: "%s.%s",
				To:   "%s.id",
			},
		},
`, tableName, getTableName(table), fk.ColumnName, tableName)
				addedJoins[joinKey] = true
			}
		}
	}

	fmt.Fprint(file, `	}
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
	
`)
	addedRelationValidations := make(map[string]bool)
	for _, fk := range fks {
		if fk.ColumnName != "id" {
			relationName := toPascalCase(strings.TrimSuffix(fk.ColumnName, "_id"))
			foreignKeyField := toPascalCase(fk.ColumnName)
			tableName := getRelationTableName(fk.ForeignTableName)
			if !addedRelationValidations[relationName] {
				fmt.Fprintf(file, `	v.UnmarshalInto("%s", m.%s)
	if m.%s != nil && m.%s.ID != uuid.Nil {
		v.UUIDExistsInDB(&m.%s.ID, "%s", "id", "%s", true)
		m.%s = m.%s.ID
	} else {
		m.%s = uuid.Nil
	}

`, strings.ToLower(relationName), relationName,
					relationName, relationName,
					relationName, fk.ColumnName, tableName,
					foreignKeyField, relationName,
					foreignKeyField)
				addedRelationValidations[relationName] = true
			}
		}
	}
	addedValidations := make(map[string]bool)
	for _, col := range ti {
		if strings.HasSuffix(col.ColumnName, "_id") && col.ColumnName != "id" {
			continue
		}

		if addedValidations[col.ColumnName] {
			continue
		}
		addedValidations[col.ColumnName] = true

		fieldName := toPascalCase(col.ColumnName)
		switch col.DataType {
		case "text", "character varying":
			if col.IsNullable == "NO" && col.ColumnName != "img" && col.ColumnName != "thumb" {
				if strings.HasPrefix(getGoDataType(col.DataType, col.IsNullable, col.ColumnName), "*") {
					fmt.Fprintf(file, "\tm.%s = v.AssignString(\"%s\", m.%s, 1, 255)\n",
						fieldName, col.ColumnName, fieldName)
				} else {
					fmt.Fprintf(file, "\tv.AssignString(\"%s\", &m.%s, 1, 255)\n",
						col.ColumnName, fieldName)
				}
			} else {
				if strings.HasPrefix(getGoDataType(col.DataType, col.IsNullable, col.ColumnName), "*") {
					fmt.Fprintf(file, "\tm.%s = v.AssignString(\"%s\", m.%s, 0, 255)\n",
						fieldName, col.ColumnName, fieldName)
				} else {
					fmt.Fprintf(file, "\tv.AssignString(\"%s\", &m.%s, 0, 255)\n",
						col.ColumnName, fieldName)
				}
			}
		case "boolean":
			fmt.Fprintf(file, "\tv.AssignBool(\"%s\", &m.%s)\n",
				col.ColumnName, fieldName)
		case "uuid":
			if col.ColumnName != "id" && !strings.HasSuffix(col.ColumnName, "_id") {
				if col.IsNullable == "YES" {
					fmt.Fprintf(file, "\tif m.%s != nil {\n", fieldName)
					fmt.Fprintf(file, "\t\tv.UUIDExistsInDB(m.%s, \"%s\", \"id\", \"%s\", true)\n",
						fieldName, col.ColumnName, getTableName(strings.TrimSuffix(col.ColumnName, "_id")))
					fmt.Fprintf(file, "\t}\n")
				} else {
					fmt.Fprintf(file, "\tv.UUIDExistsInDB(&m.%s, \"%s\", \"id\", \"%s\", true)\n",
						fieldName, col.ColumnName, getTableName(strings.TrimSuffix(col.ColumnName, "_id")))
				}
			}
		case "numeric", "real", "double precision":
			if col.IsNullable == "NO" {
				fmt.Fprintf(file, "\tv.AssignFloat(\"%s\", &m.%s)\n",
					col.ColumnName, fieldName)
			} else {
				fmt.Fprintf(file, "\tv.AssignFloat(\"%s\", m.%s)\n",
					col.ColumnName, fieldName)
			}
		case "integer", "smallint", "bigint":
			if col.IsNullable == "NO" {
				fmt.Fprintf(file, "\tv.AssignInt(\"%s\", &m.%s)\n",
					col.ColumnName, fieldName)
			} else {
				fmt.Fprintf(file, "\tv.AssignInt(\"%s\", m.%s)\n",
					col.ColumnName, fieldName)
			}
		}
	}

	fmt.Fprint(file, "\treturn v.Valid()\n")
	fmt.Fprint(file, "}\n")

	filename := fmt.Sprintf("models/%s/model.go", table)
	if err := exec.Command("go", "fmt", filename).Run(); err != nil {
		return err
	}
	log.Printf("generated: %s\n", filename)
	return nil
}
func generateQueriesGo(table string, ti []TableInfo, fks []ForeignKeyInfo) error {
	f := filepath.Clean(
		path.Join("./", "models", table, "queries.go"),
	)

	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o755)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(
		file,
		`package %s

import (
	finder "bitbucket.org/sadeemTechnology/backend-finder"
	model "bitbucket.org/sadeemTechnology/backend-model"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/url"
)

var (
	selects = &[]string{
		"%s.*",
`,
		table,
		getTableName(table),
	)

	addedSelects := make(map[string]bool)
	for _, fk := range fks {
		if fk.ColumnName != "id" {
			relationName := strings.TrimSuffix(fk.ColumnName, "_id")
			alias := relationName[:1]

			if !addedSelects[relationName] {
				fmt.Fprintf(file, `		"%s.id as \"%s.id\"",
		"%s.name as \"%s.name\"",
`, alias, relationName, alias, relationName)
				addedSelects[relationName] = true
			}
		}
	}

	fmt.Fprint(file, "\t}\n\n")
	fmt.Fprint(file, "\tinserts = &[]string{\n")

	for _, v := range ti {
		if v.ColumnName != "id" && v.ColumnName != "created_at" && v.ColumnName != "updated_at" {
			fmt.Fprintf(file, "\t\t\"%s\",\n", v.ColumnName)
		}
	}

	fmt.Fprint(file, "\t}\n")
	fmt.Fprint(file, "\tbaseJoins = &[]string{\n")
	addedJoins := make(map[string]bool)
	for _, fk := range fks {
		if fk.ColumnName != "id" {
			tableName := getRelationTableName(fk.ForeignTableName)
			joinKey := fmt.Sprintf("%s ON %s.%s = %s.id", tableName, getTableName(table), fk.ColumnName, tableName)

			if !addedJoins[joinKey] {
				fmt.Fprintf(file, `		"%s ON %s.%s = %s.id",
`, tableName, getTableName(table), fk.ColumnName, tableName)
				addedJoins[joinKey] = true
			}
		}
	}
	fmt.Fprint(file, "\t}\n)\n\n")

	fmt.Fprint(file, `func buildInput(m *Model) (*[]any, error) {
	input := &[]any{
`)

	for _, v := range ti {
		if v.ColumnName != "id" && v.ColumnName != "created_at" && v.ColumnName != "updated_at" {
			fieldName := toPascalCase(v.ColumnName)

			switch v.DataType {
			case "USER-DEFINED":
				if strings.Contains(strings.ToLower(v.ColumnName), "geom") || strings.Contains(strings.ToLower(v.ColumnName), "point") {
					fmt.Fprintf(file, "\t\tsquirrel.Expr(\"ST_GeomFromGeoJSON(?::json)\", m.%s),\n", fieldName)
				} else {
					fmt.Fprintf(file, "\t\tm.%s,\n", fieldName)
				}
			default:
				fmt.Fprintf(file, "\t\tm.%s,\n", fieldName)
			}
		}
	}

	fmt.Fprint(file, `	}
	if len(*input) != len(*inserts) {
		return nil, finder.ErrInputLengthMismatch(input, inserts)
	}
	return input, nil
}

type Queries struct {
	*model.Dependencies
}

func New(d *model.Dependencies) *Queries {
	return &Queries{d}
}

type WhereScope struct {
	IsAdmin    bool
	IsPublic   bool
	UserID     *uuid.UUID
	QueryParams url.Values
}

func getJoins(ws *WhereScope) *[]string {
	return baseJoins
}

func wheres(ws *WhereScope) *[]squirrel.Sqlizer {
	w := []squirrel.Sqlizer{}
    if ws.IsAdmin {
       return &w
    }
	
	if ws.UserID != nil {
	}
	
	if !ws.IsAdmin {
	}
	
	return &w
}

func (m *Queries) GetAll(
	ctx echo.Context,
	ws *WhereScope,
) (*finder.IndexResponse[*Model], error) {
	c := &finder.ConfigIndex{
		DB:      m.DB,
		QB:      m.QB,
		PGInfo:  m.PGInfo,
		Wheres:  wheres(ws),
		Selects: selects,
		Joins:   getJoins(ws),
	}
	return finder.IndexBuilder[*Model](ctx.QueryParams(), c)
}

func (m *Queries) GetOne(shown *Model, ws *WhereScope) error {
	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Wheres:  wheres(ws),
		Selects: selects,
		Joins:   getJoins(ws),
	}
	return finder.ShowOne(shown, c)
}

func (m *Queries) CreateOne(created *Model, tx *sqlx.Tx) error {
	input, err := buildInput(created)
	if err != nil {
		return err
	}
	c := &finder.ConfigStore{
		DB:      m.DB,
		QB:      m.QB,
		Input:   input,
		Inserts: inserts,
		Selects: selects,
	}
	return finder.CreateOne(created, c)
}

func (m *Queries) UpdateOne(updated *Model, ws *WhereScope, tx *sqlx.Tx) error {
	input, err := buildInput(updated)
	if err != nil {
		return err
	}
	c := &finder.ConfigUpdate{
		DB:      m.DB,
		QB:      m.QB,
		Input:   input,
		Wheres:  wheres(ws),
		Inserts: inserts,
		Selects: selects,
		OptimisticLock: &finder.OptimisticLock{
			Name:  "updated_at",
			Value: updated.UpdatedAt,
		},
	}
	return finder.UpdateOne(updated, c)
}

func (m *Queries) DeleteOne(deleted *Model, ws *WhereScope, tx *sqlx.Tx) error {
	c := &finder.ConfigDelete{
		DB:      m.DB,
		QB:      m.QB,
		Wheres:  wheres(ws),
		Selects: selects,
	}
	return finder.DeleteOne(deleted, c)
}
`)

	filename := fmt.Sprintf("models/%s/queries.go", table)
	if err := exec.Command("go", "fmt", filename).Run(); err != nil {
		return err
	}
	log.Printf("generated: %s\n", filename)
	return nil
}
func getGoDataType(sqlType, isNullable, columnName string) string {
	var datatype string

	switch sqlType {
	case "numeric", "real", "double precision":
		datatype = "float64"
	case "timestamp with time zone", "timestamp without time zone", "date", "time":
		datatype = "time.Time"
	case "text", "character varying", "char", "name":
		datatype = "string"
	case "interval":
		datatype = "string"
	case "jsonb", "json", "USER-DEFINED":
		if sqlType == "USER-DEFINED" {
			datatype = "pgtypes.JSONB"
		} else {
			datatype = "pgtypes.JSONB"
		}
	case "bytea":
		datatype = "[]byte"
	case "integer", "smallint":
		datatype = "int"
	case "bigint":
		datatype = "int64"
	case "uuid":
		datatype = "uuid.UUID"
	case "boolean":
		datatype = "bool"
	default:
		datatype = "string"
	}

	if sqlType == "uuid" && strings.HasSuffix(columnName, "_id") {
		return datatype
	}

	if isNullable == "YES" && sqlType != "jsonb" && sqlType != "USER-DEFINED" && sqlType != "bytea" && !strings.Contains(sqlType, "timestamp") {
		datatype = "*" + datatype
	}

	return datatype
}
func isEssentialField(columnName string) bool {
	essential := []string{"id", "name", "title", "email", "created_at", "updated_at", "price", "is_disabled"}
	for _, ess := range essential {
		if columnName == ess {
			return true
		}
	}
	return false
}
func getTableName(table string) string {
	if strings.HasSuffix(table, "y") {
		return strings.TrimSuffix(table, "y") + "ies"
	}
	if strings.HasSuffix(table, "s") || strings.HasSuffix(table, "x") || strings.HasSuffix(table, "z") ||
		strings.HasSuffix(table, "ch") || strings.HasSuffix(table, "sh") {
		return table + "es"
	}
	return table + "s"
}
func getRelationTableName(table string) string {
	return table
}

func toPascalCase(str string) string {
	var res string
	isUnderscore := false

	// Handle common ID patterns
	if strings.HasPrefix(str, "id") {
		str = strings.Replace(str, "id", "ID", 1)
	}
	if strings.HasSuffix(str, "_id") {
		str = strings.Replace(str, "_id", "ID", 1)
	}
	str = strings.Replace(str, "_id_", "ID", 1)

	for i := range str {
		l := string(str[i])

		if l == "_" {
			isUnderscore = true
		} else if isUnderscore {
			isUnderscore = false
			l = strings.ToUpper(l)
		} else if i == 0 {
			l = strings.ToUpper(l)
		} else {
			isUnderscore = false
		}
		if !isUnderscore {
			res += l
		}
	}
	return res
}
