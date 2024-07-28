package service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/room"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	errInvalidRoomID    = "invalid room ID"
	errInvalidRoomName  = "invalid room name"
	errInvalidRoomColor = "invalid room color"
	errInvalidOfficeID  = "invalid office ID"

	errRoomNameRequired  = "room name is required"
	errRoomColorRequired = "room color is required"
	errOfficeIDRequired  = "office ID is required"
	errStartDateRequired = "start date is required"
	errEndDateRequired   = "end date is required when isRepeat is true"

	errParseDate = "error parsing date"
	errParseTime = "error parsing time"

	roomCreatedSuccess = "Data has been successfully created."
	roomUpdatedSuccess = "Data has been successfully updated."
	roomDeletedSuccess = "Data has been successfully deleted."

	startBusinessHour = "09:00"
	endBusinessHour   = "17:30"

	PaginationLimitDefault = 10
	PaginationPageDefault  = 1
)

type RoomService interface {
	GetRooms(ctx context.Context, pagination *ent.PaginationInput, filter *ent.RoomFilter) (*ent.RoomDataResponse, error)
	GetRoom(ctx context.Context, id uuid.UUID) (*ent.Room, error)
	GetAvailableRooms(ctx context.Context, input ent.GetAvailableRoomInput) ([]*ent.AvailableRoomResponse, error)
	CreateRoom(ctx context.Context, input ent.CreateRoomInput) (*ent.RoomResponse, error)
	UpdateRoom(ctx context.Context, input ent.UpdateRoomInput) (*ent.RoomResponse, error)
	DeleteRoom(ctx context.Context, id uuid.UUID) (string, error)
}

type roomSvcImpl struct {
	repoRegistry repository.Repository
	logger       *zap.Logger
}

func NewRoomService(repoRegistry repository.Repository, logger *zap.Logger) RoomService {
	return &roomSvcImpl{
		repoRegistry: repoRegistry,
		logger:       logger,
	}
}

func (svc *roomSvcImpl) GetRooms(ctx context.Context, pagination *ent.PaginationInput, filter *ent.RoomFilter) (*ent.RoomDataResponse, error) {
	// Initialize the variables to store the results and error
	var results *ent.RoomDataResponse
	var err error

	// Build the query to retrieve the rooms from the repository
	query := svc.repoRegistry.Room().BuildQuery()
	svc.applyRoomFilter(query, filter)

	// Add search term filter if provided
	svc.applySearchTermFilter(query, filter)

	// Retrieve the total count of rooms from the query
	total, err := svc.repoRegistry.Room().BuildCount(ctx, query)
	if err != nil {
		return results, util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	// Apply pagination to the query if it is provided
	if pagination != nil {
		query = query.Limit(*pagination.PerPage).Offset((*pagination.Page - 1) * *pagination.PerPage)
	}

	// Retrieve the actual list of rooms from the query
	records, err := svc.repoRegistry.Room().BuildList(ctx, query)
	if err != nil {
		return results, util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	// Return the results with the total count and the actual list of rooms
	return &ent.RoomDataResponse{
		Total: total,
		Data:  records,
	}, nil
}

func (svc *roomSvcImpl) GetRoom(ctx context.Context, id uuid.UUID) (*ent.Room, error) {
	// get room
	room, err := svc.repoRegistry.Room().GetRoom(ctx, id)
	if err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	if room.IsDeleted {
		return nil, util.WrapGQLError(ctx, "room with id "+id.String()+" not found", http.StatusInternalServerError, util.ErrorFlagInternalError)
	}
	return room, nil
}

func (svc *roomSvcImpl) CreateRoom(ctx context.Context, input ent.CreateRoomInput) (*ent.RoomResponse, error) {
	//validate
	if err := svc.validateCreateRoomInput(input); err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusBadRequest, util.ErrorFlagValidateFail)
	}

	// validate office
	if err := svc.validateOfficeID(ctx, input.OfficeID); err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusBadRequest, util.ErrorFlagValidateFail)
	}

	// create room in the database
	room, err := svc.repoRegistry.Room().CreateRoom(ctx, input)
	if err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	return &ent.RoomResponse{
		Message: roomCreatedSuccess,
		Data:    room,
	}, nil
}

func (svc *roomSvcImpl) UpdateRoom(ctx context.Context, input ent.UpdateRoomInput) (*ent.RoomResponse, error) {
	// validate
	if err := svc.validateUpdateRoomInput(input); err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusBadRequest, util.ErrorFlagValidateFail)
	}

	// validate office
	if input.OfficeID != nil {
		if err := svc.validateOfficeID(ctx, *input.OfficeID); err != nil {
			return nil, util.WrapGQLError(ctx, err.Error(), http.StatusBadRequest, util.ErrorFlagValidateFail)
		}
	}

	//check room deleted
	roomId, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusBadRequest, util.ErrorFlagValidateFail)
	}
	foundedRoom, _ := svc.repoRegistry.Room().GetRoom(ctx, roomId)
	if foundedRoom.IsDeleted {
		return nil, util.WrapGQLError(ctx, errInvalidOfficeID, http.StatusNotFound, util.ErrorFlagNotFound)
	}

	// update room in the database
	room, err := svc.repoRegistry.Room().UpdateRoom(ctx, input)
	if err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	return &ent.RoomResponse{
		Message: roomUpdatedSuccess,
		Data:    room,
	}, nil
}

func (svc *roomSvcImpl) DeleteRoom(ctx context.Context, id uuid.UUID) (string, error) {
	// Validate id
	if id == uuid.Nil {
		return "", errors.New(errInvalidRoomID)
	}

	room, err := svc.repoRegistry.Room().GetRoom(ctx, id)
	if err != nil {
		return "", util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	if room == nil {
		return "", util.WrapGQLError(ctx, "room with id "+id.String()+" not found", http.StatusNotFound, util.ErrorFlagNotFound)
	}

	err = svc.repoRegistry.Room().DeleteRoom(ctx, id)
	if err != nil {
		return "", util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}
	return roomDeletedSuccess, nil
}

func (svc *roomSvcImpl) GetAvailableRooms(ctx context.Context, input ent.GetAvailableRoomInput) ([]*ent.AvailableRoomResponse, error) {
	var results []*ent.AvailableRoomResponse

	// Retrieve all rooms in the specified office
	roomFilter := ent.RoomFilter{
		OfficeID: input.OfficeID,
	}
	rooms, err := svc.repoRegistry.Room().GetRoomsByOfficeId(ctx, roomFilter)
	if err != nil {
		return nil, err
	}

	// Validate input
	if err := svc.validateGetAvailableRoomInput(input); err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusBadRequest, util.ErrorFlagValidateFail)
	}

	// Parse input times
	startDateTime, endDateTime, err := svc.parseInputTimes(input)
	if err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusBadRequest, util.ErrorFlagValidateFail)
	}

	// Find rooms with no overlapping bookings
	for _, room := range rooms {
		// Check for any overlapping bookings for this room
		conflictingBookings, err := svc.repoRegistry.Booking().GetConflictingBookings(ctx, room.ID, startDateTime, endDateTime)
		if err != nil {
			return nil, err
		}

		if !room.IsDeleted {
			roomStatus := false
			if len(conflictingBookings) == 0 {
				roomStatus = true
			}
			results = append(results, &ent.AvailableRoomResponse{
				ID:          room.ID.String(),
				Name:        room.Name,
				Color:       room.Color,
				Floor:       room.Floor,
				OfficeID:    room.OfficeID.String(),
				Description: &room.Description,
				ImageURL:    &room.ImageURL,
				Status:      roomStatus,
			})
		}
	}
	return results, nil
}

// parseInputTimes parses the input start date, start time, and end time into time.Time objects.
func (svc *roomSvcImpl) parseInputTimes(input ent.GetAvailableRoomInput) (time.Time, time.Time, error) {
	// Parse start date
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New(errParseDate)
	}

	// Parse start time
	if input.StartTime == nil {
		// Set default start time to the beginning of the business day (9:00 AM)
		defaultStartTime := startBusinessHour
		input.StartTime = &defaultStartTime
	}
	startTime, err := time.Parse("15:04", *input.StartTime)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New(errParseTime)
	}

	// Parse end time
	if input.EndTime == nil {
		// Set default end time to the end of the business day (5:30 PM)
		defaultEndTime := endBusinessHour
		input.EndTime = &defaultEndTime
	}
	endTime, err := time.Parse("15:04", *input.EndTime)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New(errParseTime)
	}

	var startDateTime, endDateTime time.Time

	// If the input is repeated, parse end date and set end time accordingly
	if input.IsRepeat != nil && *input.IsRepeat {
		if input.EndDate == nil {
			return time.Time{}, time.Time{}, errors.New(errEndDateRequired)
		}
		endDate, err := time.Parse("2006-01-02", *input.EndDate)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New(errParseDate)
		}
		startDateTime = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
		endDateTime = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)
	} else {
		// Set start and end times to the same day
		startDateTime = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
		endDateTime = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)
	}

	return startDateTime, endDateTime, nil
}

// Validation functions
func (svc *roomSvcImpl) validateCreateRoomInput(input ent.CreateRoomInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return errors.New(errRoomNameRequired)
	}
	if strings.TrimSpace(input.Color) == "" {
		return errors.New(errRoomColorRequired)
	}
	if input.OfficeID == "" {
		return errors.New(errOfficeIDRequired)
	}
	return nil
}

func (svc *roomSvcImpl) validateUpdateRoomInput(input ent.UpdateRoomInput) error {
	if input.ID == uuid.Nil.String() {
		return errors.New(errInvalidRoomID)
	}
	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return errors.New(errInvalidRoomName)
	}
	if input.Color != nil && strings.TrimSpace(*input.Color) == "" {
		return errors.New(errInvalidRoomColor)
	}
	return nil
}

func (svc *roomSvcImpl) validateGetAvailableRoomInput(input ent.GetAvailableRoomInput) error {
	if input.StartDate == "" {
		return errors.New(errStartDateRequired)
	}
	if input.OfficeID == "" {
		return errors.New(errOfficeIDRequired)
	}
	return nil
}

func (svc *roomSvcImpl) validateOfficeID(ctx context.Context, officeId string) error {
	offices, err := svc.repoRegistry.Office().GetOffices(ctx)
	if err != nil {
		return err
	}

	for _, office := range offices {
		if office.ID.String() == officeId {
			return nil
		}
	}

	return errors.New(errInvalidOfficeID)
}

func (svc *roomSvcImpl) applyRoomFilter(roomQuery *ent.RoomQuery, filter *ent.RoomFilter) error {
	if filter != nil {
		if filter.OfficeID != "" {
			roomQuery.Where(room.OfficeIDEQ(uuid.MustParse(filter.OfficeID)))
			return nil
		} else {
			return errors.New(errOfficeIDRequired)
		}
	}
	return nil
}

func (svc *roomSvcImpl) applySearchTermFilter(query *ent.RoomQuery, filter *ent.RoomFilter) {
	if filter.SearchTerm != nil && *filter.SearchTerm != "" {
		query.Where(
			room.Or(
				room.NameContainsFold(*filter.SearchTerm),
				room.FloorContainsFold(*filter.SearchTerm),
			),
		)
	}
}
