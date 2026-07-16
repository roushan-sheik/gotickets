package dto

import "time"

type Response struct {
	ID        uint      `json:"id" `
	Name      string    `json:"name" `
	Email     string    `json:"email" `
	Password  string    `json:"password" `
	CreatedAt time.Time `json:"created_at"`
}
