//go:build !js
// +build !js

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	//DownloadWasmExec()
	fileServer := http.FileServer(http.Dir("."))
	http.Handle("/", fileServer)
	println("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

//const wasmExecURL = "https://raw.githubusercontent.com/golang/go/release-branch.go1.12/misc/wasm/wasm_exec.js"
const wasmExecURL = "https://raw.githubusercontent.com/golang/go/master/misc/wasm/wasm_exec.js"
const wasmExecFile = "wasm_exec.js"

func DownloadWasmExec() {
	if _, err := os.Stat(wasmExecFile); err == nil {
		return
	}
	println("Downloading wasm_exec.js...")
	out, err := os.Create(wasmExecFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	resp, err := http.Get(wasmExecURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
