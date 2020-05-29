# Quiz
This is a solution to "Quiz Game" exercise on [Gophercises](https://courses.calhoun.io/courses/cor_gophercises),
done using test-driven development (TDD).

## Techniques and packages used
- unit-testing, mocking using dependency injection and interfaces
- channels and go-routines for timer functionality
- flag package for parsing command line arguments; csv package for parsing CSV

To play the game
```
git clone https://github.com/chammaaomar/golang-tdd.git
cd golang-tdd/cmd/quiz
go build .
./quiz -h
```

This will bring up the help menu for the CLI.

## Notes and Limitations
- Only integer answers are allowed
- No vetting is done of the _answer_ column in the CSV. It's simply taken as a string, not evaluated.
