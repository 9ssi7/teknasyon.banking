package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/state"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type AuthLogout struct {
	UserId uuid.UUID `json:"-"`
}

type AuthLogoutHandler cqrs.HandlerFunc[AuthLogout, *cqrs.Empty]

func NewAuthLogoutHandler(tracer trace.Tracer, sessionRepo abstracts.SessionRepo) AuthLogoutHandler {
	return func(ctx context.Context, cmd AuthLogout) (*cqrs.Empty, error) {
		ctx, span := tracer.Start(ctx, "AuthLogoutHandler")
		defer span.End()
		err := sessionRepo.Destroy(ctx, cmd.UserId, state.GetDeviceId(ctx))
		if err != nil {
			return nil, err
		}
		return &cqrs.Empty{}, nil
	}
}
