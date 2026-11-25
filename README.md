# Server Gram

Server Gram is a self-hosted web app for managing remote Linux servers. It handles cron jobs, files, running processes, and system resources from a browser. You run one server binary; agents are pushed to managed servers over SSH automatically.

## Repository layout

| Path       | Description                                                                 |
| ---------- | --------------------------------------------------------------------------- |
| `agent/`   | Go agent that runs on managed Linux servers.                                |
| `server/`  | Go HTTPS server, REST API, WebSocket, SSH manager, and embedded web client. |
| `client/`  | SvelteKit web UI, built into `server/dist` and embedded in the server.      |
| `shared/`  | Protobuf definitions and shared Go code used by the agent and server.       |
| `docs/`    | Mintlify documentation site.                                                |
| `scripts/` | Build and code-generation helpers.                                          |

## How it works

1. The server serves the web UI and API over HTTPS. It stores servers, SSH keys, groups, logs, notifications, and sessions in SQLite, PostgreSQL, MySQL, or SQL Server.
2. When you add a server, the server connects over SSH, detects the remote OS and architecture, and uploads the matching agent to `~/.servergram/agent`. The upload is skipped when the checksum matches.
3. The server starts `agent serve` inside the SSH session. The agent exposes gRPC over the session's standard input and output, so no port is opened on the managed server.
4. The UI uses the REST API, a WebSocket for live updates, and server-sent events for process streams. The server translates those calls into gRPC requests to the agent.
5. The agent keeps its own SQLite database at `~/.servergram/servergram.db` and performs the actual work: crontab changes, job runs, file operations, process control, and resource sampling.

The server binary embeds both the web client (`server/dist`) and the Linux agents (`server/agents/agent_linux_amd64.bin`, `server/agents/agent_linux_arm64.bin`) with `go:embed`. Both directories are build outputs and are gitignored, so they must be built before the server.

## Requirements

- Go 1.27.3 or newer
- Bun (recommended) or Node.js 20 or newer
- [mkcert](https://github.com/FiloSottile/mkcert) for locally trusted TLS certificates
- SSH access and a private key for each managed server

Managed servers must run Linux on `amd64` or `arm64`.

## Build

The Vite config reads `certs/localhost+2.pem` and `certs/localhost+2-key.pem` even for production builds, so create the certificates first:

```bash
mkdir -p certs
mkcert -install
mkcert \
  -cert-file certs/localhost+2.pem \
  -key-file certs/localhost+2-key.pem \
  localhost 127.0.0.1 ::1
```

Build the web client into `server/dist`:

```bash
cd client
bun install
bun run build
cd ..
```

Build both agents into `server/agents`:

```bash
cd agent
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../server/agents/agent_linux_amd64.bin .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ../server/agents/agent_linux_arm64.bin .
cd ..
```

Build the server:

```bash
cd server
go build -o servergram-server .
```

`scripts/build-agent-to-server.sh` builds only the current machine's architecture. Build both architectures explicitly for releases.

## Run

Create `server/.env`:

```bash
SERVER_PORT="9788"
ALLOWED_ORIGINS="https://localhost:9788"
PASSWORD="ChangeMe-Str0ng!2026"
SHH_MASTER_KEY="ChangeMe-Too-Str0ng!2026"
SSL_PUBLIC_KEY="../certs/localhost+2.pem"
SSL_PRIVATE_KEY="../certs/localhost+2-key.pem"
```

`DATABASE_URL` is optional and defaults to SQLite in the user config directory.

```bash
cd server
./servergram-server
```

Open `https://localhost:9788`. If `PASSWORD` is empty, authentication is disabled. Changing `SHH_MASTER_KEY` later makes saved SSH keys unreadable.

See [`docs/config.mdx`](docs/config.mdx) for every environment variable and its validation rules.

## Development

Run the server:

```bash
cd server
go run .
```

Run the web client with hot reload. It proxies `/api` and `/ws` to `https://localhost:9788`:

```bash
cd client
bun install
bun run dev
```

Open `https://localhost:5173`.

Other useful commands:

| Command                                  | Purpose                                                      |
| ---------------------------------------- | ------------------------------------------------------------ |
| `cd scripts && ./generate-protobuf.sh`   | Regenerate Go and TypeScript protobuf code.                  |
| `cd scripts && ./generate-openapi-ts.sh` | Regenerate `client/src/lib/types/request-schema.d.ts`.       |
| `cd scripts && ./generate-websockets-ts.sh` | Regenerate WebSocket models from Go types.                |
| `cd server && go test ./...`             | Run server tests.                                            |
| `cd docs && mint dev`                    | Preview the documentation site.                              |

The server writes `openapi.json` to the repository root on startup. Regenerate the client schema after API changes.

## Documentation

The docs live in `docs/` and are built with [Mintlify](https://mintlify.com):

```bash
cd docs
mint dev
```

## Third-party assets

File and folder icons in `client/src/lib/assets/icons` come from [Material Icon Theme](https://github.com/material-extensions/vscode-material-icon-theme) by Material Extensions, used under the MIT license (see `client/src/lib/assets/icons/LICENSE`).

## Roadmap

- [ ] Handle root process signaling
- [ ] Handle denied permissions, for example creating, reading, and deleting files or directories
- [ ] Fix action on added new server
- [ ] Add creation of files and directories with root support
- [ ] Handle protected SSH keys
- [ ] Handle permission denied errors
- [ ] Handle deleting multiple items in the client
- [ ] Handle done on stream finished
- [ ] Implement cron job 200 MB log storage limit
- [ ] Add Docker management
- [ ] Add a desktop GUI client
