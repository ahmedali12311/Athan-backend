package wallet_transaction

type TypeValue string

const (
	TypeCredit TypeValue = "credit"
	TypeDebit  TypeValue = "debit"
)

var AllWalletTrxTypes = []string{
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
