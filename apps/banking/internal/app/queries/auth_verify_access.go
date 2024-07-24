package queries

import (
	"context"
	"errors"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/9ssi7/banking/pkg/state"
	"github.com/9ssi7/banking/pkg/token"
	"go.opentelemetry.io/otel/trace"
)

type AuthVerifyAccess struct {
	AccessToken  string
	IpAddr       string
	IsUnverified bool
}

type AuthVerifyAccessRes struct {
	User *token.UserClaim
}

type AuthVerifyAccessHandler cqrs.HandlerFunc[AuthVerifyAccess, *AuthVerifyAccessRes]

func NewAuthVerifyAccessHandler(tracer trace.Tracer, sessionRepo abstracts.SessionRepo) AuthVerifyAccessHandler {
	return func(ctx context.Context, query AuthVerifyAccess) (*AuthVerifyAccessRes, error) {
		ctx, span := tracer.Start(ctx, "AuthVerifyAccessHandler")
		defer span.End()
		var claims *token.UserClaim
		var err error
		if query.IsUnverified {
			claims, err = token.Client().Parse(query.AccessToken)
		} else {
			claims, err = token.Client().VerifyAndParse(query.AccessToken)
		}
		if err != nil {
			return nil, rescode.Failed(err)
		}
		session, notExists, err := sessionRepo.FindByIds(ctx, claims.Id, state.GetDeviceId(ctx))
		if err != nil {
			return nil, err
		}
		if notExists {
			return nil, rescode.InvalidAccess(errors.New("invalid access with token and ip"))
		}
		if !session.IsAccessValid(query.AccessToken, query.IpAddr) {
			return nil, rescode.InvalidAccess(errors.New("invalid access with token and ip"))
		}
		return &AuthVerifyAccessRes{
			User: claims,
		}, nil
	}
}
