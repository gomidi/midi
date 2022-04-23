all: get build install

get:
	go get ./...

build:
	config build --plattforms='linux/amd64 windows/386/cgo' 

release:
	config release
	config build --plattforms='linux/amd64 windows/386/cgo'

test:
	go test ./... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

coverage: test
	go tool cover -html=.coverage.txt
