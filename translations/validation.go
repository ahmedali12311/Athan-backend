package translations

import (
	"fmt"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// TODO: better translation code with maps and const IDs

func (t *Translations) InputValidation() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "InputValidation",
			Description: "input validation error.",
			Other:       "input validation error.",
		},
	})
}

// Core Validations

func (t *Translations) ValidateRequired() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateRequired",
			Description: "input value must not be empty.",
			Other:       "required field.",
		},
	})
}

func (t *Translations) ValidateRequiredArray() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateRequiredArray",
			Description: "required array input.",
			Other:       "required array.",
		},
	})
}

func (t *Translations) ValidateDate() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateDate",
			Description: "validate input as a date format.",
			Other:       "must not be a valid date format.",
		},
	})
}

func (t *Translations) ValidateBool() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateBool",
			Description: "validate input as a boolean.",
			Other:       "must be a valid boolean.",
		},
	})
}

func (t *Translations) ValidateInt() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateInt",
			Description: "validate input as an integer.",
			Other:       "must be a valid integer.",
		},
	})
}

func (t *Translations) ValidateRequiredFloat() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateRequiredFloat",
			Description: "validate input as a float.",
			Other:       "must be a valid float.",
		},
	})
}

func (t *Translations) ValidateUUID() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateUUID",
			Description: "validate input as a uuid.",
			Other:       "must be valid uuid.",
		},
	})
}

func (t *Translations) ValidateID() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateID",
			Description: "validate input as a id.",
			Other:       "must be valid id.",
		},
	})
}

func (t *Translations) ValidateExistsInDB() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateExistsInDB",
			Description: "validate input exists in the database.",
			Other:       "does not exist in the records.",
		},
	})
}

func (t *Translations) ValidateNotExistsInDB() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateNotExistsInDB",
			Description: "validate input not existing in the database.",
			Other:       "already existing record.",
		},
	})
}

func (t *Translations) ValidateMustBeInList(arg *[]string) string {
	listAsString := strings.Join(*arg, ",")
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateMustBeInList",
			Description: "the element must belong to list items",
			Other:       "must be one of values: {{ .List }}.",
		},
		TemplateData: map[string]any{
			"List": listAsString,
		},
	})
}

func (t *Translations) ValidateNotEmptyRoles() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateNotEmptyRoles",
			Description: "user must have at least one role",
			Other:       "user must have at least one role",
		},
	})
}

func (t *Translations) ValidateMustHaveRole(roleName string) string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateMustHaveRole",
			Description: "user must have the selected role",
			Other:       "user must have role: {{ .Role }}.",
		},
		TemplateData: map[string]any{
			"Role": roleName,
		},
	})
}

func (t *Translations) ValidateMustBeGteZero() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateMustBeGteZero",
			Description: "value must be greater than or equal 0",
			Other:       "value must be greater than or equal 0",
		},
	})
}

func (t *Translations) ValidateMustBeGtZero() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateMustBeGtZero",
			Description: "value must be greater than 0",
			Other:       "value must be greater than 0",
		},
	})
}

func (t *Translations) ValidateMustBeLteValue(value int) string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateMustBeLteValue",
			Description: "value must be less than or equal value",
			Other:       "value must be less than or equal {{ .Value }}",
		},
		TemplateData: map[string]int{
			"Value": value,
		},
	})
}

func (t *Translations) ValidateMustBeGteFloatValue(value float64) string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateMustBeGteFloatValue",
			Description: "value must be greater than or equal value",
			Other:       "value must be greater than or equal {{ .Value }}",
		},
		TemplateData: map[string]string{
			"Value": fmt.Sprintf("%.2f", value),
		},
	})
}

func (t *Translations) ValidateEmail() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateEmail",
			Description: "input must be a valid mail address",
			Other:       "must be a valid email address.",
		},
	})
}

func (t *Translations) ValidatePasswordConfirmationNoMatch() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidatePasswordConfirmationNoMatch",
			Description: "input password must match password_confirmation",
			Other:       "password and Confirmation don't match.",
		},
	})
}

// Specific Validation

func (t *Translations) ValidateCategoryInput() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateCategoryInput",
			Description: "category not found, or does not belong to super parent.",
			Other:       "category doesn't exist or it's sub of a different super parent.",
		},
	})
}

func (t *Translations) ValidateCategoryParent() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ValidateCategoryParent",
			Description: "not allowed for category to be parent to itself.",
			Other:       "category can't be parent to itself.",
		},
	})
}

func (t *Translations) UnDestroyableCategory() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "UnDestroyableCategory",
			Description: "not allowed to delete core categories or ones with children.",
			Other:       "not allowed to delete core categories or ones with children.",
		},
	})
}

func (t *Translations) UnsupportedLocation(name string) string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "UnsupportedLocation",
			Description: "unsupported location.",
			Other:       "location point of {{.Model}} not supported.",
		},
		TemplateData: map[string]any{
			"Model": t.ModelName(name),
		},
	})
}

func (t *Translations) NotPermitted(ctxScopes, allowed []string) string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "NotPermitted",
			Description: "user does not have sufficient scopes to modify the resource",
			Other:       "user must have any of the following scopes: [ {{ .Allowed }} ] provided scopes: [ {{ .CtxScopes }} ]",
		},
		TemplateData: map[string]any{
			"Allowed":   strings.Join(allowed, ","),
			"CtxScopes": strings.Join(ctxScopes, ","),
		},
	})
}

func (t *Translations) WalletTransactionAlreadyConfirmed() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "WalletTransactionAlreadyConfirmed",
			Description: "wallet transaction already confirmed",
			Other:       "wallet transaction already confirmed",
		},
	})
}
