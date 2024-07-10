package queries

import (
	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/validation"
)

type Handlers struct {
	AuthCheck         AuthCheckHandler
	AuthVerifyAccess  AuthVerifyAccessHandler
	AuthVerifyRefresh AuthVerifyRefreshHandler

	AccountList AccountListHandler

	TransactionList TransactionListHandler
}

func NewHandler(r abstracts.Repositories, v validation.Service) Handlers {
	return Handlers{
		AuthCheck:         NewAuthCheckHandler(r.VerifyRepo),
		AuthVerifyAccess:  NewAuthVerifyAccessHandler(r.SessionRepo),
		AuthVerifyRefresh: NewAuthVerifyRefreshHandler(r.SessionRepo),

		AccountList: NewAccountListHandler(r.AccountRepo),

		TransactionList: NewTransactionListHandler(v, r.TransactionRepo, r.AccountRepo),
	}
}
