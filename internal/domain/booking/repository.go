package booking

import (
	"errors"
	"gotickets/internal/domain/event"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrBookingNotFound         = errors.New("booking not found")
	ErrNotEnoughTickets        = errors.New("not enough tickets available")
	ErrBookingAlreadyCancelled = errors.New("booking already cancelled")
	ErrForbiddenBookingAccess  = errors.New("you do not own this booking")
)

type Repository interface {
	Create(booking *Booking) error
	GetByID(bookingId uint) (*Booking, error)
	GetByUserID(userId uint) ([]*Booking, error)
	Update(booking *Booking) error
	CreateWithTicketUpdate(userId uint, eventId uint, quantity int) (*Booking, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(booking *Booking) error {
	return r.db.Create(booking).Error
}

func (r *repository) GetByID(bookingId uint) (*Booking, error) {
	var booking Booking

	err := r.db.First(&booking, bookingId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookingNotFound
		}

		return nil, err
	}

	return &booking, nil
}

func (r *repository) GetByUserID(userId uint) ([]*Booking, error) {
	var bookings []*Booking

	err := r.db.Where("user_id = ?", userId).Find(&bookings).Error
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *repository) Update(booking *Booking) error {
	return r.db.Save(booking).Error
}

func (r *repository) CreateWithTicketUpdate(userId uint, eventId uint, quantity int) (*Booking, error) {
	var booking Booking

	// start transaction
	err := r.db.Transaction(func(tx *gorm.DB) error {

		var eventData event.Event

		err := tx.Clauses(clause.Locking{Strength: "Update"}).First(&eventData, eventId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return event.ErrEventNotFound
			}
			return err
		}

		if eventData.AvailableTickets < int(quantity) {
			return ErrNotEnoughTickets
		}

		booking = Booking{
			UserID:      userId,
			EventID:     eventData.ID,
			Quantity:    quantity,
			Status:      BookingConfirmed,
			TotalPrice:  quantity * eventData.Price,
			BookingCode: generateBookingCode(),
		}

		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		eventData.AvailableTickets = eventData.AvailableTickets - int(quantity)
		if err := tx.Save(&eventData).Error; err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return &booking, nil

}
