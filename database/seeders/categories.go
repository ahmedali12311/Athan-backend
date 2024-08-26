package seeders

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Fixtures interface {
	File() string
}

// Categories ===================================================
type Category struct {
	ID            uuid.UUID  `toml:"id"`
	Name          string     `toml:"name"`
	Depth         int        `toml:"depth"`
	Sort          int        `toml:"sort"`
	IsDisabled    bool       `toml:"is_disabled"`
	IsFeatured    bool       `toml:"is_featured"`
	ParentID      *uuid.UUID `toml:"parent_id"`
	SuperParentID *uuid.UUID `toml:"super_parent_id"`
}

type Categories struct {
	Elements []Category `toml:"categories"`
}

func (Categories) File() string {
	return "categories.toml"
}

func (Categories) Table() string {
	return "categories"
}

func (s *Categories) Seed(db *sqlx.DB, qb *squirrel.StatementBuilderType) error {
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
	for _, v := range s.Elements {
		values := []any{
			v.ID,
			v.Name,
			v.Depth,
			v.Sort,
			v.IsDisabled,
			v.IsFeatured,
			v.ParentID,
			v.SuperParentID,
		}
		genericSeeder(db, qb, "categories", cols, values)
	}
	RunningSeedTable.Append(len(s.Elements), "categories")
	return nil
}

func CategoriesLoadFixtures() (*Categories, error) {
	var categories Categories

	err := loadFixtures(&categories)
	if err != nil {
		return nil, err
	}
	return &categories, nil
}
