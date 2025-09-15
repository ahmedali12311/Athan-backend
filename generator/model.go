package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"bitbucket.org/sadeemTechnology/backend-finder"
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

func generateModelFile(table string, db finder.Connection) error {
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
              -- AND table_name NOT IN (
              --     'spatial_ref_sys',
              --     'geography_columns',
              --     'geometry_columns',
              --     'schema_migrations'
              -- )
        ORDER BY table_name, ordinal_position
        `
	if err := db.Select(&ti, query, table); err != nil {
		return err
	}

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

	fmt.Fprintf(
		file,
		"package %s\n\n"+
			"import (\n"+
			"\"fmt\"\n"+
			"\"time\"\n"+
			"\"bitbucket.org/sadeemTechnology/backend-finder\"\n"+
			"pgtypes \"bitbucket.org/sadeemTechnology/backend-pgtypes\"\n"+

			"\"github.com/google/uuid\"\n"+
			")\ntype Model struct {\n",
		table,
	)
	for _, v := range ti {
		var datatype string

		// uuid
		// boolean
		// numeric
		// timestamp with time zone
		// text
		// character varying
		// interval
		// jsonb
		// bytea
		// integer
		// real

		// name
		// smallint
		// inet
		// regproc
		// bigint
		// pg_dependencies
		// xid
		// "char"
		// pg_lsn
		// pg_node_tree
		// anyarray
		// regtype
		// double precision
		// ARRAY
		// pg_ndistinct
		// pg_mcv_list
		// oid
		switch v.DataType {
		case "numeric", "real":
			datatype = "float64"
		case "timestamp with time zone":
			datatype = "time.Time"
		case "text", "character varying", "date":
			datatype = "string"
		case "interval":
			datatype = "string"
		case "jsonb", "USER-DEFINED":
			datatype = "pgtypes.JSONB"
		case "bytea":
			datatype = "[]byte"
		case "integer", "smallint":
			datatype = "int"
		case "uuid":
			datatype = "uuid.UUID"
		case "boolean":
			datatype = "bool"
		}
		if v.IsNullable == "YES" {
			datatype = "*" + datatype
		}
		if v.TableName == table {
			fmt.Fprintf(
				file,
				"\t%s %s `db:\"%s\" json:\"%s\"`\n",
				toPascalCase(v.ColumnName),
				datatype,
				v.ColumnName,
				v.ColumnName,
			)
		}
	}
	fmt.Fprintf(
		file,
		"}\n",
	)

	filename := fmt.Sprintf("models/%s/model.go", table)
	if err := exec.Command("go", "fmt", filename).Run(); err != nil {
		return err
	}
	log.Printf("generated: %s\n", filename)
	return nil
}

func toPascalCase(str string) string {
	var res string
	isUnderscore := false

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
