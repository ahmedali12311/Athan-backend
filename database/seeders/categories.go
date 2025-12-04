package seeders

import (
	"app/models/consts"

	category "bitbucket.org/sadeemTechnology/backend-model-category"
	setting "bitbucket.org/sadeemTechnology/backend-model-setting"
	pgtypes "bitbucket.org/sadeemTechnology/backend-pgtypes"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type categorySeed struct {
	category.Model
	Children []category.Model
}

func Categories(db *sqlx.DB, qb *squirrel.StatementBuilderType) {
	cols := []string{
		"id",
		"name",
		"depth",
		"sort",
		"is_disabled",
		"is_featured",
		"parent_id",
		"super_parent_id",
	}
	data := []category.Model{}

	for i := range categories {
		categories[i].Model.SuperParentID = &categories[i].Model.ID
		data = append(data, categories[i].Model)

		if len(categories[i].Children) > 0 {
			for c := range categories[i].Children {
				categories[i].Children[c].Sort = c
				categories[i].Children[c].Depth = categories[i].Model.Depth + 1
				categories[i].Children[c].ParentID = &categories[i].Model.ID
				categories[i].Children[c].SuperParentID = &categories[i].Model.ID

				data = append(data, categories[i].Children[c])
			}
		}
	}

	for i := range data {
		m := data[i]
		values := []any{
			m.ID,
			m.Name,
			m.Depth,
			m.Sort,
			m.IsDisabled,
			m.IsFeatured,
			m.ParentID,
			m.SuperParentID,
		}
		genericSeeder(db, qb, "categories", cols, values)
	}
	RunningSeedTable.Append(len(data), "categories")
}

var categories = []categorySeed{
	{
		Model: category.Model{
			ID: uuid.MustParse(setting.SuperParentCategoryID),
			Name: pgtypes.JSONB{
				"ar": "الإعدادات",
				"en": "Settings",
			},
			Depth: 0,
			Sort:  0,
		},
		Children: []category.Model{
			{
				ID: uuid.MustParse(consts.CategorySettingSocialID),
				Name: pgtypes.JSONB{
					"ar": "وسائل تواصل",
					"en": "Social Links",
				},
			},
			{
				ID: uuid.MustParse(consts.CategorySettingGeneralID),
				Name: pgtypes.JSONB{
					"ar": "عامة",
					"en": "General",
				},
			},
			{
				ID: uuid.MustParse(consts.CategorySettingPricingID),
				Name: pgtypes.JSONB{
					"ar": "تسعير",
					"en": "Pricing",
				},
			},
			{
				ID: uuid.MustParse(consts.CategorySettingTlyncID),
				Name: pgtypes.JSONB{
					"ar": "تي لنك",
					"en": "TLync",
				},
				IsDisabled: true,
			},
			{
				ID: uuid.MustParse(consts.CategorySettingResalaID),
				Name: pgtypes.JSONB{
					"ar": "خدمة رسالة",
					"en": "Resala",
				},
				IsDisabled: true,
			},
			{
				ID: uuid.MustParse(consts.CategorySettingTyrianAntID),
				Name: pgtypes.JSONB{
					"ar": "خدمة بوابة الدفع",
					"en": "Payment Gateway",
				},
				IsDisabled: true,
			},
		},
	},
	{
		Model: category.Model{
			ID: uuid.MustParse(consts.CategorySpecialTopicID),
			Name: pgtypes.JSONB{
				"ar": "مواضيع خاصة",
				"en": "Special Toopics",
			},
			Depth: 0,
			Sort:  1,
		},
		Children: []category.Model{
			{
				ID: uuid.New(),
				Name: pgtypes.JSONB{
					"ar": "موضوع هام ",
					"en": "Important topic",
				},
			},
		},
	},
	{
		Model: category.Model{
			ID: uuid.MustParse("f1e2d3c4-b5a6-7c8d-9e0f-1a2b3c4d5e6f"),
			Name: pgtypes.JSONB{
				"ar": "حديث",
				"en": "Hadith",
			},
			Depth: 0,
			Sort:  2,
		},
		Children: []category.Model{
			{
				ID: uuid.New(),
				Name: pgtypes.JSONB{
					"ar": "أحاديث صحيحه",
					"en": "Sahih Hadiths",
				},
			},
		},
	},
	{
		Model: category.Model{
			ID: uuid.MustParse("1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d"),
			Name: pgtypes.JSONB{
				"ar": "أذكار",
				"en": "Adkhar",
			},
			Depth: 0,
			Sort:  3,
		},
		Children: []category.Model{
			{
				ID: uuid.New(), // Generate a new UUID or use a specific one
				Name: pgtypes.JSONB{
					"ar": "أذكار الصباح",
					"en": "Morning Adkhar",
				},
			},
		},
	},
}
