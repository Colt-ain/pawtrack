package models

import "time"

// Permission represents a system permission
type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;uniqueIndex;not null"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserPermission represents the many-to-many relationship between users and permissions
type UserPermission struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	User         *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PermissionID uint      `json:"permission_id" gorm:"not null;index"`
	Permission   *Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID"`
	GrantedAt    time.Time `json:"granted_at"`
}

// TableName overrides the table name for UserPermission
func (UserPermission) TableName() string {
	return "user_permissions"
}
