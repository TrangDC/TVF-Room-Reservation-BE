package service

import (
	"context"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RoleService interface {
	GetRoles(ctx context.Context) ([]*ent.Role, error)
	GetRole(ctx context.Context, id uuid.UUID) (*ent.Role, error)
}

type roleSvcImpl struct {
	repoRegistry repository.Repository
	logger       *zap.Logger
}

func NewRoleService(repoRegistry repository.Repository, logger *zap.Logger) RoleService {
	return &roleSvcImpl{
		repoRegistry: repoRegistry,
		logger:       logger,
	}
}

func (svc *roleSvcImpl) GetRoles(ctx context.Context) ([]*ent.Role, error) {
	return svc.repoRegistry.Role().GetRoles(ctx)
}

func (svc *roleSvcImpl) GetRole(ctx context.Context, id uuid.UUID) (*ent.Role, error) {
	return svc.repoRegistry.Role().GetRole(ctx, id)
}
