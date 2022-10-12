# gRPC Quote Server

A Streaming gRPC Quote server as reverse proxy for [Quote Server](https://stoicquotesapi.com/).

** Deploy
* Clone the repository
* run `go mod tidy`
* `make run`
Server will be available on `50056`

Two definitions are defined in this project.
* Get Random quotes based on the number you provide.
* Get Random quotes from the author you provide (If the main server supports it)


# Improvement

There's no test in this system, so I should get to that sometimes soon but not right now :smile