package commands

import (
	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/validation"
	"go.opentelemetry.io/otel/trace"
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

	MoneyTransfer MoneyTransferHandler
}

func NewHandler(tracer trace.Tracer, r abstracts.Repositories, v validation.Service) Handlers {
	return Handlers{
		AuthLogin:    NewAuthLoginHandler(tracer, v, r.UserRepo, r.VerifyRepo, r.SessionRepo),
		AuthStart:    NewAuthStartHandler(tracer, v, r.VerifyRepo, r.UserRepo),
		AuthLogout:   NewAuthLogoutHandler(tracer, r.SessionRepo),
		AuthRefresh:  NewAuthRefreshHandler(tracer, r.SessionRepo, r.UserRepo),
		AuthRegister: NewAuthRegisterHandler(tracer, v, r.UserRepo),
		AuthVerify:   NewAuthVerifyHandler(tracer, r.UserRepo),

		AccountCreate:          NewAccountCreateHandler(tracer, v, r.AccountRepo),
		AccountActivate:        NewAccountActivateHandler(tracer, v, r.AccountRepo),
		AccountFreeze:          NewAccountFreezeHandler(tracer, v, r.AccountRepo),
		AccountLock:            NewAccountLockHandler(tracer, v, r.AccountRepo),
		AccountSuspend:         NewAccountSuspendHandler(tracer, v, r.AccountRepo),
		AccountBalanceLoad:     NewAccountBalanceLoadHandler(tracer, v, r.AccountRepo, r.TransactionRepo),
		AccountBalanceWithdraw: NewAccountBalanceWithdrawHandler(tracer, v, r.AccountRepo, r.TransactionRepo),

		MoneyTransfer: NewMoneyTransferHandler(tracer, v, r.UserRepo, r.AccountRepo, r.TransactionRepo),
	}
}
