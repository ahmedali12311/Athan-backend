package models

import (
	model "bitbucket.org/sadeemTechnology/backend-model"
	category "bitbucket.org/sadeemTechnology/backend-model-category"
	setting "bitbucket.org/sadeemTechnology/backend-model-setting"

	"app/models/city"
	"app/models/daily_prayer_times"
	"app/models/fcm_notification"
	"app/models/permission"
	"app/models/role"
	"app/models/special_topics"
	"app/models/token"
	"app/models/user"
	"app/models/user_notification"
	"app/models/wallet"

	"app/models/hadiths"

	"app/models/adhkars"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Models struct {
	Adhkars          *adhkars.Queries
	SpecialTopics    *special_topics.Queries
	Hadiths          *hadiths.Queries
	DailyPrayerTimes *daily_prayer_times.Queries
	DB               *sqlx.DB
	QB               *squirrel.StatementBuilderType

	Category         *category.Queries
	FcmNotification  *fcm_notification.Queries
	Permission       *permission.Queries
	Role             *role.Queries
	Setting          *setting.Queries
	Token            *token.Queries
	User             *user.Queries
	UserNotification *user_notification.Queries
	Wallet           *wallet.Queries
	City             *city.Queries
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
		Adhkars:          adhkars.New(d),
		SpecialTopics:    special_topics.New(d),
		Hadiths:          hadiths.New(d),
		DailyPrayerTimes: daily_prayer_times.New(d),
		DB:               db,
		QB:               &qb,

		Category:         category.New(d),
		FcmNotification:  fcm_notification.New(d),
		Permission:       permission.New(d),
		Role:             role.New(d),
		Setting:          setting.New(d),
		Token:            token.New(d),
		User:             user.New(d),
		UserNotification: user_notification.New(d),
		Wallet:           wallet.New(d),
		City:             city.New(d),
	}
}
