package service

import (
	"time"

	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/repository"
)

// EventService interface for event business logic
type EventService interface {
	CreateEvent(req *dto.CreateEventRequest) (*models.Event, error)
	ListEvents(limit, offset int, eventType string) ([]models.Event, error)
	GetEvent(id uint) (*models.Event, error)
	DeleteEvent(id uint) error
}

// eventService implementation of the event service
type eventService struct {
	repo repository.EventRepository
}

// NewEventService creates a new event service
func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{repo: repo}
}

// CreateEvent creates a new event
func (s *eventService) CreateEvent(req *dto.CreateEventRequest) (*models.Event, error) {
	when := time.Now().UTC()
	if req.At != nil {
		when = req.At.UTC()
	}

	event := &models.Event{
		Type: req.Type,
		Note: req.Note,
		At:   when,
	}

	err := s.repo.Create(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// ListEvents returns a list of events
func (s *eventService) ListEvents(limit, offset int, eventType string) ([]models.Event, error) {
	return s.repo.List(limit, offset, eventType)
}

// GetEvent returns an event by ID
func (s *eventService) GetEvent(id uint) (*models.Event, error) {
	return s.repo.GetByID(id)
}

// DeleteEvent deletes an event
func (s *eventService) DeleteEvent(id uint) error {
	return s.repo.Delete(id)
}
