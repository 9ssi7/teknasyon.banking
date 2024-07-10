package currency

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	TRY Currency = "TRY"
)

func (c Currency) String() string {
	return string(c)
}

func IsValid(c string) bool {
	switch c {
	case USD.String(), EUR.String(), TRY.String():
		return true
	}
	return false
}
