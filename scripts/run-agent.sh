#!/usr/bin/env bash

export LOG_LEVEL=debug
set -euo pipefail
PORT="${PORT:-9000}"
BIN="/tmp/servergram-agent"

command -v socat >/dev/null || { echo "socat not found: sudo apt install socat" >&2; exit 1; }

echo "Agent on 127.0.0.1:$PORT (plaintext, loopback only)"
exec socat "TCP-LISTEN:$PORT,bind=127.0.0.1,reuseaddr,fork" "EXEC:$BIN serve --stdio"
