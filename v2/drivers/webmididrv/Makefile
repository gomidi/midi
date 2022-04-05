
build:
	GOOS=js GOARCH=wasm go build

test:
	go test ./... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

coverage: test
	go tool cover -html=.coverage.txt
