create_wasm:
    cp $(go env GOROOT)/lib/wasm/wasm_exec.js _web/examples/.
    env GOOS=js GOARCH=wasm go build -o _web/examples/walking.wasm ./examples/walking
    env GOOS=js GOARCH=wasm go build -o _web/examples/display_banner.wasm ./examples/display_banner
    env GOOS=js GOARCH=wasm go build -o _web/examples/transition.wasm ./examples/transition
