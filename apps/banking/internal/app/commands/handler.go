package commands

import (
	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/validation"
)

type Handlers struct {
	AuthLogin    AuthLoginHandler
	AuthRegister AuthRegisterHandler
	AuthStart    AuthStartHandler
	AuthRefresh  AuthRefreshHandler
	AuthLogout   AuthLogoutHandler
	AuthVerify   AuthVerifyHandler

	AccountCreate          AccountCreateHandler
	AccountActivate        AccountActivateHandler
	AccountFreeze          AccountFreezeHandler
	AccountLock            AccountLockHandler
	AccountSuspend         AccountSuspendHandler
	AccountBalanceLoad     AccountBalanceLoadHandler
	AccountBalanceWithdraw AccountBalanceWithdrawHandler
}

func NewHandler(r abstracts.Repositories, v validation.Service) Handlers {
	return Handlers{
		AuthLogin:    NewAuthLoginHandler(v, r.UserRepo, r.VerifyRepo, r.SessionRepo),
		AuthStart:    NewAuthStartHandler(v, r.VerifyRepo, r.UserRepo),
		AuthLogout:   NewAuthLogoutHandler(r.SessionRepo),
		AuthRefresh:  NewAuthRefreshHandler(r.SessionRepo, r.UserRepo),
		AuthRegister: NewAuthRegisterHandler(v, r.UserRepo),
		AuthVerify:   NewAuthVerifyHandler(r.UserRepo),

		AccountCreate:          NewAccountCreateHandler(r.AccountRepo),
		AccountActivate:        NewAccountActivateHandler(r.AccountRepo),
		AccountFreeze:          NewAccountFreezeHandler(r.AccountRepo),
		AccountLock:            NewAccountLockHandler(r.AccountRepo),
		AccountSuspend:         NewAccountSuspendHandler(r.AccountRepo),
		AccountBalanceLoad:     NewAccountBalanceLoadHandler(r.AccountRepo),
		AccountBalanceWithdraw: NewAccountBalanceWithdrawHandler(r.AccountRepo),
	}
}
