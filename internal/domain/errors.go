package domain

import "errors"

var (
	ErrQuestionNotFound = errors.New("question not found")
	ErrAnswerNotFound   = errors.New("answer not found")
	ErrInvalidInput     = errors.New("invalid input data")
)
