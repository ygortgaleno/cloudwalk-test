BINARY_NAME=parse_quake_log
BYNARY_FOLDER=bin

all: build

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BYNARY_FOLDER}/${BINARY_NAME}_darwin cmd/parse_quake_log.go
	GOARCH=arm64 GOOS=darwin go build -o ${BYNARY_FOLDER}/${BINARY_NAME}_m1 cmd/parse_quake_log.go
	GOARCH=amd64 GOOS=linux go build -o ${BYNARY_FOLDER}/${BINARY_NAME}_linux cmd/parse_quake_log.go
	GOARCH=amd64 GOOS=windows go build -o ${BYNARY_FOLDER}/${BINARY_NAME}_windows cmd/parse_quake_log.go

clean_builds:
	rm -rf ${BYNARY_FOLDER}/*

test:
	go test ./... -count=1

test_coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: all build clean_builds test test_coverage