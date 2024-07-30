.PHONY: all test coverage

all: get build install

get:
	go get ./...

build:
	cd cmd/smfimage && config build --versiondir='../../'

release:
	config release --package=smfimage
	cd cmd/smfimage && config build --versiondir='../../'

test:
	go test ./... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

coverage: test
	go tool cover -html=.coverage.txt