package repository

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/dto"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/booking"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"
)

type BookingRepository interface {
	// Queries
	BuildQuery() *ent.BookingQuery
	BuildGet(ctx context.Context, query *ent.BookingQuery) (*ent.Booking, error)
	BuildList(ctx context.Context, query *ent.BookingQuery) ([]*ent.Booking, error)
	BuildCount(ctx context.Context, query *ent.BookingQuery) (int, error)

	GetBooking(ctx context.Context, id uuid.UUID) (*ent.Booking, error)

	// Mutations
	BuildUpdateOne(ctx context.Context, model *ent.Booking) *ent.BookingUpdateOne
	BuildSaveUpdateOne(ctx context.Context, update *ent.BookingUpdateOne) (*ent.Booking, error)

	CreateBooking(ctx context.Context, input dto.CreateBookingInputDTO) (*ent.Booking, error)
	UpdateBooking(ctx context.Context, input dto.UpdateBookingInputDTO) (*ent.Booking, error)
	CancelBooking(ctx context.Context, booking *ent.Booking) (*ent.Booking, error)

	// Common
	GetBookingByID(ctx context.Context, id uuid.UUID) (*ent.Booking, error)
	CheckExistingBookings(ctx context.Context, officeID, roomID uuid.UUID, startDateTime, endDateTime time.Time, excludeBookingID *uuid.UUID) error
	GetConflictingBookings(ctx context.Context, roomID uuid.UUID, startDateTime, endDateTime time.Time) ([]*ent.Booking, error)
	ValidateBookingTitle(ctx context.Context, newsId uuid.UUID, title string) error
}

type bookingRepoImpl struct {
	client *ent.Client
}

func NewBookingRepository(client *ent.Client) BookingRepository {
	return &bookingRepoImpl{
		client: client,
	}
}

// Base functions
func (rps *bookingRepoImpl) BuildQuery() *ent.BookingQuery {
	return rps.client.Booking.Query().Where(booking.DeletedAtIsNil()).WithOffice().WithRoom().WithUser()
}

func (rps *bookingRepoImpl) BuildGet(ctx context.Context, query *ent.BookingQuery) (*ent.Booking, error) {
	return query.First(ctx)
}

func (rps *bookingRepoImpl) BuildList(ctx context.Context, query *ent.BookingQuery) ([]*ent.Booking, error) {
	return query.All(ctx)
}

func (rps *bookingRepoImpl) BuildCount(ctx context.Context, query *ent.BookingQuery) (int, error) {
	return query.Count(ctx)
}

func (rps *bookingRepoImpl) BuildExist(ctx context.Context, query *ent.BookingQuery) (bool, error) {
	return query.Exist(ctx)
}

func (rps *bookingRepoImpl) BuildCreate() *ent.BookingCreate {
	return rps.client.Booking.Create().SetUpdatedAt(time.Now())
}

func (rps *bookingRepoImpl) BuildUpdateOne(ctx context.Context, model *ent.Booking) *ent.BookingUpdateOne {
	return model.Update().SetUpdatedAt(time.Now())
}

func (rps *bookingRepoImpl) BuildSaveUpdateOne(ctx context.Context, update *ent.BookingUpdateOne) (*ent.Booking, error) {
	return update.Save(ctx)
}

// Query functions
func (rps *bookingRepoImpl) GetBooking(ctx context.Context, id uuid.UUID) (*ent.Booking, error) {
	query := rps.BuildQuery().Where(booking.IDEQ(id))
	return rps.BuildGet(ctx, query)
}

// Mutation functions
func (rps *bookingRepoImpl) CreateBooking(ctx context.Context, input dto.CreateBookingInputDTO) (*ent.Booking, error) {
	userID := ctx.Value("user_id").(uuid.UUID)

	// Set default value for isRepeat if it's nil
	isRepeat := false
	if input.IsRepeat != nil {
		isRepeat = *input.IsRepeat
	}

	// Create a new booking record
	newBooking, err := rps.BuildCreate().
		SetTitle(input.Title).
		SetSlug(util.SlugGeneration(input.Title)).
		SetStartDate(input.StartDate).
		SetEndDate(input.EndDate).
		SetIsRepeat(isRepeat).
		SetOfficeID(input.OfficeID).
		SetRoomID(input.RoomID).
		SetUserID(userID).
		Save(ctx)

	if err != nil {
		return nil, util.WrapGQLError(ctx, "failed to create booking", http.StatusConflict, util.ErrorFlagCanNotCreate)
	}

	return newBooking, err
}

func (rps *bookingRepoImpl) UpdateBooking(ctx context.Context, input dto.UpdateBookingInputDTO) (*ent.Booking, error) {
	bookingToUpdate, err := rps.client.Booking.Get(ctx, input.ID)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "failed to get booking", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	updatedBooking := rps.BuildUpdateOne(ctx, bookingToUpdate)

	// Update booking fields
	if input.Title != nil {
		updatedBooking.SetTitle(*input.Title)
		updatedBooking.SetSlug(util.SlugGeneration(*input.Title))
	}

	if input.StartDate != nil && input.EndDate != nil {
		updatedBooking.SetStartDate(*input.StartDate)
		updatedBooking.SetEndDate(*input.EndDate)
	}

	if input.IsRepeat != nil {
		updatedBooking.SetIsRepeat(*input.IsRepeat)
	}

	if input.OfficeID != nil {
		updatedBooking.SetOfficeID(*input.OfficeID)
	}

	if input.RoomID != nil {
		updatedBooking.SetRoomID(*input.RoomID)
	}

	// Save updated booking
	result, err := rps.BuildSaveUpdateOne(ctx, updatedBooking)
	if err != nil {
		return nil, util.WrapGQLError(ctx, "failed to update booking", http.StatusConflict, util.ErrorFlagCanNotUpdate)
	}

	return result, nil
}

func (rps *bookingRepoImpl) CancelBooking(ctx context.Context, booking *ent.Booking) (*ent.Booking, error) {
	update := rps.BuildUpdateOne(ctx, booking).SetDeletedAt(time.Now()).SetUpdatedAt(time.Now())
	return update.Save(ctx)
}

func (rps *bookingRepoImpl) GetBookingByID(ctx context.Context, id uuid.UUID) (*ent.Booking, error) {
	query := rps.BuildQuery().Where(booking.IDEQ(id))
	return rps.BuildGet(ctx, query)
}

// Check for existing bookings in the same office, room and time range, excluding the current booking
func (rps *bookingRepoImpl) CheckExistingBookings(ctx context.Context, officeID, roomID uuid.UUID, startDateTime, endDateTime time.Time, excludeBookingID *uuid.UUID) error {
	query := rps.BuildQuery().
		Where(
			booking.And(
				booking.OfficeID(officeID),
				booking.RoomID(roomID),
				booking.Or(
					booking.And(
						booking.StartDateLT(endDateTime),
						booking.EndDateGT(startDateTime),
					),
					booking.And(
						booking.StartDateGTE(startDateTime),
						booking.StartDateLT(endDateTime),
					),
					booking.And(
						booking.EndDateGT(startDateTime),
						booking.EndDateLTE(endDateTime),
					),
				),
			),
		)

	if excludeBookingID != nil {
		query = query.Where(booking.IDNEQ(*excludeBookingID))
	}

	existingBookings, err := query.All(ctx)
	if err != nil {
		return util.WrapGQLError(ctx, "failed to query existing bookings", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	if len(existingBookings) > 0 {
		return util.WrapGQLError(ctx, "booking already exists for the given time range", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	return nil
}

func (rps *bookingRepoImpl) GetConflictingBookings(ctx context.Context, roomID uuid.UUID, startDateTime, endDateTime time.Time) ([]*ent.Booking, error) {
	return rps.BuildQuery().
		Where(
			booking.RoomID(roomID),
			booking.Or(
				booking.And(
					booking.StartDateLTE(startDateTime),
					booking.EndDateGTE(endDateTime),
				),
				booking.And(
					booking.StartDateGTE(startDateTime),
					booking.StartDateLT(endDateTime),
				),
				booking.And(
					booking.EndDateGT(startDateTime),
					booking.EndDateLTE(endDateTime),
				),
			),
		).All(ctx)
}

func (rps *bookingRepoImpl) ValidateBookingTitle(ctx context.Context, bookingID uuid.UUID, title string) error {
	query := rps.BuildQuery().Where(booking.SlugEQ(util.SlugGeneration(title)))
	if bookingID != uuid.Nil {
		query = query.Where(booking.IDNEQ(bookingID))
	}
	isExist, err := rps.BuildExist(ctx, query)
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("booking with title '%s' already exists", title)
	}
	return nil
}
