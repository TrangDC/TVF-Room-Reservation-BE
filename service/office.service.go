package service

import (
	"context"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OfficeService interface {
	GetOffices(ctx context.Context) ([]*ent.OfficeDto, error)
	GetOffice(ctx context.Context, id uuid.UUID) (*ent.OfficeDto, error)
	CreateOffice(ctx context.Context, input ent.CreateOfficeInput) (*ent.OfficeResponse, error)
	UpdateOffice(ctx context.Context, input ent.UpdateOfficeInput) (*ent.OfficeResponse, error)
	DeleteOffice(ctx context.Context, id uuid.UUID) (string, error)
}

type officeSvcImpl struct {
	repoRegistry repository.Repository
	logger       *zap.Logger
}

func NewOfficeService(repoRegistry repository.Repository, logger *zap.Logger) OfficeService {
	return &officeSvcImpl{
		repoRegistry: repoRegistry,
		logger:       logger,
	}
}

func (svc *officeSvcImpl) GetOffices(ctx context.Context) ([]*ent.OfficeDto, error) {
	offices, err := svc.repoRegistry.Office().GetOffices(ctx)
	if err != nil {
		return nil, util.WrapGQLInternalError(ctx)
	}
	var result []*ent.OfficeDto
	for _, office := range offices {
		var rooms []*ent.Room
		for _, room := range office.Edges.Rooms {
			rooms = append(rooms, &ent.Room{
				ID:          room.ID,
				Name:        room.Name,
				Color:       room.Color,
				Floor:       room.Floor,
				OfficeID:    room.OfficeID,
				Description: room.Description,
				ImageURL:    room.ImageURL,
			})
		}
		result = append(result, &ent.OfficeDto{
			ID:          office.ID.String(),
			Name:        office.Name,
			Description: &office.Description,
			Rooms:       rooms,
		})
	}
	return result, nil
}

func (svc *officeSvcImpl) GetOffice(ctx context.Context, id uuid.UUID) (*ent.OfficeDto, error) {
	office, err := svc.repoRegistry.Office().GetOffice(ctx, id)
	if err != nil {
		return nil, util.WrapGQLInternalError(ctx)
	}
	var rooms []*ent.Room
	for _, room := range office.Edges.Rooms {
		rooms = append(rooms, &ent.Room{
			ID:          room.ID,
			Name:        room.Name,
			Color:       room.Color,
			Floor:       room.Floor,
			OfficeID:    room.OfficeID,
			Description: room.Description,
			ImageURL:    room.ImageURL,
		})
	}
	officeDto := &ent.OfficeDto{
		ID:          office.ID.String(),
		Name:        office.Name,
		Description: &office.Description,
		Rooms:       rooms,
	}
	return officeDto, nil
}

func (svc *officeSvcImpl) CreateOffice(ctx context.Context, input ent.CreateOfficeInput) (*ent.OfficeResponse, error) {
	office, err := svc.repoRegistry.Office().CreateOffice(ctx, input)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "Failed to create office.", 500, util.ErrorFlag(err.Error()))
	}
	return &ent.OfficeResponse{
		Message: "Data has been successfully created.",
		Data:    office,
	}, err
}

func (svc *officeSvcImpl) UpdateOffice(ctx context.Context, input ent.UpdateOfficeInput) (*ent.OfficeResponse, error) {
	officeUpdate, err := svc.repoRegistry.Office().UpdateOffice(ctx, input)
	if err != nil {
		return nil, err
	}
	return &ent.OfficeResponse{
		Message: "Data has been successfully updated.",
		Data:    officeUpdate,
	}, err
}

func (svc *officeSvcImpl) DeleteOffice(ctx context.Context, id uuid.UUID) (string, error) {
	errorMessage := svc.repoRegistry.Office().DeleteOffice(ctx, id)
	if errorMessage != nil {
		return "", util.WrapGQLError(ctx, "Failed to delete office.", 500, util.ErrorFlag(errorMessage.Error()))
	}
	return "Data has been successfully deleted.", nil
}
