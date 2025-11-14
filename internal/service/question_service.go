package service

import (
	"fmt"
	"hitalent-test/internal/domain"
	"hitalent-test/internal/repository"
	"strings"
)

type QuestionService struct {
	repo repository.QuestionRepository
}

func NewQuestionService(repo repository.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) Create(req *domain.CreateQuestionRequest) (*domain.Question, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	question := &domain.Question{
		Text: strings.TrimSpace(req.Text),
	}

	if err := s.repo.Create(question); err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}

	return question, nil
}

func (s *QuestionService) GetByID(id uint) (*domain.Question, error) {
	return s.repo.GetByID(id)
}

func (s *QuestionService) GetAll() ([]domain.Question, error) {
	return s.repo.GetAll()
}

func (s *QuestionService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *QuestionService) validateCreateRequest(req *domain.CreateQuestionRequest) error {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return fmt.Errorf("%w: question text is required", domain.ErrInvalidInput)
	}
	if len(text) < 10 {
		return fmt.Errorf("%w: question text must be at least 10 characters", domain.ErrInvalidInput)
	}
	if len(text) > 1000 {
		return fmt.Errorf("%w: question text must not exceed 1000 characters", domain.ErrInvalidInput)
	}
	return nil
}
