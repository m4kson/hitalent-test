package service

import (
	"testing"

	"hitalent-test/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockQuestionRepository struct {
	mock.Mock
}

func (m *MockQuestionRepository) Create(question *domain.Question) error {
	args := m.Called(question)
	return args.Error(0)
}

func (m *MockQuestionRepository) GetByID(id uint) (*domain.Question, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Question), args.Error(1)
}

func (m *MockQuestionRepository) GetAll() ([]domain.Question, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Question), args.Error(1)
}

func (m *MockQuestionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestQuestionService_Create_ValidQuestion(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	mockRepo.On("Create", mock.MatchedBy(func(q *domain.Question) bool {
		return q.Text == "What is the capital of France?"
	})).Return(nil)

	service := NewQuestionService(mockRepo)

	req := &domain.CreateQuestionRequest{
		Text: "What is the capital of France?",
	}

	question, err := service.Create(req)

	require.NoError(t, err)
	assert.NotNil(t, question)
	assert.Equal(t, "What is the capital of France?", question.Text)
	mockRepo.AssertExpectations(t)
}

func TestQuestionService_Create_InvalidCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty text",
			input:       "",
			expectError: true,
			errorMsg:    "required",
		},
		{
			name:        "too short",
			input:       "Hi?",
			expectError: true,
			errorMsg:    "10 characters",
		},
		{
			name:        "too long",
			input:       string(make([]byte, 1001)),
			expectError: true,
			errorMsg:    "1000 characters",
		},
		{
			name:        "whitespace only",
			input:       "   ",
			expectError: true,
			errorMsg:    "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockQuestionRepository)
			service := NewQuestionService(mockRepo)

			req := &domain.CreateQuestionRequest{Text: tt.input}

			_, err := service.Create(req)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid input data")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQuestionService_GetByID(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	expectedQuestion := &domain.Question{ID: 1, Text: "Test"}
	mockRepo.On("GetByID", uint(1)).Return(expectedQuestion, nil)

	service := NewQuestionService(mockRepo)
	question, err := service.GetByID(1)

	require.NoError(t, err)
	assert.Equal(t, expectedQuestion, question)
	mockRepo.AssertExpectations(t)
}

func TestQuestionService_Delete(t *testing.T) {
	mockRepo := new(MockQuestionRepository)
	mockRepo.On("Delete", uint(1)).Return(nil)

	service := NewQuestionService(mockRepo)
	err := service.Delete(1)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
