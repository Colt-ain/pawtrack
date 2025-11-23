package service

import (
	"log"

	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/permissions"
	"github.com/you/pawtrack/internal/repository"
)

// MigrateExistingUsers assigns appropriate permissions to all users based on their role.
// This should be run once after deployment of the atomic permissions system.
func MigrateExistingUsers(userRepo repository.UserRepository, permRepo repository.PermissionRepository) error {
	users, err := userRepo.List()
	if err != nil {
		return err
	}

	for _, user := range users {
		var perms []string
		switch user.Role {
		case models.RoleOwner:
			perms = permissions.OwnerPermissions
		case models.RoleConsultant:
			// Base permissions for consultants (before any invitations)
			perms = permissions.ConsultantBasePermissions
		case models.RoleAdmin:
			perms = permissions.AdminPermissions
		default:
			perms = []string{}
		}

		if len(perms) > 0 {
			if err := permRepo.GrantPermissions(user.ID, perms); err != nil {
				log.Printf("failed to grant permissions to user %d: %v", user.ID, err)
				// continue migrating other users
			}
		}
	}

	return nil
}
