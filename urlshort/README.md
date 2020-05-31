# URLshort
This is a solution to "URL Shortener" exercise on [Gophercises](https://courses.calhoun.io/courses/cor_gophercises),
done using test-driven development (TDD).

## Techniques and packages used
- unit-testing HTTP servers
- YAML and JSON decoding
- BoltDB embedded (single file) key-value store 

## Usage
This is meant to be provided as a package. Its public API consists of four functions

- MapHandler
- YAMLHandler
- JSONHandler
- BoltDBHandler

An example application that uses them is provided in `/cmd/urlshort/main.go`.
