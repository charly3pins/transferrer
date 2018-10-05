# transferrer
[![Go Report Card](https://goreportcard.com/badge/github.com/charly3pins/transferrer)](https://goreportcard.com/report/github.com/charly3pins/transferrer) [![GoDoc](https://godoc.org/github.com/charly3pins/transferrer?status.svg)](https://godoc.org/github.com/charly3pins/transferrer)

## Service
- An endpoint to get the current balance of a given user.
- An endpoint to make a transfer of money between two users.

## How it works
First of all you need to clone the project:
```
go get github.com/charly3pins/transferrer
cd $GOPATH/src/github.com/charly3pins/transferrer
```
The translation service is accessible via 2 ways. A command-line servive and the API interface. 

## Run the API
Run the service via API it using `docker-compose`. Simply type the following command:
```
docker-compose up
```

This command will create two containers, one with a PostgreSQL database with dummy data on it and other with the API ready to be called.

## Obtain a JWT for work
The API is securized by JWT, so you need to negotiate initially the token. This is a silly endpoint because don't validate the user/pass againt's DB, but for mock the operation of obtaining the token it's enough. Make a POST to `/token` with `email` and `password` fields in the body.
```
curl -d '{"email":"carles@email.com", "password":"myPass"}' -H "Content-Type: application/json" -X POST http://localhost:8080/token
```

## Balance of a user 
To obtain the balance of the user logged in (via jwt) yo need to make a GET to `/balance` endpoint (The token is the obtained in the previous call)
```
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Mzg2OTEyNDAsImVtYWlsIjoiY2FybGVzQGVtYWlsLmNvbSIsInBhc3N3b3JkIjoibXlQYXNzIn0.u2IA1zHvygXC7QG49MevEdDAyllpxzhdM-Mwx8X-K9Q" -X GET http://localhost:8080/balance
```

## Transfer money 
To transfer money between two users you need to make a POST to `/transfer` with the jwt in the header.
The object expected in the body is like:
```
{
	"originUser": "carles@email.com",
	"originNumber": "XXXYYYZZZ",
	"destinationUser": "john@email.com",
	"destinationNumber": "AAABBBCCC",
	"amount": 1000
}
```

Using the JWT obtained in the first call the transfer call should look like:
```
curl -d '{"originUser":"carles@email.com", "originNumber":"XXXYYYZZZ", "destinationUser":"john@email.com", "destinationNumber":"AAABBBCCC", "amount":1000}' -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Mzg2OTI3MDEsImVtYWlsIjoiY2FybGVzQGVtYWlsLmNvbSIsInBhc3N3b3JkIjoibXlQYXNzIn0.zyj5Tx2CEbR_UtHUE79ptWO1xiEnP2ZJTq8nA_9mKd8" -X POST http://localhost:8080/transfer
```

## Test & Coverage
Taking advantage of the tools that `Go` provides for testing, the project also contains a dummy server with nginx that runs the `go test` and `go tool cover` commands to generate the html output with the corresponding coverage.
```
make docker_coverage
```
Once all is ready visit `http://localhost:8081/`
