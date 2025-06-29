create_wasm:
    cp $(go env GOROOT)/lib/wasm/wasm_exec.js _web/examples/.
    env GOOS=js GOARCH=wasm go build -o _web/examples/walking.wasm ./examples/walking.go