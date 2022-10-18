package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/muktihari/quiz_master/questionnaire"
)

func TestIntegrationScenario1(t *testing.T) {
	// Integration test, the order in table test is important, can't be parallelized.
	tt := []struct {
		Name        string
		In          string
		ExpectedOut string
	}{
		{
			Name:        "show help twice",
			In:          "help\nhelp\nexit",
			ExpectedOut: fmt.Sprintf("$ %s$ %s$ ", HelpText, HelpText),
		},
		{
			Name:        "invalid command",
			In:          "helpx\nexit",
			ExpectedOut: "$ Command \"helpx\" is not found. See \"help\"\n$ ",
		},
		// create
		{
			Name: "create question 1 success",
			In:   "create_question 1 \"How many characters are there in \"Quipper\"?\" 7\nexit",
			ExpectedOut: fmt.Sprintf("$ Question no %d created:\n"+PrintFormat+"$ ",
				1, "How many characters are there in \"Quipper\"?", "7"),
		},
		{
			Name: "create question 2 success",
			In:   "create_question 2 \"How many characters are there in \"Quipperzz\"?\" 9\nexit",
			ExpectedOut: fmt.Sprintf("$ Question no %d created:\n"+PrintFormat+"$ ",
				2, "How many characters are there in \"Quipperzz\"?", "9"),
		},
		{
			Name: "create question 3 success",
			In:   "create_question 3 \"How many characters are there in \"Engineer\"?\" 8\nexit",
			ExpectedOut: fmt.Sprintf("$ Question no %d created:\n"+PrintFormat+"$ ",
				3, "How many characters are there in \"Engineer\"?", "8"),
		},
		{
			Name:        "create question invalid ID",
			In:          "create_question X \"How many characters are there in \"Quipper\"?\" 7\nexit",
			ExpectedOut: "$ Invalid question ID, should be integer\n$ ",
		},
		{
			Name:        "create question invalid command format",
			In:          "create_question 1 \"How many characters are there in \"Quipper\"?\"\nexit",
			ExpectedOut: "$ Invalid input format. See \"help\"\n$ ",
		},
		{
			Name:        "create question failed duplicate",
			In:          "create_question 1 \"How many characters are there in \"Quipper\"?\" 7\nexit",
			ExpectedOut: "$ Could not create question: question is already exist\n$ ",
		},
		// update
		{
			Name: "update question 2 success",
			In:   "update_question 2 \"How many characters are there in 'Quipperx'?\" 8\nexit",
			ExpectedOut: fmt.Sprintf("$ Question no %d updated:\n"+PrintFormat+"$ ",
				2, "How many characters are there in 'Quipperx'?", "8"),
		},
		{
			Name:        "update question invalid ID",
			In:          "update_question X \"How many characters are there in \"Quipper\"?\" 7\nexit",
			ExpectedOut: "$ Invalid question ID, should be integer\n$ ",
		},
		{
			Name:        "update question invalid command format",
			In:          "update_question 1 \"How many characters are there in \"Quipper\"?\"\nexit",
			ExpectedOut: "$ Invalid input format. See \"help\"\n$ ",
		},
		{
			Name: "update question failed not found",
			In:   "update_question 4 \"How many characters are there in \"Quipper\"?\" 7\nexit",
			ExpectedOut: fmt.Sprintf("$ Could not update question [%d]: %v\n$ ",
				4, questionnaire.ErrQuestionNotFound),
		},
		// delete
		{
			Name:        "delete question 2 success",
			In:          "delete_question 2\nexit",
			ExpectedOut: fmt.Sprintf("$ Question no %d deleted:\n$ ", 2),
		},
		{
			Name:        "delete question invalid ID",
			In:          "delete_question X\nexit",
			ExpectedOut: "$ Invalid question ID, should be integer\n$ ",
		},
		{
			Name:        "delete question invalid command format",
			In:          "delete_question 1 \"How many characters are there in \"Quipper\"?\"\nexit",
			ExpectedOut: "$ Invalid input format. See \"help\"\n$ ",
		},
		{
			Name:        "delete question failed not found",
			In:          "delete_question 4\nexit",
			ExpectedOut: fmt.Sprintf("$ Could not delete question [%d]: %v\n$ ", 4, questionnaire.ErrQuestionNotFound),
		},
		// question
		{
			Name:        "question 1 success",
			In:          "question 1\nexit",
			ExpectedOut: "$ Q: \"How many characters are there in \"Quipper\"?\"\nA: 7\n$ ",
		},
		{
			Name:        "question invalid ID",
			In:          "question X\nexit",
			ExpectedOut: "$ Invalid question ID, should be integer\n$ ",
		},
		{
			Name:        "question invalid command format",
			In:          "question 1 \"How many characters are there in \"Quipper\"?\"\nexit",
			ExpectedOut: "$ Invalid input format. See \"help\"\n$ ",
		},
		{
			Name:        "question failed not found",
			In:          "question 4\nexit",
			ExpectedOut: fmt.Sprintf("$ Could not get question [%d]: %v\n$ ", 4, questionnaire.ErrQuestionNotFound),
		},
		// questions
		{
			Name: "questions",
			In:   "questions\nexit",
			ExpectedOut: "$ No | Question | Answer\n" +
				"1 \"How many characters are there in \"Quipper\"?\" 7\n" +
				"3 \"How many characters are there in \"Engineer\"?\" 8\n$ ",
		},
		// answer
		{
			Name:        "answer question 1 correct",
			In:          "answer_question 1 7\nexit",
			ExpectedOut: "$ Correct!\n$ ",
		},
		{
			Name:        "answer question 1 correct using number in word",
			In:          "answer_question 1 seven\nexit",
			ExpectedOut: "$ Correct!\n$ ",
		},
		{
			Name:        "answer question 1 incorrect",
			In:          "answer_question 1 8\nexit",
			ExpectedOut: "$ Incorrect!\n$ ",
		},
		{
			Name:        "answer question invalid ID",
			In:          "answer_question X 1\nexit",
			ExpectedOut: "$ Invalid question ID, should be integer$ ",
		},
		{
			Name:        "answer question invalid command format",
			In:          "answer_question 1\nexit",
			ExpectedOut: "$ Invalid input format. See \"help\"\n$ ",
		},
		{
			Name:        "answer question failed not found",
			In:          "answer_question 4 seven\nexit",
			ExpectedOut: fmt.Sprintf("$ Could not answer question [%d]: %v\n$ ", 4, questionnaire.ErrQuestionNotFound),
		},
	}

	var qs questionnaire.Service
	{
		r := questionnaire.NewRepository()
		qs = questionnaire.NewService(r)
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			var (
				in  = strings.NewReader(tc.In)
				out = new(strings.Builder)
			)
			if run(context.Background(), qs, in, out) != 0 {
				t.Fatalf("do not exit properly\n")
			}
			if diff := cmp.Diff(tc.ExpectedOut, out.String()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
