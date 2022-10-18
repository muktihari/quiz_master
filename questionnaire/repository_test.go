package questionnaire

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInmemRepositoryGetByID(t *testing.T) {
	questions := []Question{
		{ID: 1, Question: "How many characters are there in \"Quipper\"?", Answer: "7"},
	}

	tt := []struct {
		Name             string
		QuestionID       int
		ExpectedQuestion *Question
		ExpectedErr      error
	}{
		{
			Name:             "get question ID 1, found",
			QuestionID:       1,
			ExpectedQuestion: &questions[0],
			ExpectedErr:      nil,
		},
		{
			Name:             "get question ID 2, not found",
			QuestionID:       2,
			ExpectedQuestion: nil,
			ExpectedErr:      ErrQuestionNotFound,
		},
	}

	var (
		ctx = context.Background()
		r   = (Repository)(&inmemRepository{questions: questions})
	)

	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			question, err := r.GetByID(ctx, tc.QuestionID)
			if !errors.Is(tc.ExpectedErr, err) {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.ExpectedQuestion, question); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInmemRepositoryGetAll(t *testing.T) {
	expectedQuestions := []Question{
		{ID: 1, Question: "How many characters are there in \"Quipper\"?", Answer: "7"},
		{ID: 2, Question: "Guess random number, 1, 2, 3 or 4?", Answer: "4"},
	}

	tt := []struct {
		Name              string
		ExpectedQuestions []Question
		ExpectedErr       error
	}{
		{
			Name:              "get all questions, found",
			ExpectedQuestions: expectedQuestions,
			ExpectedErr:       nil,
		},
	}

	var (
		ctx = context.Background()
		r   = (Repository)(&inmemRepository{questions: expectedQuestions})
	)

	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			questions, err := r.GetAll(ctx)
			if !errors.Is(tc.ExpectedErr, err) {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.ExpectedQuestions, questions); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInmemRepositoryCreate(t *testing.T) {
	var (
		question1 = Question{ID: 1, Question: "How many characters are there in \"Quipper\"?", Answer: "7"}
		question2 = Question{ID: 2, Question: "Guess random number: 1, 2, 3 or 4?", Answer: "4"}
		question3 = Question{ID: 3, Question: "Quipper vs Ruangguru?", Answer: "Quipper"}
	)

	tt := []struct {
		Name                          string
		Question                      *Question
		QuestionInmemDB               []Question
		ExpectedQuestionsAfterCreated []Question
		ExpectedErr                   error
	}{
		{
			Name:                          "create question 3",
			Question:                      &question3,
			QuestionInmemDB:               []Question{question1, question2},
			ExpectedQuestionsAfterCreated: []Question{question1, question2, question3},
			ExpectedErr:                   nil,
		},
		{
			Name:                          "create question 1, failed duplicate",
			Question:                      &question1,
			QuestionInmemDB:               []Question{question1, question2},
			ExpectedQuestionsAfterCreated: []Question{question1, question2},
			ExpectedErr:                   ErrQuestionIsAlreadyExist,
		},
	}

	ctx := context.Background()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			var (
				r         = NewRepository()
				inmemRepo = r.(*inmemRepository)
			)
			inmemRepo.questions = tc.QuestionInmemDB

			if err := r.Create(ctx, tc.Question); !errors.Is(tc.ExpectedErr, err) {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.ExpectedQuestionsAfterCreated, inmemRepo.questions); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInmemRepositoryUpdate(t *testing.T) {
	var (
		question1 = Question{ID: 1, Question: "How many characters are there in \"Quipper\"?", Answer: "7"}
		question2 = Question{ID: 2, Question: "Guess random number, 1, 2, 3 or 4?", Answer: "4"}
		question3 = Question{ID: 3, Question: "Quipper vs Ruangguru?", Answer: "Quipper"}

		updatedQuestion1 = Question{ID: 1, Question: "How many characterss in \"Quipper\"?", Answer: "7"}
	)

	tt := []struct {
		Name                          string
		Question                      *Question
		QuestionInmemDB               []Question
		ExpectedQuestionsAfterUpdated []Question
		ExpectedErr                   error
	}{
		{
			Name:                          "update question 1",
			Question:                      &updatedQuestion1,
			QuestionInmemDB:               []Question{question1, question2},
			ExpectedQuestionsAfterUpdated: []Question{updatedQuestion1, question2},
			ExpectedErr:                   nil,
		},
		{
			Name:                          "update question 3, failed not found",
			Question:                      &question3,
			QuestionInmemDB:               []Question{question1, question2},
			ExpectedQuestionsAfterUpdated: []Question{question1, question2},
			ExpectedErr:                   ErrQuestionNotFound,
		},
	}

	ctx := context.Background()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			var (
				r         = NewRepository()
				inmemRepo = r.(*inmemRepository)
			)
			inmemRepo.questions = tc.QuestionInmemDB

			if err := r.Update(ctx, tc.Question); !errors.Is(tc.ExpectedErr, err) {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.ExpectedQuestionsAfterUpdated, inmemRepo.questions); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInmemRepositoryDelete(t *testing.T) {
	var (
		question1 = Question{ID: 1, Question: "How many characters are there in \"Quipper\"?", Answer: "7"}
		question2 = Question{ID: 2, Question: "Guess random number, 1, 2, 3 or 4?", Answer: "4"}
		question3 = Question{ID: 3, Question: "Quipper vs Ruangguru?", Answer: "Quipper"}
	)

	tt := []struct {
		Name                          string
		QuestionID                    int
		QuestionInmemDB               []Question
		ExpectedQuestionsAfterDeleted []Question
		ExpectedErr                   error
	}{
		{
			Name:                          "delete question 2",
			QuestionID:                    2,
			QuestionInmemDB:               []Question{question1, question2, question3},
			ExpectedQuestionsAfterDeleted: []Question{question1, question3},
			ExpectedErr:                   nil,
		},
		{
			Name:                          "deleted question 3, failed not found",
			QuestionID:                    3,
			QuestionInmemDB:               []Question{question1, question2},
			ExpectedQuestionsAfterDeleted: []Question{question1, question2},
			ExpectedErr:                   ErrQuestionNotFound,
		},
	}

	ctx := context.Background()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			var (
				r         = NewRepository()
				inmemRepo = r.(*inmemRepository)
			)
			inmemRepo.questions = tc.QuestionInmemDB

			if err := r.Delete(ctx, tc.QuestionID); !errors.Is(tc.ExpectedErr, err) {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.ExpectedQuestionsAfterDeleted, inmemRepo.questions); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
