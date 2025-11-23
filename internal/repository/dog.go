package repository

import (
	"time"

	"github.com/you/pawtrack/internal/models"
	"gorm.io/gorm"
)

// DogRepository interface for working with dogs
type DogRepository interface {
	Create(dog *models.Dog) error
	List() ([]models.Dog, error)
	GetByID(id uint) (*models.Dog, error)
	Update(dog *models.Dog) error
	Delete(id uint) error
	HasConsultantAccess(consultantID, dogID uint) (bool, error)
	GrantConsultantAccess(consultantID, dogID uint) error
}

// dogRepository implementation of the dog repository
type dogRepository struct {
	db *gorm.DB
}

// NewDogRepository creates a new dog repository
func NewDogRepository(db *gorm.DB) DogRepository {
	return &dogRepository{db: db}
}

// Create creates a new dog
func (r *dogRepository) Create(dog *models.Dog) error {
	return r.db.Create(dog).Error
}

// List returns a list of all dogs
func (r *dogRepository) List() ([]models.Dog, error) {
	var dogs []models.Dog
	err := r.db.Find(&dogs).Error
	return dogs, err
}

// GetByID returns a dog by ID
func (r *dogRepository) GetByID(id uint) (*models.Dog, error) {
	var dog models.Dog
	err := r.db.First(&dog, id).Error
	if err != nil {
		return nil, err
	}
	return &dog, nil
}

// Update updates a dog's data
func (r *dogRepository) Update(dog *models.Dog) error {
	return r.db.Save(dog).Error
}

// Delete deletes a dog
func (r *dogRepository) Delete(id uint) error {
	return r.db.Delete(&models.Dog{}, id).Error
}

// HasConsultantAccess checks if a consultant has access to a dog
func (r *dogRepository) HasConsultantAccess(consultantID, dogID uint) (bool, error) {
	var count int64
	err := r.db.Table("consultant_access").
		Where("consultant_id = ? AND dog_id = ? AND revoked_at IS NULL", consultantID, dogID).
		Count(&count).Error
	return count > 0, err
}

// GrantConsultantAccess grants access for a consultant to a dog
func (r *dogRepository) GrantConsultantAccess(consultantID, dogID uint) error {
	access := models.ConsultantAccess{
		ConsultantID: consultantID,
		DogID:        dogID,
		GrantedAt:    time.Now(),
	}
	return r.db.Create(&access).Error
}
