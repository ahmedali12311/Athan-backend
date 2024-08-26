package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (m *Queries) IncrementLastRef(
	id *uuid.UUID,
	tx *sqlx.Tx,
) error {
	query := `
        UPDATE users 
        SET last_ref = last_ref + 1 
        WHERE id=$1
    `
	conn := m.DB
	if tx != nil {
		conn = tx
	}
	if _, err := conn.ExecContext(context.Background(), query, id); err != nil {
		return err
	}
	return nil
}
