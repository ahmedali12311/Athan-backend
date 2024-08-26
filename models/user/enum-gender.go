package user

type GenderValue string

const (
	Male   GenderValue = "male"
	Female GenderValue = "female"
)

var AllGenders = []string{
	Male.String(),
	Female.String(),
}

func (v GenderValue) String() string {
	switch v {
	case Male:
		return "male"
	case Female:
		return "female"
	default:
		return ""
	}
}
