package dto

import "time"

type Response struct {
	ID               uint      `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Location         string    `json:"location"`
	StartsAt         time.Time `json:"starts_at"`
	TotalTickets     int       `json:"total_tickets"`
	AvailableTickets int       `json:"available_tickets"`
	Price            int       `json:"price"`
	CreatedAt        string    `json:"created_at"`
}
