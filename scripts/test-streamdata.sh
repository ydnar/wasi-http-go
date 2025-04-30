#!/bin/bash
set -euo pipefail

echo "Building streamdata example..."
tinygo build -target=wasip2-roundtrip.json -o streamdata.wasm ./examples/streamdata

echo "Running streamdata example..."
output=$(wasmtime run -Shttp -Sinherit-network -Sinherit-env streamdata.wasm)

# Verify streamdata request worked
if ! echo "$output" | grep -q "Status: 200"; then
	echo "ERROR: streamdata test failed"
	echo "$output"
	exit 1
fi

echo "streamdata test passed!"
