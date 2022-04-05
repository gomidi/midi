
build:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o main.wasm main_js.go

test:
	go test ./... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

coverage: test
	go tool cover -html=.coverage.txt
