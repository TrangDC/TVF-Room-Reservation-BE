package repository

import (
	"context"
	"net/http"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/user"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/userrole"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"

	"github.com/google/uuid"
)

type UserRepository interface {
	// Queries
	BuildQuery() *ent.UserQuery
	BuildList(ctx context.Context, query *ent.UserQuery) ([]*ent.User, error)
	BuildCount(ctx context.Context, query *ent.UserQuery) (int, error)
	BuildGet(ctx context.Context, query *ent.UserQuery) (*ent.User, error)
	BuildUserRoleQuery() *ent.UserRoleQuery
	BuildExist(ctx context.Context, query *ent.UserRoleQuery, userID, roleID uuid.UUID) (bool, error)

	GetUser(ctx context.Context, id uuid.UUID) (*ent.User, error)
	GetUserByEmail(ctx context.Context, workEmail string) (*ent.User, error)

	// Mutations
	BuildUpdateOne(ctx context.Context, model *ent.User) *ent.UserUpdateOne
	BuildSaveUpdateOne(ctx context.Context, update *ent.UserUpdateOne) (*ent.User, error)

	AssignRole(ctx context.Context, userID, roleID uuid.UUID) (*ent.User, error)
	UnassignRole(ctx context.Context, userID, roleID uuid.UUID) (*ent.User, error)
}

type userRepoImpl struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) UserRepository {
	return &userRepoImpl{
		client: client,
	}
}

// Base functions
func (rps *userRepoImpl) BuildQuery() *ent.UserQuery {
	return rps.client.User.Query().WithUserRoles()
}

func (rps *userRepoImpl) BuildGet(ctx context.Context, query *ent.UserQuery) (*ent.User, error) {
	return query.First(ctx)
}

func (rps *userRepoImpl) BuildUserRoleCreate() *ent.UserRoleCreate {
	return rps.client.UserRole.Create()
}

func (rps *userRepoImpl) BuildUserRoleDelete() *ent.UserRoleDelete {
	return rps.client.UserRole.Delete()
}

func (rps *userRepoImpl) BuildUpdate() *ent.UserUpdate {
	return rps.client.User.Update()
}

func (rps *userRepoImpl) BuildList(ctx context.Context, query *ent.UserQuery) ([]*ent.User, error) {
	return query.All(ctx)
}

func (rps *userRepoImpl) BuildCount(ctx context.Context, query *ent.UserQuery) (int, error) {
	return query.Count(ctx)
}

func (rps *userRepoImpl) BuildUpdateOne(ctx context.Context, model *ent.User) *ent.UserUpdateOne {
	return model.Update()
}

func (rps *userRepoImpl) BuildSaveUpdateOne(ctx context.Context, update *ent.UserUpdateOne) (*ent.User, error) {
	return update.Save(ctx)
}

func (rps *userRepoImpl) BuildUserRoleQuery() *ent.UserRoleQuery {
	return rps.client.UserRole.Query()
}

func (rps *userRepoImpl) BuildExist(ctx context.Context, query *ent.UserRoleQuery, userID, roleID uuid.UUID) (bool, error) {
	return query.
		Where(
			userrole.And(
				userrole.UserIDEQ(userID),
				userrole.RoleIDEQ(roleID),
			),
		).
		Exist(ctx)
}

// Query functions
func (rps userRepoImpl) GetUser(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	query := rps.BuildQuery().Where(user.IDEQ(id))
	return rps.BuildGet(ctx, query)
}

func (rps userRepoImpl) GetUserByEmail(ctx context.Context, workEmail string) (*ent.User, error) {
	query := rps.BuildQuery().Where(user.WorkEmailEQ(workEmail))
	return rps.BuildGet(ctx, query)
}

// Mutation functions
func (rps userRepoImpl) AssignRole(ctx context.Context, userID, roleID uuid.UUID) (*ent.User, error) {
	query := rps.BuildUserRoleQuery()

	// Check if the role is already assigned to the user
	exists, err := rps.BuildExist(ctx, query, userID, roleID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "error checking existing role", http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}
	if exists {
		return nil, util.WrapGQLError(ctx, "role already assigned to user", http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}

	// Assign role to the user
	_, err = rps.BuildUserRoleCreate().
		SetUserID(userID).
		SetRoleID(roleID).
		Save(ctx)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "failed to assign role", http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}

	// Fetch and return the updated user with roles
	user, err := rps.GetUser(ctx, userID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}

	return user, err
}

func (rps userRepoImpl) UnassignRole(ctx context.Context, userID, roleID uuid.UUID) (*ent.User, error) {
	query := rps.BuildUserRoleQuery()

	// Check if the role is already assigned to the user
	exists, err := rps.BuildExist(ctx, query, userID, roleID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "error checking existing role", http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}
	if !exists {
		return nil, util.WrapGQLError(ctx, "role has not assigned to user", http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}

	// Unassign the role from the user
	_, err = rps.BuildUserRoleDelete().
		Where(
			userrole.And(
				userrole.UserID(userID),
				userrole.RoleIDEQ(roleID),
			),
		).
		Exec(ctx)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "failed to unassign role", http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}

	// Fetch and return the updated user with roles
	user, err := rps.GetUser(ctx, userID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "failed to fetch user", http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}

	return user, err
}
