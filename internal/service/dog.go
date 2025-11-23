package service

import (
	"errors"
	"time"

	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/repository"
)

// DogService interface for dog business logic
type DogService interface {
	CreateDog(req *dto.CreateDogRequest, userID uint) (*models.Dog, error)
	ListDogs() ([]models.Dog, error)
	GetDog(id uint, userID uint, role models.UserRole) (*models.Dog, error)
	UpdateDog(id uint, req *dto.UpdateDogRequest, userID uint, role models.UserRole) (*models.Dog, error)
	DeleteDog(id uint, userID uint, role models.UserRole) error
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

// GetDog returns a dog by ID with RBAC check
func (s *dogService) GetDog(id uint, userID uint, role models.UserRole) (*models.Dog, error) {
	dog, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.checkAccess(dog, userID, role); err != nil {
		return nil, err
	}

	return dog, nil
}

// UpdateDog updates a dog's data with RBAC check
func (s *dogService) UpdateDog(id uint, req *dto.UpdateDogRequest, userID uint, role models.UserRole) (*models.Dog, error) {
	dog, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.checkAccess(dog, userID, role); err != nil {
		return nil, err
	}

	// Only Owner and Admin can update dog details? Or Consultant too?
	// Assuming Consultant can only VIEW for now, or maybe update notes (but notes are on events).
	// Let's restrict Update/Delete to Owner and Admin.
	if role == models.RoleConsultant {
		return nil, errors.New("unauthorized: consultants cannot update dogs")
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

// DeleteDog deletes a dog with RBAC check
func (s *dogService) DeleteDog(id uint, userID uint, role models.UserRole) error {
	dog, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.checkAccess(dog, userID, role); err != nil {
		return err
	}

	if role == models.RoleConsultant {
		return errors.New("unauthorized: consultants cannot delete dogs")
	}

	return s.repo.Delete(id)
}

func (s *dogService) checkAccess(dog *models.Dog, userID uint, role models.UserRole) error {
	if role == models.RoleAdmin {
		return nil
	}
	if role == models.RoleOwner {
		if dog.OwnerID != userID {
			return errors.New("unauthorized")
		}
		return nil
	}
	if role == models.RoleConsultant {
		hasAccess, err := s.repo.HasConsultantAccess(userID, dog.ID)
		if err != nil {
			return err
		}
		if !hasAccess {
			return errors.New("unauthorized")
		}
		return nil
	}
	return errors.New("unauthorized")
}
