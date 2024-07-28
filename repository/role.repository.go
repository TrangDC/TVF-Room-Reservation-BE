package repository

import (
	"context"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"

	"github.com/google/uuid"
)

type RoleRepository interface {
	GetRoles(ctx context.Context) ([]*ent.Role, error)
	GetRole(ctx context.Context, id uuid.UUID) (*ent.Role, error)
}

type roleRepoImpl struct {
	client *ent.Client
}

func NewRoleRepository(client *ent.Client) RoleRepository {
	return &roleRepoImpl{
		client: client,
	}
}

func (rps roleRepoImpl) GetRoles(ctx context.Context) ([]*ent.Role, error) {
	return rps.client.Role.Query().All(ctx)
}

func (rps roleRepoImpl) GetRole(ctx context.Context, id uuid.UUID) (*ent.Role, error) {
	return rps.client.Role.Get(ctx, id)
}
