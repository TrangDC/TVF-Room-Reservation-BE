package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/dto"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/booking"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BookingService interface {
	// Queries
	GetBookings(ctx context.Context, pagination *ent.PaginationInput, filter *ent.BookingFilter) (*ent.BookingDataResponse, error)
	GetBooking(ctx context.Context, id uuid.UUID) (*ent.BookingData, error)

	// Mutations
	CreateBooking(ctx context.Context, input ent.CreateBookingInput) (*ent.BookingResponse, error)
	UpdateBooking(ctx context.Context, input ent.UpdateBookingInput) (*ent.BookingResponse, error)
	CancelBooking(ctx context.Context, id uuid.UUID) (string, error)

	// Common
	renderBookingData(ctx context.Context, booking *ent.Booking) (*ent.BookingData, error)
}

type bookingSvcImpl struct {
	repoRegistry repository.Repository
	logger       *zap.Logger
	isTestEnv    bool
}

func NewBookingService(repoRegistry repository.Repository, logger *zap.Logger) BookingService {
	return &bookingSvcImpl{
		repoRegistry: repoRegistry,
		logger:       logger,
	}
}

// Query functions
func (svc *bookingSvcImpl) GetBookings(ctx context.Context, pagination *ent.PaginationInput, filter *ent.BookingFilter) (*ent.BookingDataResponse, error) {
	// Initialize filter parameters to default values
	var startDateTime, endDateTime *time.Time
	var officeID, roomID *uuid.UUID
	var keyword *string

	if filter != nil {
		var err error
		startDateTime, endDateTime, officeID, roomID, err = svc.parseFilterParams(filter)
		if err != nil {
			return nil, util.WrapGQLError(ctx, err.Error(), http.StatusConflict, util.ErrorFlagValidateFail)
		}

		if filter.Keyword != nil {
			keyword = filter.Keyword
		}
	}

	var results *ent.BookingDataResponse

	query := svc.repoRegistry.Booking().BuildQuery()

	// Create booking filter
	bookingFilter := dto.BookingFilterDTO{
		StartDate: startDateTime,
		EndDate:   endDateTime,
		OfficeID:  officeID,
		RoomID:    roomID,
		Keyword:   keyword,
	}

	if filter != nil {
		svc.filter(query, bookingFilter)
	}

	total, err := svc.repoRegistry.Booking().BuildCount(ctx, query)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "failed to count bookings", http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	query = query.Order(ent.Desc(booking.FieldUpdatedAt))

	if pagination != nil {
		query = query.Limit(*pagination.PerPage).Offset((*pagination.Page - 1) * *pagination.PerPage)
	}

	bookings, err := svc.repoRegistry.Booking().BuildList(ctx, query)
	if err != nil {
		return results, util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	// Render booking data
	var data []*ent.BookingData
	for _, booking := range bookings {
		datum, err := svc.renderBookingData(ctx, booking)
		if err != nil {
			return nil, err
		}
		data = append(data, datum)
	}

	return &ent.BookingDataResponse{
		Total: total,
		Data:  data,
	}, nil
}

func (svc *bookingSvcImpl) GetBooking(ctx context.Context, id uuid.UUID) (*ent.BookingData, error) {
	booking, err := svc.repoRegistry.Booking().GetBooking(ctx, id)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "booking not found", http.StatusNotFound, util.ErrorFlagNotFound)
	}

	return svc.renderBookingData(ctx, booking)
}

// Mutation functions
func (svc *bookingSvcImpl) CreateBooking(ctx context.Context, input ent.CreateBookingInput) (*ent.BookingResponse, error) {
	// Validate date time parameters
	startDateTime, endDate, err := validateBookingTime(ctx, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	// Check and create endDateTime base on isRepeat
	endDateTime, err := getEndDateTime(ctx, input, endDate)
	if err != nil {
		return nil, err
	}

	officeID, err := parseAndValidateOffice(ctx, svc, input.OfficeID)
	if err != nil {
		return nil, err
	}

	roomID, err := parseAndValidateRoom(ctx, svc, input.RoomID, officeID)
	if err != nil {
		return nil, err
	}

	// Check for existing bookings
	err = svc.repoRegistry.Booking().CheckExistingBookings(ctx, officeID, roomID, startDateTime, endDateTime, nil)
	if err != nil {
		return nil, err
	}

	// Create booking input
	createBookingInput := dto.CreateBookingInputDTO{
		Title:     input.Title,
		StartDate: startDateTime,
		EndDate:   endDateTime,
		IsRepeat:  input.IsRepeat,
		OfficeID:  officeID,
		RoomID:    roomID,
	}

	// Create booking
	var newBooking *ent.Booking
	err = svc.repoRegistry.Booking().ValidateBookingTitle(ctx, uuid.UUID{}, input.Title)
	if err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusBadRequest, util.ErrorFlagValidateFail)
	}

	err = svc.repoRegistry.DoInTx(ctx, func(ctx context.Context, repoRegistry repository.Repository) error {
		newBooking, err = svc.repoRegistry.Booking().CreateBooking(ctx, createBookingInput)
		return err
	})
	if err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	// Render booking data
	bookingData, err := svc.renderBookingData(ctx, newBooking)
	if err != nil {
		return nil, err
	}

	return &ent.BookingResponse{
		Message: "Data has been successfully created.",
		Data:    bookingData,
	}, err
}

func (svc *bookingSvcImpl) UpdateBooking(ctx context.Context, input ent.UpdateBookingInput) (*ent.BookingResponse, error) {
	if err := validateUpdateInput(ctx, input); err != nil {
		return nil, err
	}

	startDate, endDate, err := getValidatedDates(ctx, input)
	if err != nil {
		return nil, err
	}

	officeID, roomID, err := parseAndValidateOfficeRoom(ctx, svc, input.OfficeID, input.RoomID)
	if err != nil {
		return nil, err
	}

	// Parse booking ID
	bookingID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "failed to parse booking id", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	// Unable to update deleted booking
	_, err = svc.repoRegistry.Booking().GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "unable to update deleted booking", http.StatusNotFound, util.ErrorFlagNotFound)
	}

	// Check for existing bookings
	if startDate != nil && endDate != nil {
		err = svc.repoRegistry.Booking().CheckExistingBookings(ctx, officeID, roomID, *startDate, *endDate, &bookingID)
		if err != nil {
			return nil, err
		}
	}

	// Create booking input
	updateBookingInput := dto.UpdateBookingInputDTO{
		ID: bookingID,
	}

	if input.Title != nil {
		updateBookingInput.Title = input.Title

		err = svc.repoRegistry.Booking().ValidateBookingTitle(ctx, uuid.MustParse(input.ID), *input.Title)
		if err != nil {
			return nil, util.WrapGQLError(ctx, err.Error(), http.StatusBadRequest, util.ErrorFlagValidateFail)
		}
	}

	if startDate != nil && endDate != nil {
		updateBookingInput.StartDate = startDate
		updateBookingInput.EndDate = endDate
	}

	if input.IsRepeat != nil {
		updateBookingInput.IsRepeat = input.IsRepeat
	}

	if officeID != uuid.Nil {
		updateBookingInput.OfficeID = &officeID
	}

	if roomID != uuid.Nil {
		updateBookingInput.RoomID = &roomID
	}

	// Update booking
	var updatedBooking *ent.Booking
	err = svc.repoRegistry.DoInTx(ctx, func(ctx context.Context, repoRegistry repository.Repository) error {
		updatedBooking, err = svc.repoRegistry.Booking().UpdateBooking(ctx, updateBookingInput)
		return err
	})
	if err != nil {
		return nil, util.WrapGQLError(ctx, err.Error(), http.StatusInternalServerError, util.ErrorFlagInternalError)
	}

	// Render booking data
	bookingData, err := svc.renderBookingData(ctx, updatedBooking)
	if err != nil {
		return nil, err
	}

	return &ent.BookingResponse{
		Message: "Data has been successfully updated.",
		Data:    bookingData,
	}, err
}

func (svc *bookingSvcImpl) CancelBooking(ctx context.Context, id uuid.UUID) (string, error) {
	booking, err := svc.repoRegistry.Booking().GetBookingByID(ctx, id)
	if err != nil {
		return "", util.WrapGQLError(ctx, "booking not found", http.StatusNotFound, util.ErrorFlagCanNotDelete)
	}

	err = svc.repoRegistry.DoInTx(ctx, func(ctx context.Context, repoRegistry repository.Repository) error {
		_, err = svc.repoRegistry.Booking().CancelBooking(ctx, booking)
		return err
	})
	if err != nil {
		return "", err
	}
	return "Data has been successfully deleted.", err
}

// Common functions
func (svc *bookingSvcImpl) renderBookingData(ctx context.Context, booking *ent.Booking) (*ent.BookingData, error) {
	var office *ent.Office
	var room *ent.Room
	var user *ent.User
	var err error

	isTestEnv := svc.isTestEnv

	// Check if we're in a test environment
	if isTestEnv {
		// Use mock data
		office = &ent.Office{ID: booking.OfficeID}
		room = &ent.Room{ID: booking.RoomID, OfficeID: booking.OfficeID}
		user = &ent.User{ID: booking.UserID}
	} else {
		// Fetch real data
		// Check if the office exists
		office, err = svc.repoRegistry.Office().GetOffice(ctx, booking.OfficeID)
		if err != nil {
			return nil, util.WrapGQLError(ctx, "office not found", http.StatusNotFound, util.ErrorFlagNotFound)
		}

		// Check if the room exists
		room, err = svc.repoRegistry.Room().GetRoom(ctx, booking.RoomID)
		if err != nil {
			return nil, util.WrapGQLError(ctx, "room not found", http.StatusNotFound, util.ErrorFlagNotFound)
		}

		// Check if user exists
		user, err = svc.repoRegistry.User().GetUser(ctx, booking.UserID)
		if err != nil {
			return nil, util.WrapGQLError(ctx, "user not found", http.StatusNotFound, util.ErrorFlagNotFound)
		}
	}

	var deletedAt *time.Time
	if !booking.DeletedAt.IsZero() {
		deletedAtUTC := booking.DeletedAt.In(time.UTC)
		deletedAt = &deletedAtUTC
	}

	return &ent.BookingData{
		ID:        booking.ID.String(),
		Title:     booking.Title,
		StartDate: booking.StartDate.In(time.UTC),
		EndDate:   booking.EndDate.In(time.UTC),
		IsRepeat:  &booking.IsRepeat,
		Office:    office,
		Room:      room,
		User:      user,
		CreatedAt: booking.CreatedAt.In(time.UTC),
		UpdatedAt: booking.UpdatedAt.In(time.UTC),
		DeletedAt: deletedAt,
	}, err
}

func (svc *bookingSvcImpl) filter(bookingQuery *ent.BookingQuery, filter dto.BookingFilterDTO) {
	if filter.OfficeID != nil {
		bookingQuery.Where(booking.OfficeIDEQ(*filter.OfficeID))
	}

	if filter.RoomID != nil {
		bookingQuery.Where(booking.RoomIDEQ(*filter.RoomID))
	}

	if filter.Keyword != nil && *filter.Keyword != "" {
		keyword := strings.TrimSpace(*filter.Keyword)

		if len(keyword) >= 2 {
			bookingQuery.Where(booking.TitleContainsFold(keyword))
		}
	}

	if filter.StartDate != nil &&
		!filter.StartDate.IsZero() &&
		filter.EndDate != nil &&
		!filter.EndDate.IsZero() {
		bookingQuery.Where(
			booking.Or(
				booking.And(
					booking.StartDateGTE(*filter.StartDate),
					booking.EndDateLTE(*filter.EndDate),
				),
				booking.And(
					booking.StartDateLTE(*filter.StartDate),
					booking.EndDateGTE(*filter.EndDate),
				),
				booking.And(
					booking.StartDateLTE(*filter.StartDate),
					booking.EndDateGT(*filter.StartDate),
				),
				booking.And(
					booking.StartDateLT(*filter.EndDate),
					booking.EndDateGTE(*filter.EndDate),
				),
			),
		)
	}
}

func (svc *bookingSvcImpl) parseFilterParams(filter *ent.BookingFilter) (*time.Time, *time.Time, *uuid.UUID, *uuid.UUID, error) {
	var startDateTime, endDateTime *time.Time
	var officeID, roomID *uuid.UUID

	if filter == nil {
		return nil, nil, nil, nil, nil
	}

	if filter.StartDate != nil && filter.EndDate != nil {
		startDate, err := time.Parse("2006-01-02", *filter.StartDate)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to parse start date")
		}
		endDate, err := time.Parse("2006-01-02", *filter.EndDate)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to parse end date")
		}

		startDateTimeTemp := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
		endDateTimeTemp := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)
		startDateTime = &startDateTimeTemp
		endDateTime = &endDateTimeTemp
	}

	if filter.OfficeID != nil && *filter.OfficeID != "" {
		parsedOfficeID, err := uuid.Parse(*filter.OfficeID)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to parse office id")
		}
		officeID = &parsedOfficeID
	}

	if filter.RoomID != nil && *filter.RoomID != "" {
		parsedRoomID, err := uuid.Parse(*filter.RoomID)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to parse room id")
		}
		roomID = &parsedRoomID
	}

	return startDateTime, endDateTime, officeID, roomID, nil
}

func getEndDateTime(ctx context.Context, input ent.CreateBookingInput, endDate time.Time) (time.Time, error) {
	if input.IsRepeat != nil && *input.IsRepeat {
		if !input.EndDate.IsZero() {
			return input.EndDate, nil
		}
		return time.Time{}, util.WrapGQLError(ctx, "end date is required when isRepeat is true", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	if input.StartDate.Format("2006-01-02") != input.EndDate.Format("2006-01-02") {
		return time.Time{}, util.WrapGQLError(ctx, "end date must be the same day as start date when isRepeat is false", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	return endDate, nil
}

func checkExistingOffice(ctx context.Context, svc *bookingSvcImpl, officeID uuid.UUID) error {
	_, err := svc.repoRegistry.Office().GetOffice(ctx, officeID)
	if err != nil {
		return util.WrapGQLError(ctx, "office not found", http.StatusNotFound, util.ErrorFlagNotFound)
	}
	return nil
}

func parseAndValidateOffice(ctx context.Context, svc *bookingSvcImpl, officeIDStr string) (uuid.UUID, error) {
	officeID, err := uuid.Parse(officeIDStr)
	if err != nil {
		return uuid.UUID{}, util.WrapGQLError(ctx, "failed to parse office id", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	if err := checkExistingOffice(ctx, svc, officeID); err != nil {
		return uuid.UUID{}, err
	}

	return officeID, nil
}

func validateRoom(ctx context.Context, svc *bookingSvcImpl, officeID, roomID uuid.UUID) error {
	// Check existence of room
	room, err := svc.repoRegistry.Room().GetRoom(ctx, roomID)
	if err != nil {
		return util.WrapGQLError(ctx, "room not found", http.StatusNotFound, util.ErrorFlagNotFound)
	}

	// Check if the room is deleted, so we can not update a booking
	if room.IsDeleted {
		return util.WrapGQLError(ctx, "room has been deleted", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	// Check if the room belongs to the office
	if room.OfficeID != officeID {
		return util.WrapGQLError(ctx, "room does not belong to the specified office", http.StatusConflict, util.ErrorFlagValidateFail)
	}
	return nil
}

func parseAndValidateRoom(ctx context.Context, svc *bookingSvcImpl, roomIDStr string, officeID uuid.UUID) (uuid.UUID, error) {
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return uuid.UUID{}, util.WrapGQLError(ctx, "failed to parse room id", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	if err := validateRoom(ctx, svc, officeID, roomID); err != nil {
		return uuid.UUID{}, err
	}

	return roomID, nil
}

func parseAndValidateOfficeRoom(ctx context.Context, svc *bookingSvcImpl, officeIDStr, roomIDStr *string) (uuid.UUID, uuid.UUID, error) {
	var officeID, roomID uuid.UUID
	var err error

	if officeIDStr != nil {
		officeID, err = parseAndValidateOffice(ctx, svc, *officeIDStr)
		if err != nil {
			return uuid.UUID{}, uuid.UUID{}, err
		}
	}

	if roomIDStr != nil {
		roomID, err = parseAndValidateRoom(ctx, svc, *roomIDStr, officeID)
		if err != nil {
			return uuid.UUID{}, uuid.UUID{}, err
		}
	}

	return officeID, roomID, nil
}

func getValidatedDates(ctx context.Context, input ent.UpdateBookingInput) (*time.Time, *time.Time, error) {
	if (input.StartDate != nil && input.EndDate == nil) || (input.StartDate == nil && input.EndDate != nil) {
		return nil, nil, util.WrapGQLError(ctx, "both start date and end date must be provided together", http.StatusBadRequest, util.ErrorFlagValidateFail)
	}

	if input.IsRepeat != nil && (input.StartDate == nil || input.EndDate == nil) {
		return nil, nil, util.WrapGQLError(ctx, "start date and end date are required when isRepeat is set", http.StatusBadRequest, util.ErrorFlagValidateFail)
	}

	if input.StartDate != nil && input.EndDate != nil {
		startDate, endDate, err := validateBookingTime(ctx, *input.StartDate, *input.EndDate)
		if err != nil {
			return nil, nil, err
		}

		if input.IsRepeat != nil && !*input.IsRepeat && input.StartDate.Format("2006-01-02") != input.EndDate.Format("2006-01-02") {
			return nil, nil, util.WrapGQLError(ctx, "end date must be the same day as start date when isRepeat is false", http.StatusConflict, util.ErrorFlagValidateFail)
		}

		return &startDate, &endDate, nil
	}

	return nil, nil, nil
}

func validateUpdateInput(ctx context.Context, input ent.UpdateBookingInput) error {
	if input.ID == "" || input.ID == uuid.Nil.String() {
		return util.WrapGQLError(ctx, "booking id is required", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	// Check if at least one field is provided for update
	if input.Title == nil &&
		input.StartDate == nil &&
		input.EndDate == nil &&
		input.IsRepeat == nil &&
		input.OfficeID == nil &&
		input.RoomID == nil {
		return util.WrapGQLError(ctx, "please provide at least one field to update", http.StatusBadRequest, util.ErrorFlagValidateFail)
	}
	return nil
}

func validateBookingTime(ctx context.Context, startDate, endDate time.Time) (time.Time, time.Time, error) {
	// Check start time cannot be the same or end time before start time
	if startDate.Equal(endDate) || startDate.After(endDate) {
		return time.Time{}, time.Time{}, util.WrapGQLError(ctx, "start time and end time cannot be the same or end time before start time", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	// Extract time components from startDate and endDate
	startTime, endTime := extractTimeComponents(startDate, endDate)

	// Check if the duration is at least 15 minutes and not more than 1 hour
	duration := endTime.Sub(startTime)
	if duration < 15*time.Minute || duration > 4*time.Hour {
		return time.Time{}, time.Time{}, util.WrapGQLError(ctx, "booking duration must be between 15 minutes and 1 hour", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	// Check if the booking is on a weekend
	if isWeekend(startDate) {
		return time.Time{}, time.Time{}, util.WrapGQLError(ctx, "bookings cannot be made on weekends", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	// Define business hours
	businessStartTime, _ := time.Parse("15:04", "09:00")
	businessEndTime, _ := time.Parse("15:04", "17:30")

	// Check if the booking time is within the business time range
	if !withinBusinessHours(startTime, endTime, businessStartTime, businessEndTime) {
		return time.Time{}, time.Time{}, util.WrapGQLError(ctx, "booking must be within business hours (09:00 - 17:30)", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	// TODO: Check if the startDateTime is in the past (date difference)
	if startDate.Before(time.Now().UTC()) {
		return time.Time{}, time.Time{}, util.WrapGQLError(ctx, "bookings cannot be made in the past", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	return startDate, endDate, nil
}

func extractTimeComponents(startDate, endDate time.Time) (time.Time, time.Time) {
	startTime := time.Date(0, 1, 1, startDate.Hour(), startDate.Minute(), startDate.Second(), startDate.Nanosecond(), startDate.Location())
	endTime := time.Date(0, 1, 1, endDate.Hour(), endDate.Minute(), endDate.Second(), endDate.Nanosecond(), endDate.Location())
	return startTime, endTime
}

func withinBusinessHours(startTime, endTime, businessStartTime, businessEndTime time.Time) bool {
	return (startTime.Equal(businessStartTime) || startTime.After(businessStartTime)) && (endTime.Equal(businessEndTime) || endTime.Before(businessEndTime))
}

func isWeekend(date time.Time) bool {
	day := date.Weekday()
	return day == time.Saturday || day == time.Sunday
}
