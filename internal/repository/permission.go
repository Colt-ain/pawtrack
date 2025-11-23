package repository

import (
	"github.com/you/pawtrack/internal/models"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	// GetUserPermissions returns all permission names for a user
	GetUserPermissions(userID uint) ([]string, error)

	// GrantPermission grants a single permission to a user
	GrantPermission(userID uint, permissionName string) error

	// GrantPermissions grants multiple permissions to a user
	GrantPermissions(userID uint, permissionNames []string) error

	// RevokePermission revokes a permission from a user
	RevokePermission(userID uint, permissionName string) error

	// HasPermission checks if a user has a specific permission
	HasPermission(userID uint, permissionName string) (bool, error)

	// HasAnyPermission checks if user has at least one of the given permissions
	HasAnyPermission(userID uint, permissionNames []string) (bool, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) GetUserPermissions(userID uint) ([]string, error) {
	var permissions []string

	err := r.db.Table("user_permissions").
		Select("permissions.name").
		Joins("JOIN permissions ON permissions.id = user_permissions.permission_id").
		Where("user_permissions.user_id = ?", userID).
		Pluck("name", &permissions).Error

	return permissions, err
}

func (r *permissionRepository) GrantPermission(userID uint, permissionName string) error {
	// Find permission by name
	var permission models.Permission
	if err := r.db.Where("name = ?", permissionName).First(&permission).Error; err != nil {
		return err
	}

	// Check if already granted
	var count int64
	r.db.Model(&models.UserPermission{}).
		Where("user_id = ? AND permission_id = ?", userID, permission.ID).
		Count(&count)

	if count > 0 {
		return nil // Already granted
	}

	// Grant permission
	userPerm := &models.UserPermission{
		UserID:       userID,
		PermissionID: permission.ID,
	}

	return r.db.Create(userPerm).Error
}

func (r *permissionRepository) GrantPermissions(userID uint, permissionNames []string) error {
	for _, permName := range permissionNames {
		if err := r.GrantPermission(userID, permName); err != nil {
			// Log error but continue with other permissions
			continue
		}
	}
	return nil
}

func (r *permissionRepository) RevokePermission(userID uint, permissionName string) error {
	// Find permission by name
	var permission models.Permission
	if err := r.db.Where("name = ?", permissionName).First(&permission).Error; err != nil {
		return err
	}

	// Delete user permission
	return r.db.Where("user_id = ? AND permission_id = ?", userID, permission.ID).
		Delete(&models.UserPermission{}).Error
}

func (r *permissionRepository) HasPermission(userID uint, permissionName string) (bool, error) {
	var count int64

	err := r.db.Table("user_permissions").
		Joins("JOIN permissions ON permissions.id = user_permissions.permission_id").
		Where("user_permissions.user_id = ? AND permissions.name = ?", userID, permissionName).
		Count(&count).Error

	return count > 0, err
}

func (r *permissionRepository) HasAnyPermission(userID uint, permissionNames []string) (bool, error) {
	var count int64

	err := r.db.Table("user_permissions").
		Joins("JOIN permissions ON permissions.id = user_permissions.permission_id").
		Where("user_permissions.user_id = ? AND permissions.name IN ?", userID, permissionNames).
		Count(&count).Error

	return count > 0, err
}
