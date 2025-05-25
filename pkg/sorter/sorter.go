package sorter

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/m-row/finder"
)

type OperationType string

const (
	Create OperationType = "create"
	Update OperationType = "update"
	Delete OperationType = "delete"
)

func AdjustSort(
	operation OperationType,
	model HasSort,
	conn finder.Connection,
	qb *squirrel.StatementBuilderType,
) error {
	newSort, sortFields := model.InterfaceSortFields()

	if newSort == nil {
		return nil // No sort adjustment needed
	}

	var delta int
	var conditions squirrel.And
	var args []any

	switch operation {
	case Create:
		delta = 1
		conditions = append(conditions, squirrel.GtOrEq{"sort": *newSort})
	case Delete:
		query := fmt.Sprintf(
			"SELECT sort FROM %s WHERE id = $1",
			model.TableName(),
		)
		if err := conn.GetContext(
			context.Background(),
			&newSort,
			query,
			model.GetID(),
		); err != nil {
			return err
		}
		delta = -1
		conditions = append(conditions, squirrel.GtOrEq{"sort": *newSort})
	case Update:
		// Fetch oldSort for update operation
		var oldSort int
		query := fmt.Sprintf(
			"SELECT sort FROM %s WHERE id = $1",
			model.TableName(),
		)
		if err := conn.GetContext(
			context.Background(),
			&oldSort,
			query,
			model.GetID(),
		); err != nil {
			return err
		}

		if oldSort == *newSort {
			return nil // No change needed
		}

		if oldSort > *newSort {
			delta = 1
			conditions = append(
				conditions,
				squirrel.Lt{"sort": oldSort},
				squirrel.GtOrEq{"sort": *newSort},
			)
		} else {
			delta = -1
			conditions = append(
				conditions,
				squirrel.Gt{"sort": oldSort},
				squirrel.LtOrEq{"sort": *newSort})
		}
	default:
		return fmt.Errorf("invalid operation: %v", operation)
	}

	for key, value := range sortFields {
		conditions = append(conditions, squirrel.Eq{key: value})
	}

	query, args, err := qb.Update(model.TableName()).
		Set("sort", squirrel.Expr("sort + ?", delta)).
		Where(conditions).
		ToSql()
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(context.Background(), query, args...)
	return err
}
