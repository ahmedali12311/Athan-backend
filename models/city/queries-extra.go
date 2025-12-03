package city

import (
	"context"

	pgtypes "bitbucket.org/sadeemTechnology/backend-pgtypes"
)

func (m *Queries) GetClosestCity(city *Model, location *pgtypes.Point) error {
	query, args, err := m.QB.
		Select(*selects...).
		From("cities").
		OrderBy(
			`ST_Distance(
           location::geography,
           ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
       		)`,
		).
		Limit(1).
		ToSql()

	if err != nil {
		return err
	}

	args = append(args, location.Coordinates[0], location.Coordinates[1])

	return m.DB.GetContext(context.Background(), city, query, args...)

}
