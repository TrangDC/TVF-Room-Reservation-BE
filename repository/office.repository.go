package repository

import (
	"context"
	"net/http"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/office"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/room"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"

	"github.com/google/uuid"
)

type OfficeRepository interface {
	GetOffices(ctx context.Context) ([]*ent.Office, error)
	GetOffice(ctx context.Context, id uuid.UUID) (*ent.Office, error)
	CreateOffice(ctx context.Context, input ent.CreateOfficeInput) (*ent.Office, error)
	UpdateOffice(ctx context.Context, input ent.UpdateOfficeInput) (*ent.Office, error)
	DeleteOffice(ctx context.Context, id uuid.UUID) error
}

type officeRepoImpl struct {
	client *ent.Client
}

func NewOfficeRepository(client *ent.Client) OfficeRepository {
	return &officeRepoImpl{
		client: client,
	}
}

func (rps *officeRepoImpl) GetOffices(ctx context.Context) ([]*ent.Office, error) {
    // Query all offices with their related rooms where isDeleted is false
    offices, err := rps.client.Office.
        Query().
        WithRooms(func(query *ent.RoomQuery) {
            query.Where(room.IsDeletedEQ(false)) // Ensure correct filter method is used
        }).
        All(ctx)
    if err != nil {
        return nil, err
    }
    return offices, nil
}

func (rps *officeRepoImpl) GetOffice(ctx context.Context, id uuid.UUID) (*ent.Office, error) {
	office, err := rps.client.Office.
        Query().
        Where(office.IDEQ(id)).
        WithRooms(func(query *ent.RoomQuery) {
            query.Where(room.IsDeletedEQ(false)) // Ensure correct filter method is used
        }).
        Only(ctx)
    if err != nil {
        return nil, err
    }
    return office, nil
}

func (rps *officeRepoImpl) CreateOffice(ctx context.Context, input ent.CreateOfficeInput) (*ent.Office, error) {
	exists, err := rps.client.Office.
		Query().
		Where(office.Name(input.Name)).
		Exist(ctx)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "Office with this name already exists", 500, util.ErrorFlag(err.Error()))
	}

	if exists {
		return nil, util.WrapGQLError(ctx, "Office with this name already exists", 200, util.ErrorFlag("Invalid input provided"))
	}

	newOffice := rps.client.Office.
		Create().
		SetName(input.Name)

	if input.Description != nil {
		newOffice.SetDescription(*input.Description)
	}

	return newOffice.Save(ctx)
}

func (rps *officeRepoImpl) UpdateOffice(ctx context.Context, input ent.UpdateOfficeInput) (*ent.Office, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "failed to parse office id", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	// Check existence of the office
	_, err = rps.client.Office.Get(ctx, id)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "Office not found", http.StatusNotFound, util.ErrorFlagValidateFail)
	}

	exists, err := rps.client.Office.
		Query().
		Where(office.IDNEQ(id)).
		Exist(ctx)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "Error query office database", 500, util.ErrorFlag(err.Error()))
	}

	if exists {
		return nil, util.WrapGQLError(ctx, "Office with this name already exists", 200, "Invalid input provided")
	}

	updatedOffice := rps.client.Office.UpdateOneID(id)

	if input.Name != nil {
		updatedOffice.SetName(*input.Name)
	}

	if input.Description != nil {
		updatedOffice.SetDescription(*input.Description)
	}

	return updatedOffice.Save(ctx)
}

func (rps *officeRepoImpl) DeleteOffice(ctx context.Context, id uuid.UUID) error {
	office, err := rps.client.Office.Get(ctx, id)
	if err != nil {
		return util.WrapGQLError(ctx, "office not found", http.StatusNotFound, util.ErrorFlagCanNotDelete)
	}

	if err := rps.client.Office.DeleteOne(office).Exec(ctx); err != nil {
		return util.WrapGQLError(ctx, "failed to delete office", http.StatusBadRequest, util.ErrorFlagCanNotDelete)
	}

	return nil
}
