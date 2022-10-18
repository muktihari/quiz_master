# Quiz Master

Quiz Master is a interactive CLI app that you can play along with it.

## Specifications
### Should have following commands
- help: Shows list of available command
  - e.g.: ```sh help```
- create_question: Create a question, return error if duplicate
  - e.g.: ```sh create_question 1 "How many words is 'Quipper'?" 7```
- update_question: Update a question, return error if not found
  - e.g.: ```sh update_question 1 "How many words is 'TQIF'?" 4```
- delete_question: Delete a question, return error if not found
  - e.g.: ```sh delete_question 1```
- question: Shows a question, return error if not found
  - e.g.: ```sh question 1```
- questions: Shows all questions
  - e.g.: ```sh questions```
- answer_question: Answer a question, it will return "Correct!" or "Incorrect!"
  - e.g.: ```sh answer_question 1 7```
- exit: Exit Quiz Master CLI
  - e.g.: ```sh exit```

### Should recognize numbers
If the answer in/contains a number, it should recognize the number
```
Q : How many vowels are there in the English alphabet?
A : 5

Answer : 5 is correct
Answer : five is correct
Answer : Five is correct
Answer : 6 is incorrect
Answer : six is incorrect
```
I assume that following question-answer situasion is also work the same way
```
Q : How many vowels are there in the English alphabet?
A : five

Answer : 5 is correct
Answer : five is correct
```