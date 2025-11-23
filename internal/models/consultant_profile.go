package models

// ConsultantProfile represents extended profile information for a consultant
type ConsultantProfile struct {
	UserID      uint   `json:"user_id" gorm:"primaryKey"`
	User        *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Description string `json:"description" gorm:"type:text"`
	Services    string `json:"services" gorm:"type:text"` // Comma-separated list of services
	Breeds      string `json:"breeds" gorm:"type:text"`   // Comma-separated list of breeds
	Location    string `json:"location" gorm:"size:255"`
	Surname     string `json:"surname" gorm:"size:255"`
}
