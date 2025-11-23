package repository

import (
	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"gorm.io/gorm"
	"strings"
)

// EventRepository interface for event data access
type EventRepository interface {
	Create(event *models.Event) error
	List(filters *dto.EventFilterParams) ([]models.Event, int64, error)
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

// List returns a list of events based on filters
func (r *eventRepository) List(filters *dto.EventFilterParams) ([]models.Event, int64, error) {
	var events []models.Event
	var totalCount int64

	query := r.db.Model(&models.Event{})

	// Preload dog relationship
	query = query.Preload("Dog")

	// Role-based data access control
	switch filters.UserRole {
	case models.RoleAdmin:
		// Admin sees all events, no extra filter needed
	case models.RoleOwner:
		// Owner sees only events for their dogs
		// JOIN dogs ON events.dog_id = dogs.id WHERE dogs.owner_id = ?
		query = query.Joins("JOIN dogs ON events.dog_id = dogs.id").
			Where("dogs.owner_id = ?", filters.UserID)
	case models.RoleConsultant:
		// Consultant sees only events for dogs they have access to
		// JOIN consultant_access ca ON events.dog_id = ca.dog_id WHERE ca.consultant_id = ? AND ca.revoked_at IS NULL
		query = query.Joins("JOIN consultant_access ca ON events.dog_id = ca.dog_id").
			Where("ca.consultant_id = ? AND ca.revoked_at IS NULL", filters.UserID)
	default:
		// Unknown role or no role - return empty or error?
		// For safety, return nothing if role is not recognized
		return []models.Event{}, 0, nil
	}

	// Date range
	if filters.FromDate != nil {
		query = query.Where("events.at >= ?", filters.FromDate)
	}
	if filters.ToDate != nil {
		query = query.Where("events.at <= ?", filters.ToDate)
	}

	// Event types
	if filters.Types != "" {
		types := strings.Split(filters.Types, ",")
		query = query.Where("events.type IN ?", types)
	}

	// Text search (notes or dog name)
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		// We need to join dogs if not already joined for owner role
		// But GORM handles duplicate joins smartly usually, or we can check
		// For simplicity, let's use a subquery or left join if needed
		// Since we might have already joined dogs for Owner role, we should be careful.
		// However, for search we need dog name.
		// Let's use LEFT JOIN for dogs if not owner (Owner already has INNER JOIN)
		if filters.UserRole != models.RoleOwner {
			query = query.Joins("LEFT JOIN dogs d_search ON events.dog_id = d_search.id")
			query = query.Where("events.note LIKE ? OR d_search.name LIKE ?", searchTerm, searchTerm)
		} else {
			// Already joined as 'dogs'
			query = query.Where("events.note LIKE ? OR dogs.name LIKE ?", searchTerm, searchTerm)
		}
	}

	// Dog name exact match
	if filters.DogName != "" {
		if filters.UserRole != models.RoleOwner {
			query = query.Joins("LEFT JOIN dogs d_name ON events.dog_id = d_name.id")
			query = query.Where("d_name.name = ?", filters.DogName)
		} else {
			query = query.Where("dogs.name = ?", filters.DogName)
		}
	}

	// Count total before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	sortField := "created_at"
	if filters.SortBy != "" {
		sortField = filters.SortBy
	}
	// Prefix with table name to avoid ambiguity
	sortField = "events." + sortField

	sortOrder := "desc"
	if filters.SortOrder != "" {
		sortOrder = filters.SortOrder
	}
	query = query.Order(sortField + " " + sortOrder)

	// Pagination
	offset := (filters.Page - 1) * filters.PageSize
	query = query.Offset(offset).Limit(filters.PageSize)

	err := query.Find(&events).Error
	return events, totalCount, err
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
