package translations

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (t *Translations) TranslateModels() {
	messages := []i18n.Message{
		{ID: "banner", Other: "Banner"},
		{ID: "category", Other: "Category"},
		{ID: "profile", Other: "Profile"},
		{ID: "document", Other: "Document"},
		{ID: "permission", Other: "Permission"},
		{ID: "role", Other: "Role"},
		{ID: "setting", Other: "Setting"},
		{ID: "token", Other: "Token"},
		{ID: "user", Other: "User"},
		{ID: "user_notification", Other: "User Notification"},
		{ID: "fcm_notification", Other: "FCM Notification"},
		{ID: "image", Other: "Image"},
		{ID: "wallet_transaction", Other: "Wallet Transaction"},
		{ID: "city", Other: "City"},
		{ID: "prayer_times", Other: "Prayer Times"},
	}

	for i := range messages {
		t.Localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &(messages[i]),
		})
	}
}
