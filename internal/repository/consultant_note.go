package repository

import (
	"time"

	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"gorm.io/gorm"
)

type ConsultantNoteRepository interface {
	Create(note *models.ConsultantNote) error
	GetByID(id uint) (*models.ConsultantNote, error)
	Update(note *models.ConsultantNote) error
	Delete(id uint) error
	List(filters *dto.NoteFilterParams, consultantID uint, isAdmin bool) ([]models.ConsultantNote, int64, error)
}

type consultantNoteRepository struct {
	db *gorm.DB
}

func NewConsultantNoteRepository(db *gorm.DB) ConsultantNoteRepository {
	return &consultantNoteRepository{db: db}
}

func (r *consultantNoteRepository) Create(note *models.ConsultantNote) error {
	return r.db.Create(note).Error
}

func (r *consultantNoteRepository) GetByID(id uint) (*models.ConsultantNote, error) {
	var note models.ConsultantNote
	err := r.db.Preload("Dog.Owner").Preload("Consultant").First(&note, id).Error
	return &note, err
}

func (r *consultantNoteRepository) Update(note *models.ConsultantNote) error {
	return r.db.Save(note).Error
}

func (r *consultantNoteRepository) Delete(id uint) error {
	return r.db.Delete(&models.ConsultantNote{}, id).Error
}

func (r *consultantNoteRepository) List(filters *dto.NoteFilterParams, consultantID uint, isAdmin bool) ([]models.ConsultantNote, int64, error) {
	var notes []models.ConsultantNote
	var totalCount int64

	query := r.db.Model(&models.ConsultantNote{}).
		Joins("LEFT JOIN dogs ON dogs.id = consultant_notes.dog_id").
		Joins("LEFT JOIN users ON users.id = dogs.owner_id").
		Preload("Dog.Owner").
		Preload("Consultant")

	// RBAC: Consultants see only their notes, admins see all
	if !isAdmin {
		query = query.Where("consultant_notes.consultant_id = ?", consultantID)
	}

	// Search filter
	if filters.Search != "" {
		search := "%" + filters.Search + "%"
		query = query.Where("consultant_notes.title ILIKE ? OR consultant_notes.content ILIKE ?", search, search)
	}

	// Dog filter
	if filters.DogID > 0 {
		query = query.Where("consultant_notes.dog_id = ?", filters.DogID)
	}

	// Owner filter
	if filters.OwnerID > 0 {
		query = query.Where("dogs.owner_id = ?", filters.OwnerID)
	}

	// Date filters
	if filters.FromDate != "" {
		fromTime, err := time.Parse(time.RFC3339, filters.FromDate)
		if err == nil {
			query = query.Where("consultant_notes.created_at >= ?", fromTime)
		}
	}

	if filters.ToDate != "" {
		toTime, err := time.Parse(time.RFC3339, filters.ToDate)
		if err == nil {
			query = query.Where("consultant_notes.created_at <= ?", toTime)
		}
	}

	// Count total
	query.Count(&totalCount)

	// Sorting
	sortField := "consultant_notes.created_at"
	switch filters.SortBy {
	case "updated_at":
		sortField = "consultant_notes.updated_at"
	case "dog_name":
		sortField = "dogs.name"
	case "owner_name":
		sortField = "users.name"
	}

	order := "DESC"
	if filters.Order == "asc" {
		order = "ASC"
	}

	query = query.Order(sortField + " " + order)

	// Pagination
	offset := (filters.Page - 1) * filters.PageSize
	err := query.Offset(offset).Limit(filters.PageSize).Find(&notes).Error

	return notes, totalCount, err
}
