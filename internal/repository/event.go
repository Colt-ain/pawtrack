package repository

import (
	"github.com/you/pawtrack/internal/models"
	"gorm.io/gorm"
)

// EventRepository interface for working with events
type EventRepository interface {
	Create(event *models.Event) error
	List(limit, offset int, eventType string) ([]models.Event, error)
	GetByID(id uint) (*models.Event, error)
	Delete(id uint) error
}

// eventRepository implementation of the event repository
type eventRepository struct {
	db *gorm.DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

// Create creates a new event
func (r *eventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

// List returns a list of events with filtering and pagination
func (r *eventRepository) List(limit, offset int, eventType string) ([]models.Event, error) {
	var events []models.Event
	query := r.db.Order("at desc").Limit(limit).Offset(offset)
	
	if eventType != "" {
		query = query.Where("type = ?", eventType)
	}
	
	err := query.Find(&events).Error
	return events, err
}

// GetByID returns an event by ID
func (r *eventRepository) GetByID(id uint) (*models.Event, error) {
	var event models.Event
	err := r.db.First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Delete deletes an event
func (r *eventRepository) Delete(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
}
