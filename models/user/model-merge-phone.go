package user

import (
	"fmt"

	"bitbucket.org/sadeemTechnology/backend-validator"

	"github.com/ttacon/libphonenumber"
)

// FIX: use validator in place of this func

// MergePhone handles phone input validation, only accounts for input
// key: "region" if it was provided, otherwise it defaults to "LY"
// manually updating phone by user will set is_verified=false
func (m *Model) MergePhone(v *validator.Validator) {
	m.Phone = v.AssignString("phone", m.Phone, 0, 50)
	if m.Phone != nil {
		region := "LY"
		if v.Data.KeyExists("region") {
			region = v.Data.Values.Get("region")
			found := false
			for k := range libphonenumber.GetSupportedRegions() {
				found = found || k == region
			}
			if !found {
				v.Check(false, "phone", "no matching region for: "+region)
				return
			}
		}
		num, err := libphonenumber.Parse(*m.Phone, region)
		if err != nil {
			v.Check(false, "phone", err.Error())
			return
		}
		nn := fmt.Sprintf("%d", num.GetNationalNumber())
		if region == "LY" && len(nn) != 9 {
			v.Check(false, "phone", "LY phones must be made of 9 numbers")
			return
		}
		*m.Phone = fmt.Sprintf(
			"%d%d",
			num.GetCountryCode(),    // 218
			num.GetNationalNumber(), // 921234567
		)
		m.IsVerified = false
	}
}
