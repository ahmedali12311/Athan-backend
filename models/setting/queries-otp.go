package setting

import (
	"app/pkg/otp"
)

func (m *Queries) GetForOTP() (*otp.Settings, error) {
	keyVals, err := m.GetByKeys(
		[]string{
			otp.KeySadeemOTPKey,
			otp.KeySadeemOTPURL,
			otp.KeySadeemOTPJWT,
		},
	)
	if err != nil {
		return nil, err
	}
	var s otp.Settings
	for _, v := range keyVals {
		switch v.Key {
		case otp.KeySadeemOTPKey:
			s.Key = v.Value
		case otp.KeySadeemOTPURL:
			s.URL = v.Value
		case otp.KeySadeemOTPJWT:
			s.JWT = v.Value
		}
	}
	return &s, nil
}
