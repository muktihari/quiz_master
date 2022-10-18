package questionnaire_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/muktihari/quiz_master/questionnaire"
)

type mockRepository struct {
	questionnaire.Repository
	// injectable funcs to mock methods of questionnaire.Repository interface{}
	getByIDFunc func(ctx context.Context, ID int) (*questionnaire.Question, error)
	getAllFunc  func(ctx context.Context) ([]questionnaire.Question, error)
	createFunc  func(ctx context.Context, question *questionnaire.Question) error
	updateFunc  func(ctx context.Context, question *questionnaire.Question) error
	deleteFunc  func(ctx context.Context, ID int) error
}

func (r *mockRepository) GetByID(ctx context.Context, ID int) (*questionnaire.Question, error) {
	return r.getByIDFunc(ctx, ID)
}

func (r *mockRepository) GetAll(ctx context.Context) ([]questionnaire.Question, error) {
	return r.getAllFunc(ctx)
}

func (r *mockRepository) Create(ctx context.Context, question *questionnaire.Question) error {
	return r.createFunc(ctx, question)
}

func (r *mockRepository) Update(ctx context.Context, question *questionnaire.Question) error {
	return r.updateFunc(ctx, question)
}

func (r *mockRepository) Delete(ctx context.Context, ID int) error {
	return r.deleteFunc(ctx, ID)
}

func TestServiceGetByID(t *testing.T) {
	var predefinedQuestions = []questionnaire.Question{
		{ID: 1, Question: "How many word in \"Quipper\"?", Answer: "7"},
		{ID: 2, Question: "Guess random number, 1, 2, 3 or 4?", Answer: "4"},
	}

	tt := []struct {
		Name             string
		QuestionID       int
		MockRepository   questionnaire.Repository
		ExpectedQuestion *questionnaire.Question
		ExpectedErr      error
	}{
		{
			Name:       "valid ID, question ID 1 found",
			QuestionID: 1,
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{getByIDFunc: func(ctx context.Context, ID int) (*questionnaire.Question, error) {
					return &predefinedQuestions[0], nil
				}}
			}(),
			ExpectedQuestion: &predefinedQuestions[0],
			ExpectedErr:      nil,
		},
		{
			Name:       "question ID 3, not found",
			QuestionID: 2,
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{getByIDFunc: func(ctx context.Context, ID int) (*questionnaire.Question, error) {
					return nil, questionnaire.ErrQuestionNotFound
				}}
			}(),
			ExpectedQuestion: nil,
			ExpectedErr:      questionnaire.ErrQuestionNotFound,
		},
	}

	ctx := context.Background()
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			qs := questionnaire.NewService(tc.MockRepository)
			question, err := qs.GetByID(ctx, tc.QuestionID)
			if tc.ExpectedErr != err {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.ExpectedQuestion, question); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestServiceGetAll(t *testing.T) {
	var predefinedQuestions = []questionnaire.Question{
		{ID: 1, Question: "How many word in \"Quipper\"?", Answer: "7"},
		{ID: 2, Question: "Guess random number, 1, 2, 3 or 4?", Answer: "4"},
	}
	tt := []struct {
		Name              string
		MockRepository    questionnaire.Repository
		ExpectedQuestions []questionnaire.Question
		ExpectedErr       error
	}{
		{
			Name: "all question retrieved",
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{getAllFunc: func(ctx context.Context) ([]questionnaire.Question, error) {
					return predefinedQuestions, nil
				}}
			}(),
			ExpectedQuestions: predefinedQuestions,
			ExpectedErr:       nil,
		},
		{
			Name: "no question retrieved",
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{getAllFunc: func(ctx context.Context) ([]questionnaire.Question, error) {
					return []questionnaire.Question{}, nil
				}}
			}(),
			ExpectedQuestions: []questionnaire.Question{},
			ExpectedErr:       nil,
		},
	}

	ctx := context.Background()
	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			qs := questionnaire.NewService(tc.MockRepository)
			questions, err := qs.GetAll(ctx)
			if tc.ExpectedErr != err {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.ExpectedQuestions, questions); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestServiceCreate(t *testing.T) {
	var predefinedQuestions = []questionnaire.Question{
		{ID: 1, Question: "How many word in \"Quipper\"?", Answer: "7"},
		{ID: 2, Question: "Guess random number, 1, 2, 3 or 4?", Answer: "4"},
	}
	tt := []struct {
		Name           string
		Question       *questionnaire.Question
		MockRepository questionnaire.Repository
		ExpectedErr    error
	}{
		{
			Name:     "create question 1 success",
			Question: &predefinedQuestions[0],
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{createFunc: func(ctx context.Context, question *questionnaire.Question) error {
					return nil
				}}
			}(),
			ExpectedErr: nil,
		},
		{
			Name:     "create question 1 failed already exist",
			Question: &predefinedQuestions[0],
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{createFunc: func(ctx context.Context, question *questionnaire.Question) error {
					return questionnaire.ErrQuestionIsAlreadyExist
				}}
			}(),
			ExpectedErr: questionnaire.ErrQuestionIsAlreadyExist,
		},
	}

	ctx := context.Background()
	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			qs := questionnaire.NewService(tc.MockRepository)
			if err := qs.Create(ctx, tc.Question); tc.ExpectedErr != err {
				t.Fatal(err)
			}
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	var predefinedQuestions = []questionnaire.Question{
		{ID: 1, Question: "How many word in \"Quipper\"?", Answer: "7"},
		{ID: 2, Question: "Guess random number, 1, 2, 3 or 4?", Answer: "4"},
	}
	tt := []struct {
		Name           string
		Question       *questionnaire.Question
		MockRepository questionnaire.Repository
		ExpectedErr    error
	}{
		{
			Name:     "update question 1 success",
			Question: &predefinedQuestions[0],
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{updateFunc: func(ctx context.Context, question *questionnaire.Question) error {
					return nil
				}}
			}(),
			ExpectedErr: nil,
		},
		{
			Name:     "update question 1 failed not found",
			Question: &predefinedQuestions[0],
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{updateFunc: func(ctx context.Context, question *questionnaire.Question) error {
					return questionnaire.ErrQuestionNotFound
				}}
			}(),
			ExpectedErr: questionnaire.ErrQuestionNotFound,
		},
	}

	ctx := context.Background()
	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			qs := questionnaire.NewService(tc.MockRepository)
			if err := qs.Update(ctx, tc.Question); tc.ExpectedErr != err {
				t.Fatal(err)
			}
		})
	}
}

func TestServiceDelete(t *testing.T) {
	tt := []struct {
		Name           string
		QuestionID     int
		MockRepository questionnaire.Repository
		ExpectedErr    error
	}{
		{
			Name:       "delete question 1 success",
			QuestionID: 1,
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{deleteFunc: func(ctx context.Context, ID int) error {
					return nil
				}}
			}(),
			ExpectedErr: nil,
		},
		{
			Name:       "delete question 1 failed not found",
			QuestionID: 1,
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{deleteFunc: func(ctx context.Context, ID int) error {
					return questionnaire.ErrQuestionNotFound
				}}
			}(),
			ExpectedErr: questionnaire.ErrQuestionNotFound,
		},
	}

	ctx := context.Background()
	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			qs := questionnaire.NewService(tc.MockRepository)
			if err := qs.Delete(ctx, tc.QuestionID); tc.ExpectedErr != err {
				t.Fatal(err)
			}
		})
	}
}

func TestServiceAnswer(t *testing.T) {
	var predefinedQuestions = []questionnaire.Question{
		{ID: 1, Question: "How many word in \"Quipper\"?", Answer: "7"},
		{ID: 2, Question: "Guess random number, 1, 2, 3 or 4?", Answer: "4"},
	}
	tt := []struct {
		Name                string
		QuestionID          int
		MockRepository      questionnaire.Repository
		Answer              string
		ExpectedCorrectness bool
		ExpectedErr         error
	}{
		{
			Name:       "answer question 1 with answer \"7\", correct",
			QuestionID: 1,
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{getByIDFunc: func(ctx context.Context, ID int) (*questionnaire.Question, error) {
					return &predefinedQuestions[0], nil
				}}
			}(),
			Answer:              "7",
			ExpectedCorrectness: true,
			ExpectedErr:         nil,
		},
		{
			Name:       "answer question 1 with answer \"seven\", correct",
			QuestionID: 1,
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{getByIDFunc: func(ctx context.Context, ID int) (*questionnaire.Question, error) {
					return &predefinedQuestions[0], nil
				}}
			}(),
			Answer:              "seven",
			ExpectedCorrectness: true,
			ExpectedErr:         nil,
		},
		{
			Name:       "anwser question 1 incorrect",
			QuestionID: 1,
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{getByIDFunc: func(ctx context.Context, ID int) (*questionnaire.Question, error) {
					return &predefinedQuestions[0], nil
				}}
			}(),
			Answer:              "8",
			ExpectedCorrectness: false,
			ExpectedErr:         nil,
		},
		{
			Name:       "anwser question 3 not found",
			QuestionID: 1,
			MockRepository: func() questionnaire.Repository {
				return &mockRepository{getByIDFunc: func(ctx context.Context, ID int) (*questionnaire.Question, error) {
					return nil, questionnaire.ErrQuestionNotFound
				}}
			}(),
			Answer:              "8",
			ExpectedCorrectness: false,
			ExpectedErr:         questionnaire.ErrQuestionNotFound,
		},
	}

	ctx := context.Background()
	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			qs := questionnaire.NewService(tc.MockRepository)
			correct, err := qs.Answer(ctx, tc.QuestionID, tc.Answer)
			if tc.ExpectedErr != err {
				t.Fatal(err)
			}
			if tc.ExpectedCorrectness != correct {
				t.Fatalf("expected %v, got: %v", tc.ExpectedCorrectness, correct)
			}
		})
	}
}
