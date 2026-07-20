package event

import (
	"errors"

	"gorm.io/gorm"
)

var ErrEventNotFound = errors.New("event not found")

type Repository interface {
	Create(event *Event) error
	GetAll() ([]*Event, error)
	GetByID(eventId uint) (*Event, error)
	Update(event *Event) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(event *Event) error {
	// result := r.db.Create(event)

	// if result.Error != nil {
	// 	return result.Error
	// }

	// return nil

	return r.db.Create(event).Error

}

func (r *repository) GetAll() ([]*Event, error) {
	var events []*Event

	err := r.db.Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *repository) GetByID(eventId uint) (*Event, error) {
	var event Event

	err := r.db.First(&event, eventId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}

		return nil, err
	}

	return &event, nil
}

func (r *repository) Update(event *Event) error {
	return r.db.Save(event).Error
}
