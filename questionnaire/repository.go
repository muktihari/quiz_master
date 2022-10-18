package questionnaire

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrQuestionNotFound       = errors.New("question not found")
	ErrQuestionIsAlreadyExist = errors.New("question is already exist")
)

type Repository interface {
	// GetByID gets question by id, returns error if any
	GetByID(ctx context.Context, id int) (*Question, error)
	// GetAll gets questions, returns error if any
	GetAll(ctx context.Context) ([]Question, error)
	// Create creates question, returns error if any
	Create(ctx context.Context, question *Question) error
	// Update updates existing question, return error if any
	Update(ctx context.Context, question *Question) error
	// Delete deletes existing question, return error if any
	Delete(ctx context.Context, id int) error
}

func NewRepository() Repository {
	return &inmemRepository{
		questions: make([]Question, 0),
	}
}

type inmemRepository struct {
	mu        sync.RWMutex
	questions []Question
}

func (r *inmemRepository) GetByID(ctx context.Context, id int) (*Question, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, question := range r.questions {
		if id == question.ID {
			return &question, nil
		}
	}

	return nil, ErrQuestionNotFound
}

func (r *inmemRepository) GetAll(ctx context.Context) ([]Question, error) {
	return r.questions, nil
}

func (r *inmemRepository) Create(ctx context.Context, question *Question) error {
	_, err := r.GetByID(ctx, question.ID)
	if err == nil {
		return ErrQuestionIsAlreadyExist
	}

	r.mu.Lock()
	r.questions = append(r.questions, *question)
	r.mu.Unlock()

	return nil
}

func (r *inmemRepository) Update(ctx context.Context, question *Question) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.questions {
		if question.ID == r.questions[i].ID {
			r.questions[i] = *question
			return nil
		}
	}

	return ErrQuestionNotFound
}

func (r *inmemRepository) Delete(ctx context.Context, id int) error {
	var index *int

	r.mu.RLock()
	for i := range r.questions {
		if id == r.questions[i].ID {
			index = &i
			break
		}
	}
	r.mu.RUnlock()

	if index == nil {
		return ErrQuestionNotFound
	}

	r.mu.Lock()
	r.questions = append(r.questions[:*index], r.questions[*index+1:]...)
	r.mu.Unlock()

	return nil
}
