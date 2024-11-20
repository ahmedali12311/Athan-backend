package seeders

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/m-row/finder"
)

var RunningSeedTable = seededTable{
	{"table", "rows"},
}

func parseUUID(id string) uuid.UUID {
	parsed, err := uuid.Parse(id)
	if err != nil {
		log.Panic("uuid parse seeder")
	}
	return parsed
}

// func clearTable(db *sqlx.DB, name string) {
// 	query := fmt.Sprintf("delete from %s where true;", name)
// 	if _, err := db.ExecContext(context.Background(), query); err != nil {
// 		log.Panicf("error on deleting table %s: %s", name, err.Error())
// 	}
// }

type seededTable [][]string

func (st *seededTable) Append(count int, name string) {
	*st = append(*st, []string{name, fmt.Sprintf("%d", count)})
}

func PrintTable(table [][]string) {
	// get number of columns from the first table row
	columnLengths := make([]int, len(table[0]))
	for _, line := range table {
		for i, val := range line {
			if len(val) > columnLengths[i] {
				columnLengths[i] = len(val)
			}
		}
	}
	var lineLength int
	for _, c := range columnLengths {
		// +3 for 3 additional characters before and after each field: "| %s "
		lineLength += c + 3
	}
	lineLength += 1 // +1 for the last "|" in the line

	for i, line := range table {
		if i == 0 { // table header
			// lineLength-2 because of "+" as first and last character
			fmt.Printf("+%s+\n", strings.Repeat("-", lineLength-2))
		}
		for j, val := range line {
			fmt.Printf("| %*s ", columnLengths[j], val) // single row text
			if j == len(line)-1 {
				fmt.Printf("|\n")
			}
		}
		if i == 0 || i == len(table)-1 { // table header or last line
			// lineLength-2 because of "+" as first and last character
			fmt.Printf("+%s+\n", strings.Repeat("-", lineLength-2))
		}
	}
}

func genericSeeder(
	conn finder.Connection,
	qb *squirrel.StatementBuilderType,
	tableName string,
	columns []string,
	values []any,
) {
	lenCols := len(columns)
	lenVals := len(values)
	if lenCols != lenVals {
		log.Panicf(
			"error building sql seeding %s: %d columns but %d values provided",
			tableName,
			lenCols,
			lenVals,
		)
	}
	query, values, err := qb.
		Insert(tableName).
		Suffix(`ON CONFLICT DO NOTHING`).
		Columns(columns...).
		Values(values...).
		ToSql()
	if err != nil {
		log.Panicf("error building sql seeding %s: %s", tableName, err.Error())
	}
	if _, err := conn.ExecContext(
		context.Background(),
		query,
		values...,
	); err != nil {
		log.Panicf("error executing sql seeding %s: %s", tableName, err.Error())
	}
}

// func randomLibyaPoint() []float64 {
// 	// libya bounds
// 	// minLng := 9.374405
// 	// maxLng := 25.48516
// 	// minLat := 19.45087
// 	// maxLat := 33.57972
//
// 	// Benghazi bounds
// 	minLng := 19.986034
// 	maxLng := 20.317036
// 	minLat := 32.012181
// 	maxLat := 32.199584
//
// 	lng := getRandomFloat(minLng, maxLng)
// 	lat := getRandomFloat(minLat, maxLat)
//
// 	return []float64{lng, lat}
// }

// func getRandomFloat(min, max float64) float64 {
// 	r := min + rand.Float64()*(max-min) //nolint:gosec // doesn't matter
// 	return r
// }
