package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Documents ------------------------------------------------------------------

func (m *Queries) RemoveDocuments(
	userID *uuid.UUID,
	tx *sqlx.Tx,
) error {
	if _, err := tx.ExecContext(
		context.Background(),
		`
            DELETE FROM user_documents 
            WHERE user_id = $1
        `,
		userID,
	); err != nil {
		return err
	}
	return nil
}

func (m *Queries) AssignDocuments(
	userID *uuid.UUID,
	docIDS *[]string,
	tx *sqlx.Tx,
) error {
	inserts := m.QB.
		Insert("user_documents").
		Columns(
			"user_id",
			"document_id",
		).
		Suffix(`ON CONFLICT DO NOTHING`)

	if docIDS != nil {
		if len(*docIDS) == 0 {
			return nil
		}
		for _, v := range *docIDS {
			inserts = inserts.Values(userID, v)
		}
	}
	query, args, err := inserts.ToSql()
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(
		context.Background(),
		query,
		args...,
	); err != nil {
		return err
	}
	return nil
}
