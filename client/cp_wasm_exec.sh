#!/bin/bash
GO_VER=$(go version | awk '{print $3}' | sed 's/go//')
curl -o ./wasm_exec.js https://raw.githubusercontent.com/golang/go/go$GO_VER/misc/wasm/wasm_exec.js
echo "wasm_exec.js (version $GO_VER) downloaded to ./"
