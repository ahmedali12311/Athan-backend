// translations for standardized translated text, app specific messages
// can be added in separate files
//
//nolint:lll
package translations

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Translations struct {
	Localizer *i18n.Localizer
}

func (t *Translations) OutOfScopeError() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "OutOfScopeError",
			Description: "an error that is out of scope of api usage",
			Other:       "out of scope error.",
		},
	})
}

func (t *Translations) BadRequest() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "BadRequest",
			Description: "bad request",
			Other:       "Bad Request.",
		},
	})
}

func (t *Translations) NotFound() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "NotFound",
			Description: "a requested uri is not found",
			Other:       "resource not found.",
		},
	})
}

// ModelName translates the model name separately.
func (t *Translations) ModelName(name string) string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: name,
		},
	})
}

func (t *Translations) ModelNotFound(name string) string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ModelNotFound",
			Description: "a requested model record is not found",
			Other:       "{{.Model}} not found.",
		},
		TemplateData: map[string]interface{}{
			"Model": name,
		},
	})
}

func (m *Translations) ModelDisabled(name string) string {
	return m.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ModelDisabled",
			Description: "a requested model record is disabled",
			Other:       "{{.Model}} has been disabled.",
		},
		TemplateData: map[string]interface{}{
			"Model": name,
		},
	})
}

func (t *Translations) InternalServerError() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "InternalServerError",
			Description: "internal server error.",
			Other:       "the server encountered a problem and could not process your request.",
		},
	})
}

func (t *Translations) ConflictError() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "Conflict",
			Description: "data entry conflict, could be unique violation, data race or invalid foreign key, etc",
			Other:       "data input conflict detected.",
		},
	})
}

func (t *Translations) MethodNotAllowed() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "MethodNotAllowed",
			Description: "not allowed http verb for route.",
			Other:       "Method Not Allowed.",
		},
	})
}

func (t *Translations) InvalidCredentials() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "InvalidCredentials",
			Description: "user credentials are not valid.",
			Other:       "invalid credentials.",
		},
	})
}

func (t *Translations) UnauthorizedAccess() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "UnauthorizedAccess",
			Description: "user is not authorized to access a resource.",
			Other:       "unauthorized access.",
		},
	})
}

func (t *Translations) DisabledAccount() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "DisabledAccount",
			Description: "user account has been disabled.",
			Other:       "user account has been disabled. please contact app administrators",
		},
	})
}

func (t *Translations) DeletedAccount() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "DeletedAccount",
			Description: "user account has been deleted.",
			Other:       "user account has been deleted. please make a new account or contact app administrators",
		},
	})
}

func (t *Translations) ProfileCleared() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ProfileCleared",
			Description: "user data has been deleted.",
			Other:       "your data has been deleted from the application.",
		},
	})
}

// JWT

func (t *Translations) JwtExpired() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "JwtExpired",
			Description: "jwt token passed expiration.",
			Other:       "token expired, please login again.",
		},
	})
}

func (t *Translations) TransactionDeclined() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "TransactionDeclined",
			Description: "Transaction declined due to insufficient balance.",
			Other:       "Transaction declined: The debit amount exceeds the available balance in the wallet.",
		},
	})
}

func (t *Translations) ExternalRequestError() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ExternalRequestError",
			Description: "An issue occurred while connecting to an external service. This might be temporary, so please try again later.",
			Other:       "We're having trouble connecting to one of our services right now. This might be temporary, so please try again in a moment. If the issue persists, contact support.",
		},
	})
}
