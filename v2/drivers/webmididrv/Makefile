
build:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o main.wasm main_js.go

test:
	GOOS=js GOARCH=wasm WASM_HEADLESS=off go test -exec="/home/benny/go/bin/wasmbrowsertest"
