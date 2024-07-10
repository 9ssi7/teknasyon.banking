package valobj

type AccountStatus string

func (s AccountStatus) String() string {
	return string(s)
}

const (
	AccountStatusActive    AccountStatus = "active"
	AccountStatusLocked    AccountStatus = "locked"
	AccountStatusFrozen    AccountStatus = "frozen"
	AccountStatusSuspended AccountStatus = "suspended"
)
