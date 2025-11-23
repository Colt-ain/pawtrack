package service

import (
	"errors"
	"strings"

	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/permissions"
	"github.com/you/pawtrack/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService interface for user business logic
type UserService interface {
	CreateUser(req *dto.CreateUserRequest) (*models.User, error)
	ListUsers() ([]models.User, error)
	GetUser(id uint) (*models.User, error)
	UpdateUser(id uint, req *dto.UpdateUserRequest) (*models.User, error)
	DeleteUser(id uint) error
}

// userService implementation of the user service
type userService struct {
	repo     repository.UserRepository
	permRepo repository.PermissionRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, permRepo repository.PermissionRepository) UserService {
	return &userService{
		repo:     repo,
		permRepo: permRepo,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(req *dto.CreateUserRequest) (*models.User, error) {
	// Set default role to owner if not provided
	role := req.Role
	if role == "" {
		role = models.RoleOwner
	}

	// Prevent admin role creation via API
	if role == models.RoleAdmin {
		return nil, errors.New("cannot register admin users via API")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Grant permissions based on role
	var permissionsToGrant []string
	switch role {
	case models.RoleOwner:
		permissionsToGrant = permissions.OwnerPermissions
	case models.RoleConsultant:
		permissionsToGrant = permissions.ConsultantBasePermissions
	case models.RoleAdmin:
		permissionsToGrant = permissions.AdminPermissions
	}

	if len(permissionsToGrant) > 0 {
		s.permRepo.GrantPermissions(user.ID, permissionsToGrant)
	}

	return user, nil
}

// ListUsers returns a list of all users
func (s *userService) ListUsers() ([]models.User, error) {
	return s.repo.List()
}

// GetUser returns a user by ID
func (s *userService) GetUser(id uint) (*models.User, error) {
	return s.repo.GetByID(id)
}

// UpdateUser updates a user's data
func (s *userService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hashedPassword)
	}

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

// IsDuplicateKeyError checks if the error is a database duplicate key error
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	// Check error message for different databases
	errStr := err.Error()
	return len(errStr) > 0 && (strings.HasPrefix(errStr, "UNIQUE") || strings.HasPrefix(errStr, "duplicate"))
}
