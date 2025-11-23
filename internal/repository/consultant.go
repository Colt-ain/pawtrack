package repository

import (
	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"gorm.io/gorm"
)

type ConsultantRepository interface {
	UpdateProfile(profile *models.ConsultantProfile) error
	GetProfile(userID uint) (*models.ConsultantProfile, error)
	Search(criteria *dto.ConsultantSearchRequest) ([]models.ConsultantProfile, int64, error)
	CreateInvite(invite *models.Invite) error
	GetInviteByToken(token string) (*models.Invite, error)
	UpdateInviteStatus(invite *models.Invite) error
}

type consultantRepository struct {
	db *gorm.DB
}

func NewConsultantRepository(db *gorm.DB) ConsultantRepository {
	return &consultantRepository{db: db}
}

func (r *consultantRepository) UpdateProfile(profile *models.ConsultantProfile) error {
	return r.db.Save(profile).Error
}

func (r *consultantRepository) GetProfile(userID uint) (*models.ConsultantProfile, error) {
	var profile models.ConsultantProfile
	err := r.db.Preload("User").First(&profile, userID).Error
	return &profile, err
}

func (r *consultantRepository) Search(criteria *dto.ConsultantSearchRequest) ([]models.ConsultantProfile, int64, error) {
	var profiles []models.ConsultantProfile
	var totalCount int64

	query := r.db.Model(&models.ConsultantProfile{}).Preload("User").Joins("JOIN users ON users.id = consultant_profiles.user_id")

	// Ensure we only search active consultants (optional, but good practice)
	query = query.Where("users.role = ?", models.RoleConsultant)

	if criteria.Query != "" {
		search := "%" + criteria.Query + "%"
		query = query.Where("users.name ILIKE ? OR consultant_profiles.surname ILIKE ? OR consultant_profiles.description ILIKE ?", search, search, search)
	}

	if criteria.Services != "" {
		query = query.Where("consultant_profiles.services ILIKE ?", "%"+criteria.Services+"%")
	}

	if criteria.Breeds != "" {
		query = query.Where("consultant_profiles.breeds ILIKE ?", "%"+criteria.Breeds+"%")
	}

	if criteria.Location != "" {
		query = query.Where("consultant_profiles.location ILIKE ?", "%"+criteria.Location+"%")
	}

	query.Count(&totalCount)

	offset := (criteria.Page - 1) * criteria.PageSize
	err := query.Offset(offset).Limit(criteria.PageSize).Find(&profiles).Error

	return profiles, totalCount, err
}

func (r *consultantRepository) CreateInvite(invite *models.Invite) error {
	return r.db.Create(invite).Error
}

func (r *consultantRepository) GetInviteByToken(token string) (*models.Invite, error) {
	var invite models.Invite
	err := r.db.Where("token = ?", token).First(&invite).Error
	return &invite, err
}

func (r *consultantRepository) UpdateInviteStatus(invite *models.Invite) error {
	return r.db.Save(invite).Error
}
