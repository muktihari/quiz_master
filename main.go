package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/muktihari/quiz_master/pkg/textinput"
	"github.com/muktihari/quiz_master/questionnaire"
)

// Command is command inside the CLI.
type Command string

const (
	Exit           Command = "exit"
	Help           Command = "help"
	Question       Command = "question"
	Questions      Command = "questions"
	CreateQuestion Command = "create_question"
	UpdateQuestion Command = "update_question"
	DeleteQuestion Command = "delete_question"
	AnswerQuestion Command = "answer_question"

	HelpText = "Command | Description\n" +
		"help | Shows list of available command\n" +
		"create_question <no> <question> <answer> | Create a question\n" +
		"update_question <no> <question> <answer> | Update a question\n" +
		"delete_question <no> | Update a question\n" +
		"question <no> | Shows a question\n" +
		"questions | Shows list of question\n" +
		"exit | Exit CLI\n"

	PrintFormat = "Q: \"%s\"\nA: %s\n"
)

func main() {
	var qs questionnaire.Service
	{
		r := questionnaire.NewRepository()
		qs = questionnaire.NewService(r)
	}

	fmt.Println("Welcome to Quiz Master!")
	os.Exit(run(context.Background(), qs, os.Stdin, os.Stdout))
}

func run(ctx context.Context, qs questionnaire.Service, in io.Reader, out io.Writer) int {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, "$ ")
		scanner.Scan()

		var (
			args = textinput.Split(scanner.Text(), ' ')
			cmd  = Command(strings.ToLower(args[0]))
		)

		switch cmd {
		case Exit:
			return 0
		case Help:
			fmt.Fprint(out, HelpText)
		case Question:
			question(ctx, qs, args, out)
		case Questions:
			questions(ctx, qs, args, out)
		case CreateQuestion:
			createQuestion(ctx, qs, args, out)
		case UpdateQuestion:
			updateQuestion(ctx, qs, args, out)
		case DeleteQuestion:
			deleteQuestion(ctx, qs, args, out)
		case AnswerQuestion:
			answerQuestion(ctx, qs, args, out)
		default:
			fmt.Fprintf(out, "Command \"%s\" is not found. See \"help\"\n", cmd)
		}
	}
}

func question(ctx context.Context, qs questionnaire.Service, args []string, out io.Writer) {
	if len(args) != 2 {
		fmt.Fprintln(out, "Invalid input format. See \"help\"")
		return
	}

	id, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Fprintln(out, "Invalid question ID, should be integer")
		return
	}

	question, err := qs.GetByID(ctx, int(id))
	if err != nil {
		fmt.Fprintf(out, "Could not get question [%d]: %v\n", id, err)
		return
	}

	fmt.Fprintf(out, PrintFormat, question.Question, question.Answer)
}

func questions(ctx context.Context, qs questionnaire.Service, _ []string, out io.Writer) {
	questions, err := qs.GetAll(ctx)
	if err != nil {
		fmt.Fprintln(out, err)
		return
	}
	fmt.Fprintln(out, "No | Question | Answer")

	for _, question := range questions {
		fmt.Fprintf(out, "%d \"%s\" %s\n", question.ID, question.Question, question.Answer)
	}
}

func createQuestion(ctx context.Context, qs questionnaire.Service, args []string, out io.Writer) {
	if len(args) != 4 {
		fmt.Fprintln(out, "Invalid input format. See \"help\"")
		return
	}

	id, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Fprintln(out, "Invalid question ID, should be integer")
		return
	}

	question := &questionnaire.Question{
		ID:       int(id),
		Question: args[2],
		Answer:   args[3],
	}
	if err := qs.Create(ctx, question); err != nil {
		fmt.Fprintf(out, "Could not create question: %v\n", err)
		return
	}

	fmt.Fprintf(out, "Question no %d created:\n", question.ID)
	fmt.Fprintf(out, PrintFormat, question.Question, question.Answer)
}

func updateQuestion(ctx context.Context, qs questionnaire.Service, args []string, out io.Writer) {
	if len(args) != 4 {
		fmt.Fprintf(out, "Invalid input format. See \"help\"\n")
		return
	}

	id, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Fprintln(out, "Invalid question ID, should be integer")
		return
	}

	question := &questionnaire.Question{
		ID:       int(id),
		Question: args[2],
		Answer:   args[3],
	}
	if err := qs.Update(ctx, question); err != nil {
		fmt.Fprintf(out, "Could not update question [%d]: %v\n", id, err)
		return
	}

	fmt.Fprintf(out, "Question no %d updated:\n", question.ID)
	fmt.Fprintf(out, PrintFormat, question.Question, question.Answer)
}

func deleteQuestion(ctx context.Context, qs questionnaire.Service, args []string, out io.Writer) {
	if len(args) != 2 {
		fmt.Fprintln(out, "Invalid input format. See \"help\"")
		return
	}

	id, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Fprintln(out, "Invalid question ID, should be integer")
		return
	}

	if err := qs.Delete(ctx, int(id)); err != nil {
		fmt.Fprintf(out, "Could not delete question [%d]: %v\n", id, err)
		return
	}

	fmt.Fprintf(out, "Question no %d deleted:\n", id)
}

func answerQuestion(ctx context.Context, qs questionnaire.Service, args []string, out io.Writer) {
	if len(args) != 3 {
		fmt.Fprintln(out, "Invalid input format. See \"help\"")
		return
	}

	id, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Fprintf(out, "Invalid question ID, should be integer")
		return
	}

	correct, err := qs.Answer(ctx, int(id), args[2])
	if err != nil {
		fmt.Fprintf(out, "Could not answer question [%d]: %v\n", id, err)
		return
	}

	if correct {
		fmt.Fprintln(out, "Correct!")
	} else {
		fmt.Fprintln(out, "Incorrect!")
	}
}
