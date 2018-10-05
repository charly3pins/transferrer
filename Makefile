.PHONY: all docker_test docker_coverage deps build

BIN_NAME=transferrer

all: test docker_test docker_coverage deps build

test:
	mkdir -p cover
	go test  -coverprofile=cover/coverage.out
	go tool cover -html=cover/coverage.out -o cover/index.html

docker_test: 
	docker run --rm -v "${PWD}:/go/src/github.com/charly3pins/transferrer/"  -w "/go/src/github.com/charly3pins/transferrer/" golang:1.10 make deps test

docker_coverage:
	make docker_test
	docker run --rm -v "${PWD}/cover/.:/usr/share/nginx/html" -p 8081:80 nginx

deps:
	go get github.com/charly3pins/transferrer
	go get github.com/jmoiron/sqlx
	go get github.com/dgrijalva/jwt-go
	go get github.com/gin-gonic/gin
	go get github.com/gin-contrib/cors
	go get github.com/lib/pq

build:
	cd cmd/server; go build -o ${BIN_NAME}
	@echo "You can now use ./${BIN_NAME}"
