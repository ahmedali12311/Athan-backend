package translations

import "github.com/nicksnyder/go-i18n/v2/i18n"

func (t *Translations) FileIsNotAnImage() string {
	return t.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "FileIsNotAnImage",
			Description: "",
			Other:       "",
		},
	})
}

func (t *Translations) LoggedOut() string {
	return ""
}

func (t *Translations) NotLoggedIn() string {
	return ""
}

func (t *Translations) UserAlreadyVerified() string {
	return ""
}

func (t *Translations) ValidateAlphanumericDashUnderscoreCharactersOnly() string { //nolint:lll
	return ""
}

func (t *Translations) ValidateMinChar(value int) string {
	return ""
}

func (t *Translations) ValidateMaxChar(value int) string {
	return ""
}

func (t *Translations) ValidateStartWithLetter() string {
	return ""
}
