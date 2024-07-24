package queries

import (
	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/validation"
	"go.opentelemetry.io/otel/trace"
)

type Handlers struct {
	AuthCheck         AuthCheckHandler
	AuthVerifyAccess  AuthVerifyAccessHandler
	AuthVerifyRefresh AuthVerifyRefreshHandler

	AccountList AccountListHandler

	TransactionList TransactionListHandler
}

func NewHandler(tracer trace.Tracer, r abstracts.Repositories, v validation.Service) Handlers {
	return Handlers{
		AuthCheck:         NewAuthCheckHandler(tracer, r.VerifyRepo),
		AuthVerifyAccess:  NewAuthVerifyAccessHandler(tracer, r.SessionRepo),
		AuthVerifyRefresh: NewAuthVerifyRefreshHandler(tracer, r.SessionRepo),

		AccountList: NewAccountListHandler(tracer, r.AccountRepo),

		TransactionList: NewTransactionListHandler(tracer, v, r.TransactionRepo, r.AccountRepo),
	}
}
