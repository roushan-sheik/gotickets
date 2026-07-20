package event

import (
	"gotickets/internal/domain/event/dto"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title            string    `json:"title" gorm:"type:varchar(150);not null"`
	Description      string    `json:"description" gorm:"type:text"`
	Location         string    `json:"location" gorm:"type:varchar(150);not null"`
	StartsAt         time.Time `json:"starts_at" gorm:"not null"`
	TotalTickets     int       `json:"total_tickets" gorm:"not null"`
	AvailableTickets int       `json:"available_tickets" gorm:"not null"`
	Price            int       `json:"price" gorm:"not null"`
}

func (e *Event) ToResponse() *dto.Response {
	return &dto.Response{
		ID:               e.ID,
		Title:            e.Title,
		Description:      e.Description,
		Location:         e.Location,
		StartsAt:         e.StartsAt,
		TotalTickets:     e.TotalTickets,
		AvailableTickets: e.AvailableTickets,
		Price:            e.Price,
		CreatedAt:        e.CreatedAt.String(),
	}
}
