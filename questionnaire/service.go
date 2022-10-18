package questionnaire

import (
	"context"
	"strings"

	"github.com/muktihari/quiz_master/pkg/textinput"
)

type Service interface {
	// GetByID gets question by id, returns error if any
	GetByID(ctx context.Context, ID int) (*Question, error)
	// GetAll gets all questions, returns error if any
	GetAll(ctx context.Context) ([]Question, error)
	// Create creates question, returns error if any
	Create(ctx context.Context, question *Question) error
	// Update updates existing question, return error if any
	Update(ctx context.Context, question *Question) error
	// Delete deletes existing question, return error if any
	Delete(ctx context.Context, ID int) error
	// Answer checks whether given answer to specific question is correct
	Answer(ctx context.Context, ID int, answer string) (bool, error)
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

type service struct {
	repository Repository
}

func (s *service) GetByID(ctx context.Context, ID int) (*Question, error) {
	return s.repository.GetByID(ctx, ID)
}

func (s *service) GetAll(ctx context.Context) ([]Question, error) {
	return s.repository.GetAll(ctx)
}

func (s *service) Create(ctx context.Context, question *Question) error {
	question.Question = strings.Trim(question.Question, "\"")
	question.Answer = strings.Trim(question.Answer, "\"")

	return s.repository.Create(ctx, question)
}

func (s *service) Update(ctx context.Context, question *Question) error {
	question.Question = strings.Trim(question.Question, "\"")
	question.Answer = strings.Trim(question.Answer, "\"")

	return s.repository.Update(ctx, question)
}

func (s *service) Delete(ctx context.Context, ID int) error {
	return s.repository.Delete(ctx, ID)
}

func (s *service) Answer(ctx context.Context, ID int, answer string) (bool, error) {
	question, err := s.GetByID(ctx, ID)
	if err != nil {
		return false, err
	}
	question.Answer = textinput.RecognizedAsNumber(question.Answer)
	answer = textinput.RecognizedAsNumber(strings.Trim(answer, "\""))

	return question.Answer == answer, nil
}
