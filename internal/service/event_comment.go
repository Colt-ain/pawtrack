package service

import (
	"errors"

	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/repository"
)

type EventCommentService interface {
	CreateComment(req *dto.CreateCommentRequest, userID uint, role models.UserRole) (*models.EventComment, error)
	GetComment(id uint, userID uint, role models.UserRole) (*dto.CommentResponse, error)
	UpdateComment(id uint, req *dto.UpdateCommentRequest, userID uint, role models.UserRole) (*models.EventComment, error)
	DeleteComment(id uint, userID uint, role models.UserRole) error
	ListComments(eventID uint, userID uint, role models.UserRole) (*dto.CommentListResponse, error)
}

type eventCommentService struct {
	commentRepo repository.EventCommentRepository
	eventRepo   repository.EventRepository
	dogRepo     repository.DogRepository
}

func NewEventCommentService(
	commentRepo repository.EventCommentRepository,
	eventRepo repository.EventRepository,
	dogRepo repository.DogRepository,
) EventCommentService {
	return &eventCommentService{
		commentRepo: commentRepo,
		eventRepo:   eventRepo,
		dogRepo:     dogRepo,
	}
}

// checkEventAccess verifies if user can access the event (owner or consultant with access)
func (s *eventCommentService) checkEventAccess(eventID uint, userID uint, role models.UserRole) error {
	// Admin has access to everything
	if role == models.RoleAdmin {
		return nil
	}

	// Get event with dog info
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	// If event has no dog, deny access (shouldn't happen but just in case)
	if event.DogID == nil {
		return errors.New("event has no associated dog")
	}

	// Owner check
	if role == models.RoleOwner {
		dog, err := s.dogRepo.GetByID(*event.DogID)
		if err != nil {
			return err
		}
		if dog.OwnerID == userID {
			return nil
		}
		return errors.New("not authorized")
	}

	// Consultant check
	if role == models.RoleConsultant {
		hasAccess, err := s.dogRepo.HasConsultantAccess(userID, *event.DogID)
		if err != nil {
			return err
		}
		if hasAccess {
			return nil
		}
		return errors.New("not authorized")
	}

	return errors.New("not authorized")
}

func (s *eventCommentService) CreateComment(req *dto.CreateCommentRequest, userID uint, role models.UserRole) (*models.EventComment, error) {
	// Check if user has access to the event
	if err := s.checkEventAccess(req.EventID, userID, role); err != nil {
		return nil, err
	}

	comment := &models.EventComment{
		EventID: req.EventID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *eventCommentService) GetComment(id uint, userID uint, role models.UserRole) (*dto.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user has access to the event
	if err := s.checkEventAccess(comment.EventID, userID, role); err != nil {
		return nil, err
	}

	return s.toDTO(comment), nil
}

func (s *eventCommentService) UpdateComment(id uint, req *dto.UpdateCommentRequest, userID uint, role models.UserRole) (*models.EventComment, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Only author or admin can update
	if role != models.RoleAdmin && comment.UserID != userID {
		return nil, errors.New("only comment author can update")
	}

	comment.Content = req.Content

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *eventCommentService) DeleteComment(id uint, userID uint, role models.UserRole) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Only author or admin can delete
	if role != models.RoleAdmin && comment.UserID != userID {
		return errors.New("only comment author can delete")
	}

	return s.commentRepo.Delete(id)
}

func (s *eventCommentService) ListComments(eventID uint, userID uint, role models.UserRole) (*dto.CommentListResponse, error) {
	// Check if user has access to the event
	if err := s.checkEventAccess(eventID, userID, role); err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.ListByEvent(eventID)
	if err != nil {
		return nil, err
	}

	var commentDTOs []dto.CommentResponse
	for _, comment := range comments {
		commentDTOs = append(commentDTOs, *s.toDTO(&comment))
	}

	return &dto.CommentListResponse{
		Comments: commentDTOs,
		Count:    len(commentDTOs),
	}, nil
}

func (s *eventCommentService) toDTO(comment *models.EventComment) *dto.CommentResponse {
	resp := &dto.CommentResponse{
		ID:        comment.ID,
		EventID:   comment.EventID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	if comment.User != nil {
		resp.UserName = comment.User.Name
		resp.UserRole = string(comment.User.Role)
	}

	return resp
}
