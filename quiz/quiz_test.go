package quiz

import (
	"bytes"
	"encoding/csv"
	"os"
	"path"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const testDir = "test"

type spySleeper struct {
	args []time.Duration
}

type spyPrinter struct {
	called int
}

func (s *spyPrinter) Println(a ...interface{}) (int, error) {
	s.called++
	return 1, nil
}

func (s *spySleeper) Sleep(d time.Duration) {
	s.args = append(s.args, d)
}

func TestParseCSV(t *testing.T) {
	t.Run("CSV with more than two columns should be gracefully rejected", func(t *testing.T) {
		_, errParse := setupParseCSV("3_column.csv", false)
		if errParse != errBadColumns {
			t.Fatalf("Expected error %v, got error %v", errBadColumns, errParse)
		}
	})

	t.Run("CSV with fewer than two columns should be gracefully rejected", func(t *testing.T) {
		_, errParse := setupParseCSV("1_column.csv", false)
		if errParse != errBadColumns {
			t.Fatalf("Expected error %v, got error %v", errBadColumns, errParse)
		}
	})

	t.Run("CSV with non-integer answers should be gracefully rejected", func(t *testing.T) {
		_, errParse := setupParseCSV("non_int.csv", false)
		if errParse != strconv.ErrSyntax {
			t.Fatalf("Expected error %v, got error %v", strconv.ErrSyntax, errParse)
		}
	})

	t.Run("CSV with header should be accepted", func(t *testing.T) {
		_, errParse := setupParseCSV("header.csv", true)
		if errParse != nil {
			t.Fatalf("Expected no error, got %v", errParse)
		}
	})

	t.Run("Acceptable CSV should be correctly parsed", func(t *testing.T) {
		qaMap, err := setupParseCSV("correct.csv", false)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expectedMap := map[string]int{
			"2+5":                       7,
			"What does 3+9 equal, sir?": 12,
		}

		if !reflect.DeepEqual(qaMap, expectedMap) {
			t.Fatalf("Expected parsed CSV to be %v, got %v", expectedMap, qaMap)
		}
	})
}

func TestPlayGame(t *testing.T) {
	t.Run("Basic game loop of pose question then accept answer then update score then pose next question, should work", func(t *testing.T) {
		qaMap := map[string]int{
			"1+4":  5,
			"10/5": 2,
			"5*6":  30,
		}
		var score int
		done := make(chan int, 1)
		outSpy := &spyPrinter{}
		// real := realPrinter{}
		userResponse := bytes.NewBufferString("5\n3\nq\n")
		gameLoop(qaMap, userResponse, outSpy, &score, done)

		expectedResponses := 3

		// testing only the number of responses since the content of the response is implementation / likely to change
		// but the basic fact of the game responding shouldn't change
		if outSpy.called < expectedResponses {
			t.Fatalf("Expected %d responses, instead got %d", expectedResponses, outSpy.called)
		}
	})

	t.Run("Game should exit after timer has run out and show user final score", func(t *testing.T) {
		sleepySpy := &spySleeper{args: make([]time.Duration, 0, 5)}
		printingSpy := &spyPrinter{}
		userResponse := bytes.NewBufferString("\n")
		playGame(path.Join(testDir, "correct.csv"), 30, false, userResponse, sleepySpy, printingSpy)
		expectedSleep := time.Duration(30) * time.Second
		if len(sleepySpy.args) != 1 || sleepySpy.args[0] != expectedSleep {
			t.Fatalf("time.Sleep got called with args %v", sleepySpy.args)
		}
		if printingSpy.called < 2 {
			t.Fatalf("Expected Println to be called at least twice, got called 0 times")
		}
	})
}

func setupParseCSV(filename string, header bool) (map[string]int, error) {
	csvPath := path.Join(testDir, filename)
	csvFile, errOpen := os.Open(csvPath)
	if errOpen != nil {
		return map[string]int{}, errOpen
	}
	reader := csv.NewReader(csvFile)

	return parseCSV(reader, header)
}
