package booking

import (
	"gotickets/internal/domain/booking/dto"
	"gotickets/internal/domain/event"

	"github.com/google/uuid"
)

type service struct {
	bookingRepo Repository
	eventRepo   event.Repository
}

func NewService(bookingRepo Repository, eventRepo event.Repository) *service {
	return &service{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

func generateBookingCode() string {
	return "GT-" + uuid.New().String()
}

func (s *service) CreateBooking(userId uint, req dto.CreateRequest) (*dto.Response, error) {
	booking, err := s.bookingRepo.CreateWithTicketUpdate(userId, req.EventID, req.Quantity)
	if err != nil {
		return nil, err
	}

	return booking.ToResponse(), nil
}

func (s *service) GetMyBookings(userId uint) ([]*dto.Response, error) {
	bookings, err := s.bookingRepo.GetByUserID(userId)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.Response, len(bookings)) // Initialize the slice with the correct length

	for i, b := range bookings {
		responses[i] = b.ToResponse()
	}

	return responses, nil
}
