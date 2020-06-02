# Choose-Your-Own-Adventure
This is a solution to "cyoa" exercise on [Gophercises](https://courses.calhoun.io/courses/cor_gophercises),
done using test-driven development (TDD).

## Techniques and packages used
- html templates
- unit-testing http servers
- flag package for parsing command line arguments; json package for parsing JSON

To start an example app
```
git clone https://github.com/chammaaomar/golang-tdd.git
cd golang-tdd/cmd/cyoa
go build .
./coya stories
```

This will start the app on `localhost:8080` with a sample choose-your-own-adventure.
