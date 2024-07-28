package repository

import (
	"context"
	"fmt"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"

	"github.com/pkg/errors"
)

type Repository interface {
	User() UserRepository
	Role() RoleRepository
	Office() OfficeRepository
	Room() RoomRepository
	Booking() BookingRepository

	// DoInTx executes the given function in a transaction.
	DoInTx(ctx context.Context, txFunc func(ctx context.Context, repoRegistry Repository) error) error
}

type RepoImpl struct {
	entClient *ent.Client
	entTx     *ent.Tx
	user      UserRepository
	role      RoleRepository
	office    OfficeRepository
	room      RoomRepository
	booking   BookingRepository
}

func NewRepository(entClient *ent.Client) Repository {
	return &RepoImpl{
		entClient: entClient,
		user:      NewUserRepository(entClient),
		role:      NewRoleRepository(entClient),
		office:    NewOfficeRepository(entClient),
		room:      NewRoomRepository(entClient),
		booking:   NewBookingRepository(entClient),
	}
}

func (r *RepoImpl) User() UserRepository {
	return r.user
}

func (r *RepoImpl) Role() RoleRepository {
	return r.role
}

func (r *RepoImpl) Office() OfficeRepository {
	return r.office
}

func (r *RepoImpl) Room() RoomRepository {
	return r.room
}

func (r *RepoImpl) Booking() BookingRepository {
	return r.booking
}

func (r *RepoImpl) DoInTx(ctx context.Context, txFunc func(ctx context.Context, repoRegistry Repository) error) error {
	if r.entTx != nil {
		return errors.WithStack(errors.New("invalid tx state, no nested tx allowed"))
	}

	tx, err := r.entClient.Tx(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	committed := false

	defer func() {
		if committed {
			return
		}
		// rollback if not commited
		_ = tx.Rollback()
	}()

	impl := &RepoImpl{
		entTx:   tx,
		user:    NewUserRepository(tx.Client()),
		role:    NewRoleRepository(tx.Client()),
		office:  NewOfficeRepository(tx.Client()),
		room:    NewRoomRepository(tx.Client()),
		booking: NewBookingRepository(tx.Client()),
	}

	if err := txFunc(ctx, impl); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.WithStack(fmt.Errorf("failed to commit tx: %s", err.Error()))
	}

	committed = true
	return nil
}
