package repository

import (
	"errors"
	"hitalent-test/internal/domain"

	"gorm.io/gorm"
)

type answerRepository struct {
	db *gorm.DB
}

type AnswerRepository interface {
	Create(answer *domain.Answer) error
	GetByID(id uint) (*domain.Answer, error)
	GetByQuestionID(questionID uint) ([]domain.Answer, error)
	Delete(id uint) error
}

func NewAnswerRepository(db *gorm.DB) AnswerRepository {
	return &answerRepository{db: db}
}

func (r *answerRepository) Create(answer *domain.Answer) error {
	return r.db.Create(answer).Error
}

func (r *answerRepository) GetByID(id uint) (*domain.Answer, error) {
	var answer domain.Answer
	err := r.db.First(&answer, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrAnswerNotFound
	}
	return &answer, err
}

func (r *answerRepository) GetByQuestionID(questionID uint) ([]domain.Answer, error) {
	var answers []domain.Answer
	err := r.db.Where("question_id = ?", questionID).Find(&answers).Error
	return answers, err
}

func (r *answerRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Answer{}, id)
	if result.RowsAffected == 0 {
		return domain.ErrAnswerNotFound
	}
	return result.Error
}
