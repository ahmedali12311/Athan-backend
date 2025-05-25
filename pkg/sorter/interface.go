package sorter

// HasSort is an interface that defines methods required for models that support
// sorting.
// Models implementing this interface can be used with sorting logic to adjust
// their position
// in a list or table based on a "sort" field.
type HasSort interface {
	// TableName returns the name of the database table associated with model.
	// This is used to construct SQL queries dynamically.

	TableName() string

	// InterfaceSortFields returns two values:
	//   1. A pointer to an integer representing the new sort value for model.
	//      If the pointer is nil, no sorting adjustment is needed.
	//   2. A map of field names to their values, used to filter records in the
	//      database.
	//      These fields are typically used in the WHERE clause of SQL queries
	//      to ensure
	//      that sorting adjustments are applied only to relevant records.
	//
	// Example:
	//   func (m MyModel) InterfaceSortFields() (*int, map[string]any) {
	//       return &m.Sort, map[string]any{"category_id": m.Category.ID}
	//   }
	InterfaceSortFields() (*int, map[string]any)

	// GetID returns the unique identifier (ID) of the model.
	GetID() string
}
