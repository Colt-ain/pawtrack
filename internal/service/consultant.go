package service

import (
	"errors"
	"time"

	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/permissions"
	"github.com/you/pawtrack/internal/repository"
	"github.com/you/pawtrack/internal/utils"
	"gorm.io/gorm"
)

type ConsultantService interface {
	UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*models.ConsultantProfile, error)
	GetProfile(userID uint) (*dto.ConsultantProfileResponse, error)
	SearchConsultants(req *dto.ConsultantSearchRequest) ([]dto.ConsultantProfileResponse, int64, error)
	InviteConsultant(ownerID uint, consultantID uint, req *dto.CreateInviteRequest) (*dto.InviteResponse, error)
	AcceptInvite(token string, consultantID uint) error
}

type consultantService struct {
	repo     repository.ConsultantRepository
	dogRepo  repository.DogRepository
	permRepo repository.PermissionRepository
}

func NewConsultantService(repo repository.ConsultantRepository, dogRepo repository.DogRepository, permRepo repository.PermissionRepository) ConsultantService {
	return &consultantService{
		repo:     repo,
		dogRepo:  dogRepo,
		permRepo: permRepo,
	}
}

func (s *consultantService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*models.ConsultantProfile, error) {
	profile, err := s.repo.GetProfile(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new profile if not exists
			profile = &models.ConsultantProfile{UserID: userID}
		} else {
			return nil, err
		}
	}

	profile.Description = req.Description
	profile.Services = req.Services
	profile.Breeds = req.Breeds
	profile.Location = req.Location
	profile.Surname = req.Surname

	err = s.repo.UpdateProfile(profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *consultantService) GetProfile(userID uint) (*dto.ConsultantProfileResponse, error) {
	profile, err := s.repo.GetProfile(userID)
	if err != nil {
		return nil, err
	}

	return s.toDTO(profile), nil
}

func (s *consultantService) SearchConsultants(req *dto.ConsultantSearchRequest) ([]dto.ConsultantProfileResponse, int64, error) {
	profiles, count, err := s.repo.Search(req)
	if err != nil {
		return nil, 0, err
	}

	var dtos []dto.ConsultantProfileResponse
	for _, p := range profiles {
		dtos = append(dtos, *s.toDTO(&p))
	}

	return dtos, count, nil
}

func (s *consultantService) InviteConsultant(ownerID uint, consultantID uint, req *dto.CreateInviteRequest) (*dto.InviteResponse, error) {
	// Generate token (simple UUID or random string)
	token := utils.GenerateRandomString(32)

	invite := &models.Invite{
		OwnerID:      ownerID,
		ConsultantID: consultantID,
		DogID:        req.DogID,
		Token:        token,
		Status:       models.InvitePending,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour), // 24h expiry
	}

	err := s.repo.CreateInvite(invite)
	if err != nil {
		return nil, err
	}

	// Stub email sending
	// In a real app, we would use an email service here
	// fmt.Printf("Sending email to consultant %d: Invite link: /invites/accept?token=%s\n", consultantID, token)

	return &dto.InviteResponse{
		ID:           invite.ID,
		Token:        invite.Token, // Return token for E2E testing convenience, usually don't return it
		Status:       string(invite.Status),
		ConsultantID: invite.ConsultantID,
		DogID:        invite.DogID,
		ExpiresAt:    invite.ExpiresAt,
	}, nil
}

func (s *consultantService) AcceptInvite(token string, consultantID uint) error {
	invite, err := s.repo.GetInviteByToken(token)
	if err != nil {
		return err
	}

	if invite.Status != models.InvitePending {
		return errors.New("invite is not pending")
	}

	if time.Now().After(invite.ExpiresAt) {
		invite.Status = models.InviteRejected // Or expired
		s.repo.UpdateInviteStatus(invite)
		return errors.New("invite expired")
	}

	if invite.ConsultantID != consultantID {
		return errors.New("invite not for this consultant")
	}

	// Grant access
	// We need a method in DogRepository to grant access.
	// Since we don't have it exposed in DogRepository interface yet for writing,
	// let's assume we can add it or use a direct DB call if we were in repo.
	// But we are in service. We need to update DogRepository interface.
	// Wait, we have `consultant_access` table. We should use it.
	// Let's add `GrantConsultantAccess` to DogRepository.

	// Grant consultant access to the dog
	err = s.dogRepo.GrantConsultantAccess(consultantID, invite.DogID)
	if err != nil {
		return err
	}

	// Grant assigned permissions to consultant
	s.permRepo.GrantPermissions(consultantID, permissions.ConsultantAssignedPermissions)

	invite.Status = models.InviteAccepted
	return s.repo.UpdateInviteStatus(invite)
}

func (s *consultantService) toDTO(p *models.ConsultantProfile) *dto.ConsultantProfileResponse {
	return &dto.ConsultantProfileResponse{
		ID:          p.UserID,
		Name:        p.User.Name,
		Email:       p.User.Email,
		Surname:     p.Surname,
		Description: p.Description,
		Services:    p.Services,
		Breeds:      p.Breeds,
		Location:    p.Location,
	}
}
