package queries

import "github.com/9ssi7/banking/internal/domain/abstracts"

type Handlers struct {
	AuthCheck         AuthCheckHandler
	AuthVerifyAccess  AuthVerifyAccessHandler
	AuthVerifyRefresh AuthVerifyRefreshHandler

	AccountList AccountListHandler
}

func NewHandler(r abstracts.Repositories) Handlers {
	return Handlers{
		AuthCheck:         NewAuthCheckHandler(r.VerifyRepo),
		AuthVerifyAccess:  NewAuthVerifyAccessHandler(r.SessionRepo),
		AuthVerifyRefresh: NewAuthVerifyRefreshHandler(r.SessionRepo),

		AccountList: NewAccountListHandler(r.AccountRepo),
	}
}
