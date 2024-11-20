package user

import (
	"app/pkg/validator"
)

func (m *Model) MergeEmailPassword(v *validator.Validator, isInsert, isResetPassword bool) {
	var oldpw string
	data := v.Data

	keyEmail := "email"
	keyPW := "password"
	keyPWC := "password_confirmation"
	keyOLDPW := "previous_password"

	hasEmail := data.KeyExists(keyEmail)
	hasPW := data.KeyExists(keyPW)
	hasPWC := data.KeyExists(keyPWC)
	hasOLDPW := data.KeyExists(keyOLDPW)

	if !hasPW && hasPWC {
		v.Check(false, keyPW, v.T.ValidateRequired())
		return
	}

	if hasPW && !hasPWC {
		v.Check(false, keyPWC, v.T.ValidateRequired())
		return
	}

	if hasEmail {
		if mEmail := data.Values.Get(keyEmail); mEmail != "" {
			m.Email = &mEmail
		}
		if isInsert && !hasPW {
			v.Check(false, keyPW, v.T.ValidateRequired())
			v.Check(false, keyPWC, v.T.ValidateRequired())
			return
		}
	}
	if !isInsert && !hasOLDPW && hasPW && !isResetPassword {
		v.Check(false, keyOLDPW, v.T.ValidateRequired())
		return
	}

	if hasPW {
		pw := data.Values.Get(keyPW)
		if pw == "" {
			v.Check(false, keyPW, v.T.ValidateRequired())
			return
		}

		pwc := data.Values.Get(keyPWC)
		if pwc == "" {
			v.Check(false, keyPWC, v.T.ValidateRequired())
			return
		}

		if pw != pwc {
			v.Check(
				false,
				keyPW,
				v.T.ValidatePasswordConfirmationNoMatch(),
			)
			return
		}

		if !isInsert && hasOLDPW && !isResetPassword {
			oldpw = data.Values.Get(keyOLDPW)
			if oldpw == "" {
				v.Check(false, keyOLDPW, v.T.ValidateRequired())
				return
			}
		}
		m.Password = password{
			Plaintext:    &pw,
			Confirmation: &pwc,
			Previous:     &oldpw,
		}
		if err := m.Password.Set(pw); err != nil {
			v.Check(false, keyPW, err.Error())
		}
		if !isInsert && !isResetPassword {
			err := m.Password.CheckPreviousPassword(v, m.ID.String())
			if err != nil {
				v.Check(false, keyOLDPW, err.Error())
			}
		}

		if pw != "" && m.Password.Hash == nil {
			panic("missing password Hash for user")
		}
	}
}
