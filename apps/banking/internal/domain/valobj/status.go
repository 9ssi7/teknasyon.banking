package valobj

type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "active"
	AccountStatusLocked    AccountStatus = "locked"
	AccountStatusFrozen    AccountStatus = "frozen"
	AccountStatusSuspended AccountStatus = "suspended"
)
