package otp

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ttacon/libphonenumber"
)

func ValidatePhoneNumber(phoneNumber, region string) (*string, error) {

	found := false
	for k := range libphonenumber.GetSupportedRegions() {
		found = found || k == region
	}

	if !found {
		return nil, errors.Errorf("no matching region for: %#v", region)
	}

	num, err := libphonenumber.Parse(phoneNumber, region)
	if err != nil {
		return nil, errors.Errorf("Error: %s", err.Error())
	}

	nn := fmt.Sprintf("%d", num.GetNationalNumber())
	if region == "LY" && len(nn) != 9 {
		return nil, errors.Errorf("LY phones must be made of 9 numbers")
	}

	phone := fmt.Sprintf(
		"%d%d",
		num.GetCountryCode(),    // 218
		num.GetNationalNumber(), // 921234567
	)

	return &phone, nil
}
