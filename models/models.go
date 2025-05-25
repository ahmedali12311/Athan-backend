package models

import (
	"app/models/category"
	"app/models/fcm_notification"
	"app/models/permission"
	"app/models/role"
	"app/models/setting"
	"app/models/token"
	"app/models/user"
	"app/models/user_notification"
	"app/models/wallet"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/m-row/model"
)

type Models struct {
	DB *sqlx.DB
	QB *squirrel.StatementBuilderType

	Category         *category.Queries
	FcmNotification  *fcm_notification.Queries
	Permission       *permission.Queries
	Role             *role.Queries
	Setting          *setting.Queries
	Token            *token.Queries
	User             *user.Queries
	UserNotification *user_notification.Queries
	Wallet           *wallet.Queries
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

		Category:         category.New(d),
		FcmNotification:  fcm_notification.New(d),
		Permission:       permission.New(d),
		Role:             role.New(d),
		Setting:          setting.New(d),
		Token:            token.New(d),
		User:             user.New(d),
		UserNotification: user_notification.New(d),
		Wallet:           wallet.New(d),
	}
}
