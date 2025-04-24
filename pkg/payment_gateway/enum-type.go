package payment_gateway

type TypeValue string

const (
	TypeCredit TypeValue = "credit"
	TypeDebit  TypeValue = "debit"
)

var AllPaymentPaymentWalletTrxTypes = []string{
	TypeCredit.String(),
	TypeDebit.String(),
}

func (v TypeValue) String() string {
	switch v {
	case TypeCredit:
		return "credit"
	case TypeDebit:
		return "debit"
	default:
		return ""
	}
}
