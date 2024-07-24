package commands

import (
	"context"
	"errors"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/9ssi7/banking/pkg/state"
	"github.com/9ssi7/banking/pkg/token"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type AuthRefresh struct {
	AccessToken  string
	RefreshToken string
	IpAddress    string
	UserId       uuid.UUID
}

type AuthRefreshRes struct {
	AccessToken string
}

type AuthRefreshHandler cqrs.HandlerFunc[AuthRefresh, *AuthRefreshRes]

func NewAuthRefreshHandler(tracer trace.Tracer, sessionRepo abstracts.SessionRepo, userRepo abstracts.UserRepo) AuthRefreshHandler {
	return func(ctx context.Context, cmd AuthRefresh) (*AuthRefreshRes, error) {
		ctx, span := tracer.Start(ctx, "AuthRefreshHandler")
		defer span.End()
		session, notFound, err := sessionRepo.FindByIds(ctx, cmd.UserId, state.GetDeviceId(ctx))
		if err != nil {
			return nil, err
		}
		if notFound {
			return nil, rescode.InvalidRefreshOrAccessTokens(errors.New("invalid refresh with access token and ip"))
		}
		if !session.IsRefreshValid(cmd.AccessToken, cmd.RefreshToken, cmd.IpAddress) {
			return nil, rescode.InvalidRefreshOrAccessTokens(errors.New("invalid refresh with access token and ip"))
		}
		user, err := userRepo.FindById(ctx, cmd.UserId)
		if err != nil {
			return nil, err
		}
		accessToken, err := token.Client().GenerateAccessToken(token.User{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Roles: user.Roles,
		})
		if err != nil {
			return nil, err
		}
		session.Refresh(accessToken)
		if err := sessionRepo.Save(ctx, user.Id, session); err != nil {
			return nil, err
		}
		return &AuthRefreshRes{
			AccessToken: accessToken,
		}, nil
	}
}
