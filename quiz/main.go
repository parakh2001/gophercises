package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// input is csv file wiht questions and answers
// and time limit from command line
// output-> questions displayed one by one
// user enters answers
// final score
// quiz ends after timer expires
// need to count correct answers
type Quiz struct {
	Question string
	Answer   string
}

func ReadCsv(reader *csv.Reader, quizCh chan<- Quiz) {
	defer close(quizCh)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return
		}
		if err != nil {
			continue
		}
		if len(record) < 2 {
			continue
		}
		quizCh <- Quiz{
			Question: record[0],
			Answer:   record[1],
		}
	}
}

// / This function should run independently and shouldn't be blocked by users quiz
// / this function should return to main and tell time's up.
func StartTimer(eventCh chan bool) {
	time.Sleep(time.Second * 30)
	eventCh <- true
}
func AskQuestion(question Quiz, index int) string {
	var answer string
	fmt.Printf("Question %d:\n", index)
	fmt.Printf("%s = ", question.Question)
	fmt.Scanln(&answer)
	return answer
}
func checkAnswer(answer Quiz, userAnswer string) bool {
	if answer.Answer == userAnswer {
		return true
	}
	return false
}
func PrintScore(score int) {
	fmt.Printf("Final score is : %d", score)
}
func main() {
	var score int
	// load the csv file so taht ReadCsv function can read that file and other function can extract the question and answer
	// open the csv file
	file, err := os.Open("problems.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()
	// initialise the csv reader
	reader := csv.NewReader(file)
	quizCh := make(chan Quiz)
	go ReadCsv(reader, quizCh)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()
	questionNumber := 1
	for quiz := range quizCh {
		answerCh := make(chan string)
		go func(q Quiz, index int) {
			answerCh <- AskQuestion(q, index)
		}(quiz, questionNumber)
		select {
		case userAnswer := <-answerCh:
			if checkAnswer(quiz, userAnswer) {
				score++
			}
		case <-ctx.Done():
			fmt.Println("\nTime's up!")
			PrintScore(score)
			return
		}
		questionNumber++
	}
	PrintScore(score)
}
