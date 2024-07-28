package repository

import (
	"context"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/room"

	"github.com/google/uuid"
)

type RoomRepository interface {
	GetRoomsByOfficeId(ctx context.Context, filter ent.RoomFilter) ([]*ent.Room, error)
	GetRoom(ctx context.Context, id uuid.UUID) (*ent.Room, error)
	CreateRoom(ctx context.Context, input ent.CreateRoomInput) (*ent.Room, error)
	UpdateRoom(ctx context.Context, input ent.UpdateRoomInput) (*ent.Room, error)
	DeleteRoom(ctx context.Context, id uuid.UUID) error

	//query
	BuildQuery() *ent.RoomQuery
	BuildCount(ctx context.Context, query *ent.RoomQuery) (int, error)
	BuildList(ctx context.Context, query *ent.RoomQuery) ([]*ent.Room, error)
	BuildGet(ctx context.Context, query *ent.RoomQuery) (*ent.Room, error)
}

type roomRepoImpl struct {
	client *ent.Client
}

func NewRoomRepository(client *ent.Client) RoomRepository {
	return &roomRepoImpl{
		client: client,
	}
}

func (rps *roomRepoImpl) BuildCount(ctx context.Context, query *ent.RoomQuery) (int, error) {
	return query.Where(room.IsDeleted(false)).Count(ctx)
}

func (rps *roomRepoImpl) BuildGet(ctx context.Context, query *ent.RoomQuery) (*ent.Room, error) {
	return query.First(ctx)
}

func (rps *roomRepoImpl) BuildList(ctx context.Context, query *ent.RoomQuery) ([]*ent.Room, error) {
	return query.Where(room.IsDeleted(false)).All(ctx)
}

func (rps *roomRepoImpl) BuildQuery() *ent.RoomQuery {
	return rps.client.Room.Query()
}

func (rps *roomRepoImpl) GetRoomsByOfficeId(ctx context.Context, filter ent.RoomFilter) ([]*ent.Room, error) {
	filterOfficeID, err := uuid.Parse(filter.OfficeID)
	if err != nil {
		return nil, err
	}
	return rps.client.Room.Query().Where(room.OfficeID(filterOfficeID)).All(ctx)
}

func (rps *roomRepoImpl) GetRoom(ctx context.Context, id uuid.UUID) (*ent.Room, error) {
	return rps.client.Room.Get(ctx, id)
}

func (rps *roomRepoImpl) CreateRoom(ctx context.Context, input ent.CreateRoomInput) (*ent.Room, error) {
	officeID, err := uuid.Parse(input.OfficeID)
	if err != nil {
		return nil, err
	}

	roomCreate := rps.client.Room.Create().
		SetName(input.Name).
		SetColor(input.Color).
		SetFloor(input.Floor).
		SetOfficeID(officeID)

	if input.Description != nil {
		roomCreate.SetDescription(*input.Description)
	}

	if input.ImageURL != nil {
		roomCreate.SetImageURL(*input.ImageURL)
	}

	return roomCreate.Save(ctx)
}

func (rps *roomRepoImpl) UpdateRoom(ctx context.Context, input ent.UpdateRoomInput) (*ent.Room, error) {
	roomID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, err
	}

	upd := rps.client.Room.UpdateOneID(roomID)
	if input.Name != nil {
		upd.SetName(*input.Name)
	}
	if input.Color != nil {
		upd.SetColor(*input.Color)
	}
	if input.Floor != nil {
		upd.SetFloor(*input.Floor)
	}
	if input.Description != nil {
		upd.SetDescription(*input.Description)
	}
	if input.ImageURL != nil {
		upd.SetImageURL(*input.ImageURL)
	}

	if input.OfficeID != nil {
		officeID, err := uuid.Parse(*input.OfficeID)
		if err != nil {
			return nil, err
		}
		upd.SetOfficeID(officeID)
	}

	return upd.Save(ctx)
}

func (rps *roomRepoImpl) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	// Update the room's IsDeleted field to true
	_, err := rps.client.Room.
		UpdateOneID(id).
		SetIsDeleted(true).
		Save(ctx)

	if err != nil {
		return err
	}

	return err
}
