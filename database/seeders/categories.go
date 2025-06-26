package seeders

import (
	"app/models"

	"bitbucket.org/sadeemTechnology/backend-model-category"
	"bitbucket.org/sadeemTechnology/backend-model-setting"
	"bitbucket.org/sadeemTechnology/backend-pgtypes"
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
		data = append(data, categories[i].Model)

		if len(categories[i].Children) > 0 {
			for c := range categories[i].Children {
				categories[i].Children[c].Sort = c
				categories[i].Children[c].Depth = categories[i].Depth + 1
				categories[i].Children[c].ParentID = &categories[i].ID
				categories[i].Children[c].SuperParentID = &categories[i].ID

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
				ID: uuid.MustParse(models.CategorySettingSocialID),
				Name: pgtypes.JSONB{
					"ar": "وسائل تواصل",
					"en": "Social Links",
				},
			},
			{
				ID: uuid.MustParse(models.CategorySettingGeneralID),
				Name: pgtypes.JSONB{
					"ar": "عامة",
					"en": "General",
				},
			},
			{
				ID: uuid.MustParse(models.CategorySettingPricingID),
				Name: pgtypes.JSONB{
					"ar": "تسعير",
					"en": "Pricing",
				},
			},
		},
	},
}
