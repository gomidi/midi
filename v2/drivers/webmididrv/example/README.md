# Example for webmididrv

Build the `main.wasm` file with 

```
GOOS=js GOARCH=wasm go build -o main.wasm main_js.go
```

Start the webserver with 

```
go run main.go
```

And then point your browser to `http://localhost:8080`.

If you need a fresh `wasm_exec.js`, you can get it with

```
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

If you create your own html file, make sure, it contains the lines


```html
<script src="wasm_exec.js"></script>
<script>
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
    });
</script>
```
