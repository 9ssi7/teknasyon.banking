package queries

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/9ssi7/banking/pkg/state"
	"github.com/9ssi7/banking/pkg/token"
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

func NewAuthVerifyAccessHandler(sessionRepo abstracts.SessionRepo) AuthVerifyAccessHandler {
	return func(ctx context.Context, query AuthVerifyAccess) (*AuthVerifyAccessRes, error) {
		var claims *token.UserClaim
		var err error
		if query.IsUnverified {
			claims, err = token.Client().Parse(query.AccessToken)
		} else {
			claims, err = token.Client().VerifyAndParse(query.AccessToken)
		}
		if err != nil {
			return nil, rescode.Failed
		}
		session, err := sessionRepo.FindByIds(ctx, claims.Id, state.GetDeviceId(ctx))
		if err != nil {
			return nil, err
		}
		if !session.IsAccessValid(query.AccessToken, query.IpAddr) {
			return nil, rescode.InvalidAccess
		}
		return &AuthVerifyAccessRes{
			User: claims,
		}, nil
	}
}
