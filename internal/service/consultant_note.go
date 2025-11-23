package service

import (
	"errors"
	"math"

	"github.com/you/pawtrack/internal/dto"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/repository"
)

type ConsultantNoteService interface {
	CreateNote(req *dto.CreateNoteRequest, consultantID uint) (*models.ConsultantNote, error)
	GetNote(id uint, userID uint, role models.UserRole) (*dto.NoteResponse, error)
	UpdateNote(id uint, req *dto.UpdateNoteRequest, userID uint, role models.UserRole) (*models.ConsultantNote, error)
	DeleteNote(id uint, userID uint, role models.UserRole) error
	ListNotes(filters *dto.NoteFilterParams, userID uint, role models.UserRole) (*dto.NoteListResponse, error)
}

type consultantNoteService struct {
	noteRepo repository.ConsultantNoteRepository
	dogRepo  repository.DogRepository
}

func NewConsultantNoteService(noteRepo repository.ConsultantNoteRepository, dogRepo repository.DogRepository) ConsultantNoteService {
	return &consultantNoteService{
		noteRepo: noteRepo,
		dogRepo:  dogRepo,
	}
}

func (s *consultantNoteService) CreateNote(req *dto.CreateNoteRequest, consultantID uint) (*models.ConsultantNote, error) {
	// Verify consultant has access to the dog
	hasAccess, err := s.dogRepo.HasConsultantAccess(consultantID, req.DogID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New("consultant does not have access to this dog")
	}

	note := &models.ConsultantNote{
		ConsultantID: consultantID,
		DogID:        req.DogID,
		Title:        req.Title,
		Content:      req.Content,
	}

	err = s.noteRepo.Create(note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *consultantNoteService) GetNote(id uint, userID uint, role models.UserRole) (*dto.NoteResponse, error) {
	note, err := s.noteRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// RBAC: Admins can see all, consultants only their own
	if role != models.RoleAdmin && note.ConsultantID != userID {
		return nil, errors.New("unauthorized")
	}

	return s.toDTO(note), nil
}

func (s *consultantNoteService) UpdateNote(id uint, req *dto.UpdateNoteRequest, userID uint, role models.UserRole) (*models.ConsultantNote, error) {
	note, err := s.noteRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// RBAC: Admins can update all, consultants only their own
	if role != models.RoleAdmin && note.ConsultantID != userID {
		return nil, errors.New("unauthorized")
	}

	// Update fields if provided
	if req.Title != "" {
		note.Title = req.Title
	}
	if req.Content != "" {
		note.Content = req.Content
	}

	err = s.noteRepo.Update(note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *consultantNoteService) DeleteNote(id uint, userID uint, role models.UserRole) error {
	note, err := s.noteRepo.GetByID(id)
	if err != nil {
		return err
	}

	// RBAC: Admins can delete all, consultants only their own
	if role != models.RoleAdmin && note.ConsultantID != userID {
		return errors.New("unauthorized")
	}

	return s.noteRepo.Delete(id)
}

func (s *consultantNoteService) ListNotes(filters *dto.NoteFilterParams, userID uint, role models.UserRole) (*dto.NoteListResponse, error) {
	// Set defaults
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}

	isAdmin := role == models.RoleAdmin
	notes, totalCount, err := s.noteRepo.List(filters, userID, isAdmin)
	if err != nil {
		return nil, err
	}

	var noteDTOs []dto.NoteResponse
	for _, note := range notes {
		noteDTOs = append(noteDTOs, *s.toDTO(&note))
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(filters.PageSize)))

	return &dto.NoteListResponse{
		Notes:      noteDTOs,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (s *consultantNoteService) toDTO(note *models.ConsultantNote) *dto.NoteResponse {
	resp := &dto.NoteResponse{
		ID:           note.ID,
		ConsultantID: note.ConsultantID,
		DogID:        note.DogID,
		Title:        note.Title,
		Content:      note.Content,
		CreatedAt:    note.CreatedAt,
		UpdatedAt:    note.UpdatedAt,
	}

	if note.Dog != nil {
		resp.DogName = note.Dog.Name
		if note.Dog.Owner != nil {
			resp.OwnerID = note.Dog.OwnerID
			resp.OwnerName = note.Dog.Owner.Name
		}
	}

	return resp
}
