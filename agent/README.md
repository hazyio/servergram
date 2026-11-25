# Servergram agent

The agent is a small Go service that runs on each managed Linux server. It performs file operations, process control, resource sampling, and cron job management on behalf of the Servergram server.

The server uploads the agent to `~/.servergram/agent` over SSH and starts it. You normally do not install or run the agent by hand.

## Commands

```bash
agent serve
agent run-job -id=<id> [-manual] "command"
```

- `serve` starts the gRPC server over standard input and output. The server starts this command inside an SSH session.
- `run-job` executes a cron job command, records the run in the local database, and marks it as scheduled or manually triggered. Managed crontab entries call this command.

## How it runs

- The agent listens only on the SSH session's standard input and output. It does not open a TCP port.
- It runs as the SSH user, so its permissions are that user's permissions.
- It stores its own SQLite database at `~/.servergram/servergram.db` and migrates it at startup.
- It reads and writes the user's crontab with `crontab`. Managed entries are prefixed so the agent can recognize them.
- It samples system resources with `gopsutil`.
- Logs are written to standard error as JSON. The server forwards them to its own log output.

## Environment

| Variable    | Default | Description                    |
| ----------- | ------- | ------------------------------ |
| `LOG_LEVEL` | `info`  | Set to `debug` for verbose logs. |

## Build

From the repository root:

```bash
cd agent
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../server/agents/agent_linux_amd64.bin .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ../server/agents/agent_linux_arm64.bin .
```

The server embeds these two files at compile time. To build only the current machine's architecture, run `scripts/build-agent-to-server.sh`.

## Development

Build a local binary:

```bash
go build -o /tmp/servergram-agent .
```

`scripts/run-agent.sh` exposes the agent on `127.0.0.1:9000` through `socat` for transport debugging. It requires `socat` and does not replace a real SSH session.

The agent shares protobuf definitions with the server through the `shared` module. From the `scripts/` directory, run `./generate-protobuf.sh` to regenerate them.

## Platform support

The runtime targets Linux on `amd64` and `arm64`. Commands such as `crontab` and the resource APIs are Linux-specific.
