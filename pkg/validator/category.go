package validator

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CategoryValidator checks if category with super parent exists in database
func (v *Validator) CategoryValidator(
	id *uuid.UUID,
	fieldName, superParentId string,
) {
	var exists bool
	query := `
        SELECT EXISTS(
            SELECT 1 
            FROM Categories 
            WHERE id=$1
                AND super_parent_id=$2
        ) AS exists
    `
	if err := v.DB.GetContext(
		context.Background(),
		&exists,
		query,
		id,
		superParentId,
	); err != nil {
		exists = false
	}
	if !exists {
		v.Check(exists, fieldName, v.T.ValidateCategoryInput())
	}
}

func (v *Validator) ValidateCategoryArray(
	fieldName, superParentID string,
	required bool,
) *[]string {
	arr, ok := v.Data.Values[fieldName]
	if required && !ok {
		v.Check(false, fieldName, v.T.ValidateRequiredArray())
	}
	if required && len(arr) == 0 {
		v.Check(
			false,
			fmt.Sprintf("%s.0", fieldName),
			v.T.ValidateUUID(),
		)
	}
	if ok && len(arr) > 0 {
		for index, id := range arr {
			if _, err := uuid.Parse(id); err != nil {
				v.Check(
					false,
					fmt.Sprintf("%s.%d", fieldName, index),
					v.T.ValidateUUID(),
				)
			} else {
				var exists bool
				query := `
                    SELECT EXISTS(
                        SELECT 1 
                        FROM categories 
                        WHERE id=$1 
                              AND super_parent_id=$2
                    )
                `
				if err := v.DB.GetContext(
					context.Background(),
					&exists,
					query,
					id,
					superParentID,
				); err != nil {
					exists = false
				}
				if required && !exists {
					v.Check(
						exists,
						fmt.Sprintf("%s.%d", fieldName, index),
						v.T.ValidateCategoryInput(),
					)
				}
			}
		}
	}
	return &arr
}

// ValidateListUUIDs unmarshalls a key to a string slice
func (v *Validator) ValidateListUUIDs(
	fieldName, tableName string,
	required bool,
	allowedScopes ...string,
) *[]string {
	arr := []string{}
	v.UnmarshalInto(fieldName, &arr, allowedScopes...)
	if required && len(arr) == 0 {
		v.Check(false, fieldName, v.T.ValidateRequiredArray())
	}
	if len(arr) > 0 {
		for index, id := range arr {
			if _, err := uuid.Parse(id); err != nil {
				v.Check(
					false,
					fmt.Sprintf("%s.%d", fieldName, index),
					v.T.ValidateUUID(),
				)
			} else {
				var exists bool
				query := fmt.Sprintf(
					`SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1) AS exists`,
					tableName,
				)
				if err := v.DB.GetContext(
					context.Background(),
					&exists,
					query,
					id,
				); err != nil {
					exists = false
				}
				if required && !exists {
					v.Check(
						exists,
						fmt.Sprintf("%s.%d", fieldName, index),
						v.T.ValidateExistsInDB(),
					)
				}
			}
		}
	}
	return &arr
}
