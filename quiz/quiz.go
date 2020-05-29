package quiz

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// defined for mocking and dependency injection
type printer interface {
	Println(a ...interface{}) (int, error)
}

type realPrinter struct{}

func (r *realPrinter) Println(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

// defined for mocking and dependency injection
type sleeper interface {
	Sleep(d time.Duration)
}

type realSleeper struct{}

func (s *realSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}

var errBadColumns = errors.New("CSV file has a record with the wrong number of columns, expected two")
var greetingMessage = "Welcome to the maths quiz! Press any button to continue, or enter 'q' at any time to exit"
var byeMessage = "Thank you for playing. your final score is"
var timeOutMessage = "You ran out of time. Thank you for playing. Your final score is"
var outOf = "out of"
var endGame = "q"

func parseCSV(reader *csv.Reader, header bool) (map[string]int, error) {
	parsedCSV := make(map[string]int)

	if header {
		// skip header
		reader.Read()
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return parsedCSV, err
		}
		if len(record) != 2 {
			return parsedCSV, errBadColumns
		}
		errExtract := extractQA(record, parsedCSV)
		if errExtract != nil {
			return parsedCSV, errExtract
		}
	}

	return parsedCSV, nil
}

func extractQA(record []string, csvMap map[string]int) error {
	question, answer := record[0], record[1]
	intAnswer, errAtoi := strconv.Atoi(strings.TrimSpace(answer))

	if errAtoi != nil {
		return strconv.ErrSyntax
	}

	csvMap[question] = intAnswer

	return nil
}

// gameLoop controls the basic loop of the quiz: Pose question,
// check answer, update score, and post next question
func gameLoop(qaMap map[string]int, input io.Reader, output printer, scorePtr *int, done chan int) {
	scanner := bufio.NewScanner(input)
	for question, answer := range qaMap {
		output.Println(question)
		scanner.Scan()
		userInput := scanner.Text()
		if userInput == endGame {
			done <- 1
			return
		}
		userAnswer, err := strconv.Atoi(strings.TrimSpace(userInput))
		if err == nil && userAnswer == answer {
			*scorePtr++
		}
	}
	done <- 1
	return
}

func playGame(csvPath string, timer int, header bool, input io.Reader, sleepy sleeper, output printer) (int, error) {
	var score int
	done := make(chan int)
	quit := make(chan int)
	// set up question-answer CSV
	csvFile, errOpen := os.Open(csvPath)
	if errOpen != nil {
		return score, errOpen
	}

	// parse question-answer CSV to produce a Map of QA
	reader := csv.NewReader(csvFile)
	qaMap, errParse := parseCSV(reader, header)
	if errParse != nil {
		return score, errParse
	}
	maxScore := len(qaMap)

	// greet and wait for user input to start game
	scanner := bufio.NewScanner(input)
	output.Println(greetingMessage)
	scanner.Scan()
	userInput := scanner.Text()
	if userInput == endGame {
		output.Println(byeMessage, score)
		return score, nil
	}

	go gameLoop(qaMap, input, output, &score, done)
	go func() {
		sleepy.Sleep(time.Duration(timer) * time.Second)
		quit <- 1
	}()

	select {
	case <-done:
		output.Println(byeMessage, score, outOf, maxScore)
		return score, nil
	case <-quit:
		output.Println(timeOutMessage, score, outOf, maxScore)
		return score, nil
	}
}

// PlayGame reads the csv at csvPath for question/answer pairs,
// skipping the header if there is one, and plays the game
// for a maximum of timer seconds
func PlayGame(csvPath string, timer int, header bool) (int, error) {
	return playGame(csvPath, timer, header, os.Stdin, &realSleeper{}, &realPrinter{})
}
