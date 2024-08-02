package service

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/role"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/user"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/userrole"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService interface {
	// Queries
	GetUser(ctx context.Context, id uuid.UUID) (*ent.User, error)
	GetMe(ctx context.Context) (*ent.UserData, error)
	GetAdminUsers(ctx context.Context, pagination *ent.PaginationInput, keyword *string) (*ent.UserDataResponse, error)

	// Mutations
	AssignRole(ctx context.Context, input ent.AssignRoleInput) (*ent.UserResponse, error)
	UnassignRole(ctx context.Context, input ent.UnassignRoleInput) (*ent.UserResponse, error)

	validateAssignRole(ctx context.Context, workEmail, roleID string) (uuid.UUID, uuid.UUID, error)
	renderUserResponse(ctx context.Context, user *ent.User, mode string) *ent.UserResponse
}

type userSvcImpl struct {
	repoRegistry repository.Repository
	logger       *zap.Logger
}

func NewUserService(repoRegistry repository.Repository, logger *zap.Logger) UserService {
	return &userSvcImpl{
		repoRegistry: repoRegistry,
		logger:       logger,
	}
}

// Query functions
func (svc *userSvcImpl) GetUser(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	return svc.repoRegistry.User().GetUser(ctx, id)
}

func (svc *userSvcImpl) GetMe(ctx context.Context) (*ent.UserData, error) {
	userID := ctx.Value("user_id").(uuid.UUID)
	user, err := svc.repoRegistry.User().GetUser(ctx, userID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusNotFound, util.ErrorFlagNotFound)
	}

	roles := user.QueryUserRoles().QueryRole().AllX(ctx)

	return &ent.UserData{
		ID:        user.ID.String(),
		Name:      user.Name,
		WorkEmail: user.WorkEmail,
		Roles:     sortRoles(roles),
	}, nil
}

func (svc *userSvcImpl) GetAdminUsers(ctx context.Context, pagination *ent.PaginationInput, keyword *string) (*ent.UserDataResponse, error) {
	var results *ent.UserDataResponse
	var err error

	query := svc.repoRegistry.User().BuildQuery()

	// Filter users by their work email address or username and admin role
	svc.filter(query, keyword)

	total, err := svc.repoRegistry.User().BuildCount(ctx, query)
	if err != nil {
		return results, util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	if pagination != nil {
		query = query.Limit(*pagination.PerPage).Offset((*pagination.Page - 1) * *pagination.PerPage)
	}

	users, err := svc.repoRegistry.User().BuildList(ctx, query)
	if err != nil {
		return results, util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	var userData []*ent.UserData
	for _, user := range users {
		roles := user.QueryUserRoles().QueryRole().AllX(ctx)

		userData = append(userData, &ent.UserData{
			ID:        user.ID.String(),
			Name:      user.Name,
			WorkEmail: user.WorkEmail,
			Roles:     sortRoles(roles),
		})
	}

	return &ent.UserDataResponse{
		Total: total,
		Data:  userData,
	}, nil
}

// Mutation functions
func (svc *userSvcImpl) AssignRole(ctx context.Context, input ent.AssignRoleInput) (*ent.UserResponse, error) {
	// Validate input
	userID, roleID, err := svc.validateAssignRole(ctx, input.WorkEmail, input.RoleID)
	if err != nil {
		return nil, err
	}

	userRole := svc.repoRegistry.User().BuildUserRoleQuery()

	// Check if user already has the role
	exists, err := svc.repoRegistry.User().BuildExist(ctx, userRole, userID, roleID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "error checking existing role", http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}

	// Check if the role is already assigned to the user
	if exists {
		return nil, util.WrapGQLError(ctx, "user already has that role", http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}

	// Assign role
	updatedUser, err := svc.repoRegistry.User().AssignRole(ctx, userID, roleID)
	if err != nil {
		return nil, err
	}

	return svc.renderUserResponse(ctx, updatedUser, "assigned"), nil
}

func (svc *userSvcImpl) UnassignRole(ctx context.Context, input ent.UnassignRoleInput) (*ent.UserResponse, error) {
	// Validate input
	userID, roleID, err := svc.validateAssignRole(ctx, input.WorkEmail, input.RoleID)
	if err != nil {
		return nil, err
	}

	// Unassign role
	updatedUser, err := svc.repoRegistry.User().UnassignRole(ctx, userID, roleID)
	if err != nil {
		return nil, err
	}

	return svc.renderUserResponse(ctx, updatedUser, "unassigned"), err
}

func (svc *userSvcImpl) validateAssignRole(ctx context.Context, workEmail, roleID string) (uuid.UUID, uuid.UUID, error) {
	// Fetch the user by work email
	user, err := svc.repoRegistry.User().GetUserByEmail(ctx, workEmail)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, util.WrapGQLError(ctx, "user not found by that email", http.StatusNotFound, util.ErrorFlagNotFound)
	}
	userID := user.ID

	// Parse the role input
	parseRoleID, err := uuid.Parse(roleID)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, util.WrapGQLError(ctx, "invalid role ID", http.StatusBadRequest, util.ErrorFlagValidateFail)
	}
	// Check if the role is existing
	_, err = svc.repoRegistry.Role().GetRole(ctx, parseRoleID)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, util.WrapGQLError(ctx, "role not found", http.StatusNotFound, util.ErrorFlagNotFound)
	}

	return userID, parseRoleID, err
}

// Common functions
func (svc *userSvcImpl) renderUserResponse(ctx context.Context, user *ent.User, mode string) *ent.UserResponse {
	roles := user.QueryUserRoles().QueryRole().AllX(ctx)

	return &ent.UserResponse{
		Message: "Role " + mode + " successfully",
		User: &ent.UserData{
			ID:        user.ID.String(),
			Name:      user.Name,
			WorkEmail: user.WorkEmail,
			Roles:     sortRoles(roles),
		},
	}
}

func (svc *userSvcImpl) filter(userQuery *ent.UserQuery, keyword *string) {
	// Filter users by admin role
	userQuery.Where(
		user.HasUserRolesWith(
			userrole.HasRoleWith(
				role.MachineNameEQ("administrator"),
			),
		),
	)

	// Filter users by name or work email if provided
	if keyword != nil && *keyword != "" {
		trimmedKeyword := strings.TrimSpace(*keyword)

		if len(trimmedKeyword) >= 2 {
			userQuery.Where(
				user.Or(
					user.NameContainsFold(trimmedKeyword),
					user.WorkEmailContainsFold(trimmedKeyword),
				),
			)
		}
	}
}

func sortRoles(roles []*ent.Role) []*ent.Role {
	order := map[string]int{
		"super_admin":   1,
		"administrator": 2,
		"user":          3,
	}
	sort.SliceStable(roles, func(i, j int) bool {
		return order[roles[i].MachineName] < order[roles[j].MachineName]
	})

	return roles
}
