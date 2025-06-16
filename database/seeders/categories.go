package seeders

import (
	"bitbucket.org/sadeemTechnology/backend-model-category"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

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
	for i := range categories {
		m := categories[i]
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
	RunningSeedTable.Append(len(categories), "categories")
}

// documents = parseUUID(category.DocumentSuperParent)
// trainers  = parseUUID(category.TrainerSuperParent)
// programs  = parseUUID(category.ProgramSuperParent)

var categories = []category.Model{
	// 	{
	// 		ID:            documents,
	// 		Name:          "مستندات",
	// 		Depth:         0,
	// 		Sort:          0,
	// 		IsDisabled:    false,
	// 		IsFeatured:    false,
	// 		ParentID:      nil,
	// 		SuperParentID: nil,
	// 	},
	// 	// Children
	// 	{
	// 		ID:            parseUUID("22fb185c-0af1-492e-baa5-7ced47fe66a9"),
	// 		Name:          "جواز سفر",
	// 		Depth:         1,
	// 		Sort:          0,
	// 		IsDisabled:    false,
	// 		IsFeatured:    false,
	// 		ParentID:      &documents,
	// 		SuperParentID: &documents,
	// 	},
	// 	{
	// 		ID:            parseUUID("100e38b9-50fe-4ffe-9c35-00646bba4d0c"),
	// 		Name:          "اختبار",
	// 		Depth:         1,
	// 		Sort:          1,
	// 		IsDisabled:    false,
	// 		IsFeatured:    false,
	// 		ParentID:      &documents,
	// 		SuperParentID: &documents,
	// 	},
	// 	{
	// 		ID:            parseUUID("895d5537-cdc5-43b9-9fda-77e01ef4c5f0"),
	// 		Name:          "درس",
	// 		Depth:         1,
	// 		Sort:          2,
	// 		IsDisabled:    false,
	// 		IsFeatured:    false,
	// 		ParentID:      &documents,
	// 		SuperParentID: &documents,
	// 	},

	// 	{
	// 		ID:            trainers,
	// 		Name:          "مدربين",
	// 		Depth:         0,
	// 		Sort:          1,
	// 		IsDisabled:    false,
	// 		IsFeatured:    false,
	// 		ParentID:      nil,
	// 		SuperParentID: nil,
	// 	},
	// 	// Children
	// 	{
	// 		ID:            parseUUID("3896b02e-e098-413f-96fc-3d9835b77401"),
	// 		Name:          "صيدلة",
	// 		Depth:         1,
	// 		Sort:          0,
	// 		IsDisabled:    false,
	// 		IsFeatured:    false,
	// 		ParentID:      &trainers,
	// 		SuperParentID: &trainers,
	// 	},
	// 	{
	// 		ID:            programs,
	// 		Name:          "برامج",
	// 		Depth:         0,
	// 		Sort:          2,
	// 		IsDisabled:    false,
	// 		IsFeatured:    false,
	// 		ParentID:      nil,
	// 		SuperParentID: nil,
	// 	},
	// 	// Children
	// 	{
	// 		ID:            parseUUID("ddbd3db0-a34e-4066-a2fb-68606e331b8e"),
	// 		Name:          "طب",
	// 		Depth:         1,
	// 		Sort:          0,
	// 		IsDisabled:    false,
	// 		IsFeatured:    false,
	// 		ParentID:      &programs,
	// 		SuperParentID: &programs,
	// 	},
}
