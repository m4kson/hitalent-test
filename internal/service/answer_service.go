package service

import (
	"fmt"
	"hitalent-test/internal/domain"
	"hitalent-test/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type AnswerService struct {
	answerRepo   repository.AnswerRepository
	questionRepo repository.QuestionRepository
}

func NewAnswerService(answerRepo repository.AnswerRepository, questionRepo repository.QuestionRepository) *AnswerService {
	return &AnswerService{
		answerRepo:   answerRepo,
		questionRepo: questionRepo,
	}
}

func (s *AnswerService) Create(questionID uint, req *domain.CreateAnswerRequest) (*domain.Answer, error) {
	_, err := s.questionRepo.GetByID(questionID)
	if err != nil {
		return nil, err
	}

	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	answer := &domain.Answer{
		QuestionID: questionID,
		UserID:     strings.TrimSpace(req.UserID),
		Text:       strings.TrimSpace(req.Text),
	}

	if err := s.answerRepo.Create(answer); err != nil {
		return nil, fmt.Errorf("failed to create answer: %w", err)
	}

	return answer, nil
}

func (s *AnswerService) GetByID(id uint) (*domain.Answer, error) {
	return s.answerRepo.GetByID(id)
}

func (s *AnswerService) Delete(id uint) error {
	return s.answerRepo.Delete(id)
}

func (s *AnswerService) validateCreateRequest(req *domain.CreateAnswerRequest) error {
	userID := strings.TrimSpace(req.UserID)
	if userID == "" {
		return fmt.Errorf("%w: user_id is required", domain.ErrInvalidInput)
	}

	if _, err := uuid.Parse(userID); err != nil {
		return fmt.Errorf("%w: user_id must be a valid UUID", domain.ErrInvalidInput)
	}

	text := strings.TrimSpace(req.Text)
	if text == "" {
		return fmt.Errorf("%w: answer text is required", domain.ErrInvalidInput)
	}
	if len(text) < 5 {
		return fmt.Errorf("%w: answer text must be at least 5 characters", domain.ErrInvalidInput)
	}
	if len(text) > 1000 {
		return fmt.Errorf("%w: answer text must not exceed 1000 characters", domain.ErrInvalidInput)
	}

	return nil
}
