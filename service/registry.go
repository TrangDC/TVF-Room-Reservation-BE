package service

import (
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/azuread"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/repository"

	"go.uber.org/zap"
)

type Service interface {
	Auth() AuthService
	User() UserService
	Role() RoleService
	Office() OfficeService
	Room() RoomService
	Booking() BookingService
}

type serviceImpl struct {
	authService    AuthService
	UserService    UserService
	RoleService    RoleService
	OfficeService  OfficeService
	RoomService    RoomService
	BookingService BookingService
}

func NewService(azureADOAuthClient azuread.AzureADOAuth, entClient *ent.Client, logger *zap.Logger) Service {
	repoRegistry := repository.NewRepository(entClient)

	return &serviceImpl{
		authService:    NewAuthService(azureADOAuthClient, logger),
		UserService:    NewUserService(repoRegistry, logger),
		RoleService:    NewRoleService(repoRegistry, logger),
		OfficeService:  NewOfficeService(repoRegistry, logger),
		RoomService:    NewRoomService(repoRegistry, logger),
		BookingService: NewBookingService(repoRegistry, logger),
	}
}

func (i serviceImpl) Auth() AuthService {
	return i.authService
}

func (i serviceImpl) User() UserService {
	return i.UserService
}

func (i serviceImpl) Role() RoleService {
	return i.RoleService
}

func (i serviceImpl) Office() OfficeService {
	return i.OfficeService
}

func (i serviceImpl) Room() RoomService {
	return i.RoomService
}

func (i serviceImpl) Booking() BookingService {
	return i.BookingService
}
