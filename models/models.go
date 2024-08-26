package models

import (
	"app/model"
	"app/models/category"
	"app/models/permission"
	"app/models/role"
	"app/models/setting"
	"app/models/token"
	"app/models/user"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/Masterminds/squirrel"
	"github.com/ahmedalkabir/finder"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Models struct {
	DB *sqlx.DB
	QB *squirrel.StatementBuilderType

	Category   *category.Queries
	Permission *permission.Queries
	Role       *role.Queries
	Setting    *setting.Queries
	Token      *token.Queries
	User       *user.Queries
}

func Setup(
	db *sqlx.DB,
	fb *firebase.App,
	fbm *messaging.Client,
	info map[string][]string,
) *Models {
	dbCache := squirrel.NewStmtCache(db)

	qb := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(dbCache)

	d := &model.Dependencies{
		DB:     db,
		QB:     &qb,
		FB:     fb,
		FBM:    fbm,
		PGInfo: info,
	}

	return &Models{
		DB: db,
		QB: &qb,

		Category: category.New(d),

		Permission: permission.New(d),

		Role:    role.New(d),
		Setting: setting.New(d),
		Token:   token.New(d),

		User: user.New(d),
	}
}

// Transaction
func (m *Models) Transaction(fn func(tx *sqlx.Tx) (finder.Model, error)) (finder.Model, error) {
	// TODO: log the operation of database transactions
	tx, err := m.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		// inner function panic
		if ex := recover(); ex != nil {
			_ = tx.Rollback()
			panic(ex)
		}
	}()

	model, err := fn(tx)

	if err != nil {
		// TODO: should we panic here? or just return error!
		// of database transaction
		_ = tx.Rollback()
	} else {
		_ = tx.Commit()
	}
	return model, err
}
