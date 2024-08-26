package translations

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (t *Translations) TranslateModels() {
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "banner", Other: "Banner"},
	})
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "category", Other: "Category"},
	})
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "document", Other: "Document"},
	})
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "permission", Other: "Permission"},
	})
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "role", Other: "Role"},
	})
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "setting", Other: "Setting"},
	})
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "token", Other: "Token"},
	})
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "user", Other: "User"},
	})
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "user_notification",
			Other: "User Notification",
		},
	})
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "fcm_notification",
			Other: "FCM Notification",
		},
	})
	t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "image",
			Other: "Image",
		},
	})
}
