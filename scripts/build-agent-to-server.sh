#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

case "$(uname -m)" in
	x86_64|amd64)
		GOARCH=amd64
		;;
	aarch64|arm64)
		GOARCH=arm64
		;;
	*)
		echo "Unsupported host architecture: $(uname -m)" >&2
		exit 1
		;;
esac

cd "$ROOT_DIR/agent"
GOOS=linux GOARCH="$GOARCH" go build -o "$ROOT_DIR/server/agents/agent_linux_${GOARCH}.bin" .
