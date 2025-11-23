package service

import (
	"math"
	"time"

	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/repository"
)

// EventService interface for event business logic
type EventService interface {
	CreateEvent(req *dto.CreateEventRequest) (*models.Event, error)
	ListEvents(filters *dto.EventFilterParams) (*dto.EventListResponse, error)
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
		DogID: req.DogID,
		Type:  req.Type,
		Note:  req.Note,
		At:    when,
		AttachmentURL: req.AttachmentURL,
	}

	err := s.repo.Create(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// ListEvents returns a list of events with filtering and pagination
func (s *eventService) ListEvents(filters *dto.EventFilterParams) (*dto.EventListResponse, error) {
	// Set defaults
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}

	events, totalCount, err := s.repo.List(filters)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(filters.PageSize)))

	return &dto.EventListResponse{
		Events:     events,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// GetEvent returns an event by ID
func (s *eventService) GetEvent(id uint) (*models.Event, error) {
	return s.repo.GetByID(id)
}

// DeleteEvent deletes an event
func (s *eventService) DeleteEvent(id uint) error {
	return s.repo.Delete(id)
}
