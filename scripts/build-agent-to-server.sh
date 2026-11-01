#!/bin/bash
set -e
cd ../agent
go build -o ../server/agents/agent_x86_64.bin
