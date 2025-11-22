package service

import (
	"time"

	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/repository"
)

// DogService interface for dog business logic
type DogService interface {
	CreateDog(req *dto.CreateDogRequest, userID uint) (*models.Dog, error)
	ListDogs() ([]models.Dog, error)
	GetDog(id uint) (*models.Dog, error)
	UpdateDog(id uint, req *dto.UpdateDogRequest) (*models.Dog, error)
	DeleteDog(id uint) error
}

// dogService implementation of the dog service
type dogService struct {
	repo repository.DogRepository
}

// NewDogService creates a new dog service
func NewDogService(repo repository.DogRepository) DogService {
	return &dogService{repo: repo}
}

// CreateDog creates a new dog
func (s *dogService) CreateDog(req *dto.CreateDogRequest, userID uint) (*models.Dog, error) {
	var birthDate time.Time
	if req.BirthDate != "" {
		var err error
		birthDate, err = time.Parse(time.RFC3339, req.BirthDate)
		if err != nil {
			return nil, err
		}
	}

	dog := &models.Dog{
		OwnerID:   userID,
		Name:      req.Name,
		Breed:     req.Breed,
		BirthDate: birthDate,
	}

	err := s.repo.Create(dog)
	if err != nil {
		return nil, err
	}

	return dog, nil
}

// ListDogs returns a list of all dogs
func (s *dogService) ListDogs() ([]models.Dog, error) {
	return s.repo.List()
}

// GetDog returns a dog by ID
func (s *dogService) GetDog(id uint) (*models.Dog, error) {
	return s.repo.GetByID(id)
}

// UpdateDog updates a dog's data
func (s *dogService) UpdateDog(id uint, req *dto.UpdateDogRequest) (*models.Dog, error) {
	dog, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	dog.Name = req.Name
	dog.Breed = req.Breed

	if req.BirthDate != "" {
		birthDate, err := time.Parse(time.RFC3339, req.BirthDate)
		if err != nil {
			return nil, err
		}
		dog.BirthDate = birthDate
	}

	err = s.repo.Update(dog)
	if err != nil {
		return nil, err
	}

	return dog, nil
}

// DeleteDog deletes a dog
func (s *dogService) DeleteDog(id uint) error {
	return s.repo.Delete(id)
}
