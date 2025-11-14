package repository

import (
	"errors"
	"hitalent-test/internal/domain"

	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

type QuestionRepository interface {
	Create(question *domain.Question) error
	GetByID(id uint) (*domain.Question, error)
	GetAll() ([]domain.Question, error)
	Delete(id uint) error
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(question *domain.Question) error {
	return r.db.Create(question).Error
}

func (r *questionRepository) GetByID(id uint) (*domain.Question, error) {
	var question domain.Question
	err := r.db.Preload("Answers").First(&question, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrQuestionNotFound
	}
	return &question, err
}

func (r *questionRepository) GetAll() ([]domain.Question, error) {
	var questions []domain.Question
	err := r.db.Find(&questions).Error
	return questions, err
}

func (r *questionRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Question{}, id)
	if result.RowsAffected == 0 {
		return domain.ErrQuestionNotFound
	}
	return result.Error
}
