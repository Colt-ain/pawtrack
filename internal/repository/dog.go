package repository

import (
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
