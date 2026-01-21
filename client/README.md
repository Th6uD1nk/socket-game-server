# RTGS Client

## Run the client

```bash
go run rtgs-client
````

## Build the client

```bash
go build rtgs-client
```

### Alternatively, build for WebAssembly (WASM)

```bash
GOOS=js GOARCH=wasm go build -o game.wasm
```

Make sure you copied `wasm_exec.js`:

```bash
chmod +x ./cp_wasm_exec.sh
./cp_wasm_exec.sh
```

Then start a local HTTP server:

```bash
python3 -m http.server 8080
```

Open in your browser via:

```
http://localhost:8080/index.html
```
