package main

import (
	"flag"
	"github.com/chammaaomar/golang-tdd/quiz"
)

var timerPtr = flag.Int("timer", 30, "time limit in seconds")
var csvPathPtr = flag.String("questions", "problems.csv", "path to CSV with question/answer pairs")
var headerPtr = flag.Bool("header", false, "whether the questions CSV has a header")

func main() {
	flag.Parse()
	quiz.PlayGame(*csvPathPtr, *timerPtr, *headerPtr)
}
