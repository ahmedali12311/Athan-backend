package seeders

import (
	"log"

	"bitbucket.org/sadeemTechnology/backend-model-setting"
	pgtypes "bitbucket.org/sadeemTechnology/backend-pgtypes"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func Settings(db *sqlx.DB, qb *squirrel.StatementBuilderType) {
	cols := []string{
		"id",
		"name",
		"key",
		"value",
		"is_disabled",
		"is_readonly",
		"field_type",
		"data_type",
		"category_id",
	}
	for i := range settings {
		values := []any{
			settings[i].ID,
			settings[i].Name,
			settings[i].Key,
			settings[i].Value,
			settings[i].IsDisabled,
			settings[i].IsReadOnly,
			settings[i].FieldType,
			settings[i].DataType,
			settings[i].CategoryID,
		}
		genericSeeder(db, qb, "settings", cols, values)
	}
	if _, err := db.Exec(
		`SELECT setval('settings_id_seq', (SELECT MAX(id) FROM settings));`,
	); err != nil {
		log.Panicf(
			"error executing sql sequence update settings: %s",
			err.Error(),
		)
	}
	RunningSeedTable.Append(len(settings), "settings")
}

var settings = []setting.Model{
	{
		ID: 1,
		Name: pgtypes.JSONB{
			"ar": "اسم التطبيق",
			"en": "App Name",
		},
		Key:        setting.KeyAppName,
		Value:      "اسم التطبيق",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 2,
		Name: pgtypes.JSONB{
			"ar": "عن التطبيق",
			"en": "ِAbout",
		},
		Key:        setting.KeyAbout,
		Value:      "عن التطبيق",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "textarea",
		DataType:   "string",
	},
	{
		ID: 3,
		Name: pgtypes.JSONB{
			"ar": "قوانين التطبيق",
			"en": "rules",
		},
		Key:        setting.KeyRules,
		Value:      "قوانين التطبيق",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "textarea",
		DataType:   "string",
	},
	{
		ID: 4,
		Name: pgtypes.JSONB{
			"ar": "هاتف",
			"en": "Phone",
		},
		Key:        setting.KeyAppPhone,
		Value:      "+218910001122",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 5,
		Name: pgtypes.JSONB{
			"ar": "رابط واتساب",
			"en": "Whatsapp",
		},
		Key:        setting.KeyAppWhatsappUrl,
		Value:      "whatsapp/910001122",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 6,
		Name: pgtypes.JSONB{
			"ar": "رابط فيسبوك",
			"en": "Facebook",
		},
		Key:        setting.KeyAppFacebookUrl,
		Value:      "facebook.com/app",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 7,
		Name: pgtypes.JSONB{
			"ar": "رابط تيليغرام",
			"en": "Telegram",
		},
		Key:        setting.KeyAppTelegramUrl,
		Value:      "telegram.com/910001122",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 8,
		Name: pgtypes.JSONB{
			"ar": "رابط انستاغرام",
			"en": "Instagram",
		},
		Key:        setting.KeyAppInstagramUrl,
		Value:      "instagram.com/app",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 9,
		Name: pgtypes.JSONB{
			"ar": "رابط الموقع",
			"en": "Website",
		},
		Key:        setting.KeyAppWebsiteUrl,
		Value:      "wwww.app.ly",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 10,
		Name: pgtypes.JSONB{
			"ar": "البريد الالكتروني",
			"en": "Email",
		},
		Key:        setting.KeyAppEmailUrl,
		Value:      "info@app.ly",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 11,
		Name: pgtypes.JSONB{
			"ar": "رابط اكس",
			"en": "X",
		},
		Key:        "app_x_url", // TODO:
		Value:      "x.com/app",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 12,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        setting.KeyAppLogo,
		Value:      "path/image.jpg",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "file",
		DataType:   "string",
	},
	{
		ID: 13,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        setting.KeyAppColorPrimary,
		Value:      "#6C04FC",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "color",
		DataType:   "string",
	},
	{
		ID: 14,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        setting.KeyAppColorSecondary,
		Value:      "#4E6E5D",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "color",
		DataType:   "string",
	},
	{
		ID: 15,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        setting.KeyAppColorOnPrimary,
		Value:      "#C2E7DA",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "color",
		DataType:   "string",
	},
	{
		ID: 16,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        setting.KeyAppColorOnSecondary,
		Value:      "#F1FFE7",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "color",
		DataType:   "string",
	},
	{
		ID: 17,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        setting.KeyAppAppStoreUrl,
		Value:      "https://apps.apple.com/us/app",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 18,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        setting.KeyAppGooglePlayUrl,
		Value:      "https://play.google.com/store/apps",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 19,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        setting.KeyAppPrivacyPolicy,
		Value:      "privacy policy",
		IsDisabled: false,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{ // TLync settings -------------------------------------------------------
		ID: 20,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        "tlync_endpoint",
		Value:      "https://c7drkx2ege.execute-api.eu-west-2.amazonaws.com",
		IsDisabled: true,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 21,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        "tlync_token",
		Value:      "qS00000000000000000000000000000000000000",
		IsDisabled: true,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 22,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        "tlync_store_id",
		Value:      "wL00000000000000000000000000000000000000000000000000000000000000",
		IsDisabled: true,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 23,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        "tlync_front_url",
		Value:      "https://sadeem-tech.com",
		IsDisabled: true,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{ // OTP settings ---------------------------------------------------------
		ID: 24,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        "sadeem_otp_url",
		Value:      "https://otp.sadeem-factory.com",
		IsDisabled: true,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 25,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        "sadeem_otp_jwt",
		Value:      "qS00000000000000000000000000000000000000",
		IsDisabled: true,
		IsReadOnly: false,
		FieldType:  "textarea",
		DataType:   "string",
	},
	{
		ID: 26,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        "sadeem_otp_env",
		Value:      "development",
		IsDisabled: true,
		IsReadOnly: false,
		FieldType:  "textarea",
		DataType:   "string",
	},
	{ // tyrian-ant settings --------------------------------------------------
		ID: 27,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        "payment_gateway_endpoint",
		Value:      "https://tyrian-ant.sadeem-lab.com/api/v1",
		IsDisabled: true,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
	{
		ID: 28,
		Name: pgtypes.JSONB{
			"ar": "الشعار",
			"en": "Logo",
		},
		Key:        "payment_gateway_api_key",
		Value:      "qS00000000000000000000000000000000000000",
		IsDisabled: true,
		IsReadOnly: false,
		FieldType:  "text",
		DataType:   "string",
	},
}
