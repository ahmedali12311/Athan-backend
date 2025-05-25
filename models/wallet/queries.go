//nolint:lll
package wallet

import (
	"context"

	"app/models/wallet_transaction"
	"github.com/m-row/model"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder"
)

type Queries struct {
	*model.Dependencies
}

func New(d *model.Dependencies) *Queries {
	return &Queries{d}
}

var (
	selects = &[]string{
		"wallet_transactions.*",

		"users.id as \"user.id\"",
		"users.name as \"user.name\"",

		"recharger.id as \"recharged_by.id\"",
		"recharger.name as \"recharged_by.name\"",
	}
	joins = &[]string{
		"users ON wallet_transactions.wallet_id = users.id",
		"users as recharger ON wallet_transactions.recharged_by_id = recharger.id",
	}
	inserts = &[]string{
		"id",
		"wallet_id",
		"type",
		"amount",
		"payment_method",
		"payment_reference",
		"notes",
		"recharged_by_id",
		"is_confirmed",
		"tlync_url",
		"tlync_response",
	}
)

func wheres(userID *uuid.UUID) *[]squirrel.Sqlizer {
	w := &[]squirrel.Sqlizer{}
	if userID != nil {
		*w = append(
			*w,
			squirrel.Expr("wallet_transactions.wallet_id=?", userID),
		)
	}
	return w
}

// buildInput match the returned array to defaultInsertColumns order.
func buildInput(wt *wallet_transaction.Model) (*[]any, error) {
	input := &[]any{
		wt.ID,
		wt.User.ID,
		wt.Type.String(),
		wt.Amount,
		wt.PaymentMethod,
		wt.PaymentReference,
		wt.Notes,
		wt.RechargedBy.ID,
		wt.IsConfirmed,
		wt.TLyncURL,
		wt.TLyncResponse,
	}
	if len(*input) != len(*inserts) {
		return nil, finder.ErrInputLengthMismatch(input, inserts)
	}
	return input, nil
}

func (m *Queries) GetOne(id *uuid.UUID) (*Model, error) {
	var wallet Model
	query, args, err := m.QB.
		Select("wallets.*").
		From("wallets").
		Where("wallets.id=?", id).
		ToSql()
	if err != nil {
		return nil, err
	}
	if err := m.DB.
		GetContext(context.Background(), &wallet, query, args...); err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (m *Queries) GetAllTransactions(
	ctx echo.Context,
	userID *uuid.UUID,
) (*finder.IndexResponse[*wallet_transaction.Model], error) {
	c := &finder.ConfigIndex{
		DB:      m.DB,
		QB:      m.QB,
		PGInfo:  m.PGInfo,
		Joins:   joins,
		Wheres:  wheres(userID),
		Selects: selects,
		GroupBys: &[]string{
			"wallet_transactions.id",
			"users.id",
			"recharger.id",
		},
	}
	return finder.IndexBuilder[*wallet_transaction.Model](ctx.QueryParams(), c)
}

func (m *Queries) GetTransaction(
	shown *wallet_transaction.Model,
	userID *uuid.UUID,
) error {
	c := &finder.ConfigShow{
		DB:      m.DB,
		QB:      m.QB,
		Joins:   joins,
		Wheres:  wheres(userID),
		Selects: selects,
	}
	return finder.ShowOne(shown, c)
}

func (m *Queries) UpdateOne(
	walletID *uuid.UUID,
	conn finder.Connection,
) error {
	query := `
        UPDATE 
            wallets 
        SET 
            trx_count_credit = t.trx_count_credit, 
            trx_count_debit = t.trx_count_debit, 
            trx_total_credit = COALESCE(t.trx_total_credit, 0), 
            trx_total_debit = COALESCE(t.trx_total_debit, 0), 
            credit = COALESCE(t.trx_total_credit, 0) - COALESCE(t.trx_total_debit, 0) 
        FROM (
            SELECT (
                SELECT COUNT(amount) 
                FROM wallet_transactions
                WHERE wallet_id = $1 
                  AND type = 'credit'
                  AND is_confirmed = true) AS trx_count_credit, 

              (SELECT COUNT(amount) 
                FROM wallet_transactions
                WHERE wallet_id = $1 
                  AND type = 'debit') AS trx_count_debit, 

              (SELECT SUM(amount) 
                FROM "wallet_transactions" 
                WHERE wallet_id = $1 
                  AND type = 'credit'
                  AND is_confirmed = true) AS trx_total_credit, 

              (SELECT SUM(amount) 
                FROM "wallet_transactions" 
                WHERE 
                  wallet_id = $1 
                  AND type = 'debit') AS trx_total_debit) AS t 
        WHERE 
          id = $1 RETURNING wallets.*
    `

	if _, err := conn.ExecContext(
		context.Background(),
		query,
		walletID,
	); err != nil {
		return err
	}
	return nil
}

func (m *Queries) CreateTransaction(
	created *wallet_transaction.Model,
	conn finder.Connection,
) error {
	input, err := buildInput(created)
	if err != nil {
		return err
	}
	c := &finder.ConfigStore{
		DB:      conn,
		QB:      m.QB,
		Input:   input,
		Joins:   joins,
		Selects: selects,
		Inserts: inserts,
	}
	if err := finder.CreateOne(created, c); err != nil {
		return err
	}
	return m.UpdateOne(created.User.ID, conn)
}

func (m *Queries) UpdateTransaction(
	updated *wallet_transaction.Model,
	userID *uuid.UUID,
	conn finder.Connection,
) error {
	input, err := buildInput(updated)
	if err != nil {
		return err
	}
	c := &finder.ConfigUpdate{
		DB:      conn,
		QB:      m.QB,
		Input:   input,
		Joins:   joins,
		Selects: selects,
		Inserts: inserts,
		Wheres:  wheres(userID),
		OptimisticLock: &finder.OptimisticLock{
			Name:  "updated_at",
			Value: updated.UpdatedAt,
		},
	}
	if err := finder.UpdateOne(updated, c); err != nil {
		return err
	}
	return m.UpdateOne(updated.User.ID, conn)
}

func (m *Queries) DestroyTransaction(
	deleted *wallet_transaction.Model,
	userID *uuid.UUID,
	conn finder.Connection,
) error {
	c := &finder.ConfigDelete{
		DB:      conn,
		QB:      m.QB,
		Joins:   joins,
		Selects: selects,
		Wheres:  wheres(userID),
	}
	if err := finder.DeleteOne(deleted, c); err != nil {
		return err
	}
	return m.UpdateOne(deleted.User.ID, conn)
}
