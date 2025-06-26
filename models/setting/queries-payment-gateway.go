package setting

// func (m *Queries) GetForPaymentGateway(s *payment_gateway.Settings) error {
// 	keyVals, err := m.GetByKeys(
// 		[]string{
// 			KeyPaymentGatewayAPIKey,
// 			KeyPaymentGatewayEndpoint,
// 		},
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	for _, v := range keyVals {
// 		switch v.Key {
// 		case KeyPaymentGatewayAPIKey:
// 			s.APIKey = v.Value
// 		case KeyPaymentGatewayEndpoint:
// 			s.Endpoint = v.Value
// 		}
// 	}
// 	return nil
// }
