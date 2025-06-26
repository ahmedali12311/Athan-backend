package seeders

import (
	"log"

	"app/models"

	setting "bitbucket.org/sadeemTechnology/backend-model-setting"
	pgtypes "bitbucket.org/sadeemTechnology/backend-pgtypes"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type settingSeed struct {
	CategoryID *uuid.UUID
	Settings   []setting.Model
}

func Settings(db *sqlx.DB, qb *squirrel.StatementBuilderType) {
	cols := []string{
		"sort",
		"name",
		"key",
		"value",
		"is_disabled",
		"is_readonly",
		"field_type",
		"data_type",
		"category_id",
	}

	data := []setting.Model{}

	for i := range settings {
		if len(settings[i].Settings) > 0 {
			for j := range settings[i].Settings {
				settings[i].Settings[j].Sort = j
				settings[i].Settings[j].CategoryID = settings[i].CategoryID
				data = append(data, settings[i].Settings[j])
			}
		}
	}

	for i := range data {
		values := []any{
			data[i].Sort,
			data[i].Name,
			data[i].Key,
			data[i].Value,
			data[i].IsDisabled,
			data[i].IsReadOnly,
			data[i].FieldType,
			data[i].DataType,
			data[i].CategoryID,
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
	RunningSeedTable.Append(len(data), "settings")
}

var settings = []settingSeed{
	{
		CategoryID: parsePointerUUID(models.CategorySettingGeneralID),
		Settings: []setting.Model{
			{
				Name: pgtypes.JSONB{
					"ar": "اسم التطبيق",
					"en": "App Name",
				},
				Key:       setting.KeyAppName,
				Value:     "اسم التطبيق",
				FieldType: "text",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "عن التطبيق",
					"en": "ِAbout",
				},
				Key:       setting.KeyAbout,
				Value:     "عن التطبيق",
				FieldType: "textarea",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "قوانين التطبيق",
					"en": "rules",
				},
				Key:       setting.KeyRules,
				Value:     "قوانين التطبيق",
				FieldType: "textarea",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "الشعار",
					"en": "Logo",
				},
				Key:       setting.KeyAppLogo,
				Value:     "path/image.jpg",
				FieldType: "file",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "لون التطبيق الأساسي",
					"en": "App Color Primary",
				},
				Key:       setting.KeyAppColorPrimary,
				Value:     "#6C04FC",
				FieldType: "color",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "لون التطبيق على الأساسي",
					"en": "App Color On Primary",
				},
				Key:       setting.KeyAppColorOnPrimary,
				Value:     "#C2E7DA",
				FieldType: "color",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "لون التطبيق الثانوي",
					"en": "App Color Secondary",
				},
				Key:       setting.KeyAppColorSecondary,
				Value:     "#4E6E5D",
				FieldType: "color",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "لون التطبيق على الثانوي",
					"en": "App Color On Secondary",
				},
				Key:       setting.KeyAppColorOnSecondary,
				Value:     "#F1FFE7",
				FieldType: "color",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "سياسات الخصوصية",
					"en": "Privacy Policy",
				},
				Key:       setting.KeyAppPrivacyPolicy,
				Value:     "privacy policy",
				FieldType: "text",
				DataType:  "string",
			},
		},
	},
	{
		CategoryID: parsePointerUUID(models.CategorySettingSocialID),
		Settings: []setting.Model{
			{
				Name: pgtypes.JSONB{
					"ar": "رابط فيسبوك",
					"en": "Facebook",
				},
				Key:       setting.KeyAppFacebookUrl,
				Value:     "facebook.com/app",
				FieldType: "text",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رابط تيليغرام",
					"en": "Telegram",
				},
				Key:       setting.KeyAppTelegramUrl,
				Value:     "telegram.com/910001122",
				FieldType: "text",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رابط انستاغرام",
					"en": "Instagram",
				},
				Key:       setting.KeyAppInstagramUrl,
				Value:     "instagram.com/app",
				FieldType: "text",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رابط الموقع",
					"en": "Website",
				},
				Key:       setting.KeyAppWebsiteUrl,
				Value:     "wwww.app.ly",
				FieldType: "text",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "البريد الالكتروني",
					"en": "Email",
				},
				Key:       setting.KeyAppEmailUrl,
				Value:     "info@app.ly",
				FieldType: "text",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رابط اكس",
					"en": "X",
				},
				Key:       setting.KeyAppXUrl,
				Value:     "x.com/app",
				FieldType: "text",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "هاتف",
					"en": "Phone",
				},
				Key:       setting.KeyAppPhone,
				Value:     "+218910001122",
				FieldType: "text",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رابط واتساب",
					"en": "Whatsapp",
				},
				Key:       setting.KeyAppWhatsappUrl,
				Value:     "whatsapp/910001122",
				FieldType: "text",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رابط الأب ستور",
					"en": "app_app_store_url",
				},
				Key:       setting.KeyAppAppStoreUrl,
				Value:     "https://apps.apple.com/us/app",
				FieldType: "text",
				DataType:  "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رابط غووغل بلاي",
					"en": "app_google_play_url",
				},
				Key:       setting.KeyAppGooglePlayUrl,
				Value:     "https://play.google.com/store/apps",
				FieldType: "text",
				DataType:  "string",
			},
		},
	},
	{
		CategoryID: parsePointerUUID(models.CategorySettingTlyncID),
		Settings: []setting.Model{
			{
				Name: pgtypes.JSONB{
					"ar": "رابط تي لنك",
					"en": "TLync Endpoint",
				},
				Key:        setting.KeyTlyncEndpoint,
				Value:      "https://c7drkx2ege.execute-api.eu-west-2.amazonaws.com",
				IsDisabled: true,
				FieldType:  "text",
				DataType:   "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رمز تي لنك",
					"en": "TLync Token",
				},
				Key:        setting.KeyTlyncToken,
				Value:      "qS00000000000000000000000000000000000000",
				IsDisabled: true,
				FieldType:  "text",
				DataType:   "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "معرف المتجر لتي لنك",
					"en": "TLync Store ID",
				},
				Key:        setting.KeyTlyncStoreID,
				Value:      "wL00000000000000000000000000000000000000000000000000000000000000",
				IsDisabled: true,
				FieldType:  "text",
				DataType:   "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رابط العودى من تي لنك",
					"en": "TLync Front URL",
				},
				Key:        setting.KeyTlyncFrontUrl,
				Value:      "https://sadeem-tech.com",
				IsDisabled: true,
				FieldType:  "text",
				DataType:   "string",
			},
		},
	},
	{
		CategoryID: parsePointerUUID(models.CategorySettingResalaID),
		Settings: []setting.Model{
			{
				Name: pgtypes.JSONB{
					"ar": "رابط رسالة",
					"en": "Resala URL",
				},
				Key:        setting.KeyResalaURL,
				Value:      "https://otp.sadeem-factory.com",
				IsDisabled: true,
				FieldType:  "text",
				DataType:   "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رمز رسالة",
					"en": "Resala JWT",
				},
				Key:        setting.KeyResalaJWT,
				Value:      "qS00000000000000000000000000000000000000",
				IsDisabled: true,
				FieldType:  "textarea",
				DataType:   "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "بيئة رسالة",
					"en": "Resala Environment",
				},
				Key:        setting.KeyResalaEnv,
				Value:      "development",
				IsDisabled: true,
				FieldType:  "text",
				DataType:   "string",
			},
		},
	},
	{
		CategoryID: parsePointerUUID(models.CategorySettingTyrianAntID),
		Settings: []setting.Model{
			{
				Name: pgtypes.JSONB{
					"ar": "رابط بوابة الدفع",
					"en": "Payment Gateway Endpoint",
				},
				Key:        setting.KeyPaymentGatewayEndpoint,
				Value:      "https://tyrian-ant.sadeem-lab.com/api/v1",
				IsDisabled: true,
				FieldType:  "text",
				DataType:   "string",
			},
			{
				Name: pgtypes.JSONB{
					"ar": "رمز بوابة الدفع",
					"en": "Payment Gateway Api Key",
				},
				Key:        setting.KeyPaymentGatewayAPIKey,
				Value:      "qS00000000000000000000000000000000000000",
				IsDisabled: true,
				FieldType:  "text",
				DataType:   "string",
			},
		},
	},
}
