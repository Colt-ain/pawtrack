package repository

import (
	"github.com/you/pawtrack/internal/models"
	"gorm.io/gorm"
)

type EventCommentRepository interface {
	Create(comment *models.EventComment) error
	GetByID(id uint) (*models.EventComment, error)
	Update(comment *models.EventComment) error
	Delete(id uint) error
	ListByEvent(eventID uint) ([]models.EventComment, error)
}

type eventCommentRepository struct {
	db *gorm.DB
}

func NewEventCommentRepository(db *gorm.DB) EventCommentRepository {
	return &eventCommentRepository{db: db}
}

func (r *eventCommentRepository) Create(comment *models.EventComment) error {
	return r.db.Create(comment).Error
}

func (r *eventCommentRepository) GetByID(id uint) (*models.EventComment, error) {
	var comment models.EventComment
	err := r.db.Preload("User").Preload("Event.Dog").First(&comment, id).Error
	return &comment, err
}

func (r *eventCommentRepository) Update(comment *models.EventComment) error {
	return r.db.Save(comment).Error
}

func (r *eventCommentRepository) Delete(id uint) error {
	return r.db.Delete(&models.EventComment{}, id).Error
}

func (r *eventCommentRepository) ListByEvent(eventID uint) ([]models.EventComment, error) {
	var comments []models.EventComment
	err := r.db.Where("event_id = ?", eventID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}
