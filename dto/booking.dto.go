package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookingFilterDTO struct {
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	OfficeID  *uuid.UUID `json:"officeId,omitempty"`
	RoomID    *uuid.UUID `json:"roomId,omitempty"`
	Keyword   *string    `json:"keyword,omitempty"`
}

type CreateBookingInputDTO struct {
	Title     string    `json:"title"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	IsRepeat  *bool     `json:"isRepeat,omitempty"`
	OfficeID  uuid.UUID `json:"officeId"`
	RoomID    uuid.UUID `json:"roomId"`
}

type UpdateBookingInputDTO struct {
	ID        uuid.UUID  `json:"id"`
	Title     *string    `json:"title,omitempty"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	IsRepeat  *bool      `json:"isRepeat,omitempty"`
	OfficeID  *uuid.UUID `json:"officeId,omitempty"`
	RoomID    *uuid.UUID `json:"roomId,omitempty"`
}
