//nolint:lll
package translations

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	PgeDefaultErr                 = "PgeDefaultErr"
	PgeNotNullConstraintViolation = "23502"
	PgeForeignKeyViolation        = "23503"
	PgeUniqueConstraintViolation  = "23505"
	PgeCheckConstraintViolation   = "23514"
)

func (t *Translations) PGError(
	method, pgErrCode, constraint string,
) string {
	var message string
	switch pgErrCode {
	case PgeNotNullConstraintViolation:
		message = t.Localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "PgeNotNullConstraintViolation",
				Description: "database: not null constraint violation",
				Other:       "some required values are empty.",
			},
		})
	case PgeForeignKeyViolation:
		switch method {
		case "DELETE":
			message = t.Localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:          "PgeForeignKeyViolationDelete",
					Description: "database: foreign key violation",
					Other:       "unable to delete record, violates foreign key: {{ .ForeignKey }}",
				},
				TemplateData: map[string]any{
					"ForeignKey": constraint,
				},
			})
		default:
			message = t.Localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:          "PgeForeignKeyViolationModify",
					Description: "database: foreign key violation",
					Other:       "unable to modify record, violates foreign key: {{ .ForeignKey }}",
				},
				TemplateData: map[string]any{
					"ForeignKey": constraint,
				},
			})
		}
	case PgeUniqueConstraintViolation:
		message = t.Localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "PgeUniqueConstraintViolation",
				Description: "database: unique constraint violation",
				Other:       `this record contains duplicated data that conflicts with what is already in the database.`, //nolint:lll
			},
		})
	case PgeCheckConstraintViolation:
		message = t.Localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "PgeCheckConstraintViolation",
				Description: "database: check constraint violation",
				Other:       "this record contains inconsistent or out-of-range data.",
			},
		})
	default:
		message = t.Localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          PgeDefaultErr,
				Description: "database: default database error.",
				Other:       "database error.",
			},
		})
	}

	return message
}
