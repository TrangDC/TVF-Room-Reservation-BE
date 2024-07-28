package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/dto"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/repository"
	repoMock "gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/repository/mock"
	"go.uber.org/zap"
)

// Booking cancel testing function
func TestCancelBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	bookingID := "38ec1605-5c21-4c47-a0f4-e8b01a375e2d"
	parsedBookingID := uuid.MustParse(bookingID)

	mockRepo := repoMock.NewMockRepository(ctrl)
	mockBookingRepo := repoMock.NewMockBookingRepository(ctrl)

	// Set up expectations
	mockRepo.EXPECT().Booking().Return(mockBookingRepo).AnyTimes()
	mockBookingRepo.EXPECT().GetBookingByID(ctx, parsedBookingID).Return(&ent.Booking{ID: parsedBookingID}, nil)
	mockRepo.EXPECT().DoInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context, repo repository.Repository) error) error {
			return fn(ctx, mockRepo)
		},
	)

	mockBookingRepo.EXPECT().CancelBooking(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, booking *ent.Booking) (*ent.Booking, error) {
			booking.DeletedAt = time.Now()
			return booking, nil
		},
	)

	// Create service
	logger, _ := zap.NewProduction()
	svc := &bookingSvcImpl{repoRegistry: mockRepo, logger: logger}

	// Call cancel booking function
	message, err := svc.CancelBooking(ctx, parsedBookingID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "Data has been successfully deleted.", message)
}

// Booking create testing function
func TestCreateBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	startDate := time.Date(2024, 7, 15, 10, 0, 0, 0, time.UTC)
	endDate := startDate.Add(30 * time.Minute)
	input := ent.CreateBookingInput{
		Title:     "Meeting",
		StartDate: startDate,
		EndDate:   endDate,
		IsRepeat:  nil,
		OfficeID:  "f42ba67c-fbb6-4c21-a307-9f18357fe332",
		RoomID:    "31051962-10a8-4aa8-a291-e5dee5aa0070",
	}

	officeID, _ := uuid.Parse(input.OfficeID)
	roomID, _ := uuid.Parse(input.RoomID)

	mockRepo := repoMock.NewMockRepository(ctrl)
	mockBookingRepo := repoMock.NewMockBookingRepository(ctrl)
	mockOfficeRepo := repoMock.NewMockOfficeRepository(ctrl)
	mockRoomRepo := repoMock.NewMockRoomRepository(ctrl)

	// Set up expectations
	mockRepo.EXPECT().Booking().Return(mockBookingRepo).AnyTimes()
	mockRepo.EXPECT().Office().Return(mockOfficeRepo).AnyTimes()
	mockRepo.EXPECT().Room().Return(mockRoomRepo).AnyTimes()

	mockOfficeRepo.EXPECT().GetOffice(ctx, officeID).Return(&ent.Office{ID: officeID}, nil)
	mockRoomRepo.EXPECT().GetRoom(ctx, roomID).Return(&ent.Room{ID: roomID, OfficeID: officeID}, nil)
	mockBookingRepo.EXPECT().CheckExistingBookings(ctx, officeID, roomID, startDate, endDate, nil).Return(nil)
	mockBookingRepo.EXPECT().ValidateBookingTitle(ctx, uuid.Nil, input.Title).Return(nil)

	// Mock DoInTx to execute the function within the transaction
	mockRepo.EXPECT().DoInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context, repo repository.Repository) error) error {
			return fn(ctx, mockRepo)
		},
	)

	mockBookingRepo.EXPECT().CreateBooking(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, input dto.CreateBookingInputDTO) (*ent.Booking, error) {
			return &ent.Booking{
				ID:        uuid.New(),
				Title:     input.Title,
				StartDate: input.StartDate,
				EndDate:   input.EndDate,
				OfficeID:  input.OfficeID,
				RoomID:    input.RoomID,
			}, nil
		},
	)

	// Create service
	logger, _ := zap.NewProduction()
	svc := &bookingSvcImpl{repoRegistry: mockRepo, logger: logger, isTestEnv: true}

	// Call create booking function
	response, err := svc.CreateBooking(ctx, input)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Data has been successfully created.", response.Message)
	assert.NotNil(t, response.Data)
	assert.Equal(t, input.Title, response.Data.Title)
	assert.Equal(t, officeID, response.Data.Office.ID)
	assert.Equal(t, roomID, response.Data.Room.ID)
}

// Booking update testing function
func TestUpdateBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	bookingID := uuid.New()
	startDate := time.Date(2024, 7, 15, 10, 0, 0, 0, time.UTC)
	endDate := startDate.Add(30 * time.Minute)
	input := ent.UpdateBookingInput{
		ID:        bookingID.String(),
		Title:     stringPtr("Updated Meeting"),
		StartDate: &startDate,
		EndDate:   &endDate,
		IsRepeat:  boolPtr(false),
		OfficeID:  stringPtr("f42ba67c-fbb6-4c21-a307-9f18357fe332"),
		RoomID:    stringPtr("31051962-10a8-4aa8-a291-e5dee5aa0070"),
	}

	officeID, _ := uuid.Parse(*input.OfficeID)
	roomID, _ := uuid.Parse(*input.RoomID)

	mockRepo := repoMock.NewMockRepository(ctrl)
	mockBookingRepo := repoMock.NewMockBookingRepository(ctrl)
	mockOfficeRepo := repoMock.NewMockOfficeRepository(ctrl)
	mockRoomRepo := repoMock.NewMockRoomRepository(ctrl)

	// Set up expectations
	mockRepo.EXPECT().Booking().Return(mockBookingRepo).AnyTimes()
	mockRepo.EXPECT().Office().Return(mockOfficeRepo).AnyTimes()
	mockRepo.EXPECT().Room().Return(mockRoomRepo).AnyTimes()

	mockBookingRepo.EXPECT().GetBookingByID(ctx, bookingID).Return(&ent.Booking{ID: bookingID}, nil)
	mockOfficeRepo.EXPECT().GetOffice(ctx, officeID).Return(&ent.Office{ID: officeID}, nil)
	mockRoomRepo.EXPECT().GetRoom(ctx, roomID).Return(&ent.Room{ID: roomID, OfficeID: officeID}, nil)
	mockBookingRepo.EXPECT().CheckExistingBookings(ctx, officeID, roomID, startDate, endDate, &bookingID).Return(nil)
	mockBookingRepo.EXPECT().ValidateBookingTitle(ctx, bookingID, *input.Title).Return(nil)

	// Mock DoInTx to execute the function within the transaction
	mockRepo.EXPECT().DoInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context, repo repository.Repository) error) error {
			return fn(ctx, mockRepo)
		},
	)

	mockBookingRepo.EXPECT().UpdateBooking(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, input dto.UpdateBookingInputDTO) (*ent.Booking, error) {
			return &ent.Booking{
				ID:        input.ID,
				Title:     *input.Title,
				StartDate: *input.StartDate,
				EndDate:   *input.EndDate,
				IsRepeat:  *input.IsRepeat,
				OfficeID:  *input.OfficeID,
				RoomID:    *input.RoomID,
			}, nil
		},
	)

	// Create service
	logger, _ := zap.NewProduction()
	svc := &bookingSvcImpl{repoRegistry: mockRepo, logger: logger, isTestEnv: true}

	// Call update booking function
	response, err := svc.UpdateBooking(ctx, input)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Data has been successfully updated.", response.Message)
	assert.NotNil(t, response.Data)
	assert.Equal(t, *input.Title, response.Data.Title)
	assert.Equal(t, officeID, response.Data.Office.ID)
	assert.Equal(t, roomID, response.Data.Room.ID)
	assert.Equal(t, startDate.UTC(), response.Data.StartDate)
	assert.Equal(t, endDate.UTC(), response.Data.EndDate)
}

func TestUpdateBooking_DuplicateTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	bookingID := uuid.MustParse("5ed95ef8-7cdf-4ae7-b249-5e66f02df0d9")
	existingTitle := "123123 1"
	input := ent.UpdateBookingInput{
		ID:    bookingID.String(),
		Title: stringPtr(existingTitle),
	}

	mockRepo := repoMock.NewMockRepository(ctrl)
	mockBookingRepo := repoMock.NewMockBookingRepository(ctrl)

	mockRepo.EXPECT().Booking().Return(mockBookingRepo).AnyTimes()

	mockBookingRepo.EXPECT().GetBookingByID(ctx, bookingID).Return(&ent.Booking{ID: bookingID}, nil)

	// Expect ValidateBookingTitle to return an error (duplicate title)
	mockBookingRepo.EXPECT().ValidateBookingTitle(ctx, bookingID, *input.Title).Return(fmt.Errorf("Booking with title '%s' already exists", existingTitle))

	logger, _ := zap.NewProduction()
	svc := &bookingSvcImpl{repoRegistry: mockRepo, logger: logger, isTestEnv: true}

	// Call update booking function
	response, err := svc.UpdateBooking(ctx, input)

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), fmt.Sprintf("input: Booking with title '%s' already exists", existingTitle))
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
