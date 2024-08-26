package setting

import (
	"app/pkg/tlync"
)

func (m *Queries) GetForTlync(s *tlync.Settings) error {
	keyVals, err := m.GetByKeys(
		[]string{
			KeyTlyncToken,
			KeyTlyncStoreID,
			KeyTlyncEndpoint,
			KeyTlyncFrontUrl,
		},
	)
	if err != nil {
		return err
	}
	for _, v := range keyVals {
		switch v.Key {
		case KeyTlyncToken:
			s.Token = v.Value
		case KeyTlyncStoreID:
			s.StoreID = v.Value
		case KeyTlyncEndpoint:
			s.Endpoint = v.Value
		case KeyTlyncFrontUrl:
			s.FrontURL = v.Value
		}
	}
	return nil
}
