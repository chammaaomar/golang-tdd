package quiz

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type printer interface {
	Println(a ...interface{}) (int, error)
}

type realPrinter struct{}

func (r *realPrinter) Println(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

var errBadColumns = errors.New("CSV file has a record with the wrong number of columns, expected two")
var greetingMessage = "Welcome to the maths quiz! Press any button to continue, or 'q' at any time to exit"
var byeMessage = "Thank you for playing, your final score is "
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

// PlayGame controls the basic loop of the quiz: Greets, then presents
// questions, updates score, and shows final score to user
func PlayGame(qaMap map[string]int, input io.Reader, output printer) {
	scanner := bufio.NewScanner(input)
	score := 0
	output.Println(greetingMessage)

	scanner.Scan()
	userInput := strings.TrimSpace(scanner.Text())
	if userInput == endGame {
		output.Println(byeMessage, score)
		return
	}

	for question, answer := range qaMap {
		output.Println(question)
		scanner.Scan()
		userInput := scanner.Text()
		if userInput == endGame {
			output.Println(byeMessage, score)
			return
		}
		userAnswer, err := strconv.Atoi(strings.TrimSpace(userInput))
		if err == nil && userAnswer == answer {
			score++
		}
	}
	output.Println(byeMessage, score)
	return
}
