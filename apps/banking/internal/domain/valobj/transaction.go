package valobj

type TransactionFilters struct {
	StartDate string `query:"startDate" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `query:"endDate" validate:"omitempty,datetime=2006-01-02"`
	Kind      string `query:"kind" validate:"omitempty,oneof=withdrawal deposit transfer"`
}

type TransactionDirection string

type TransactionKind string

func (tk TransactionKind) String() string {
	return string(tk)
}

func (td TransactionDirection) String() string {
	return string(td)
}

const (
	TransactionKindWithdrawal TransactionKind = "withdrawal"
	TransactionKindDeposit    TransactionKind = "deposit"
	TransactionKindTransfer   TransactionKind = "transfer"
)

const (
	TransactionDirectionIncoming TransactionDirection = "incoming"
	TransactionDirectionOutgoing TransactionDirection = "outgoing"
	TransactionDirectionInternal TransactionDirection = "internal"
)
