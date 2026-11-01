#!/bin/bash
set -e
cd ../agent
go build -o ../server/agents/agent_linux_amd64.bin
