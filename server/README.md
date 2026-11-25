# Servergram server

The server is the control plane for Servergram. It serves the web UI and REST API over HTTPS, keeps the application database, and manages SSH connections and agents on remote servers.

A single server binary embeds:

- the production web client from `server/dist`
- the Linux agents from `server/agents/agent_linux_amd64.bin` and `server/agents/agent_linux_arm64.bin`

Both are build outputs and are gitignored. Build the client and agents before building the server.

## Requirements

- Go 1.27.3 or newer
- A TLS certificate and key in PEM format
- A database: SQLite, PostgreSQL, MySQL, or SQL Server

## Build

From the repository root:

```bash
cd client && bun install && bun run build && cd ..
cd agent
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../server/agents/agent_linux_amd64.bin .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ../server/agents/agent_linux_arm64.bin .
cd ..
cd server
go build -o servergram-server .
```

The full instructions, including certificate setup, are in the root [`README.md`](../README.md).

## Run

Create `server/.env` or export the variables. A `.env` file is loaded from the current working directory.

```bash
SERVER_PORT="9788"
ALLOWED_ORIGINS="https://localhost:9788"
PASSWORD="ChangeMe-Str0ng!2026"
SHH_MASTER_KEY="ChangeMe-Too-Str0ng!2026"
SSL_PUBLIC_KEY="../certs/localhost+2.pem"
SSL_PRIVATE_KEY="../certs/localhost+2-key.pem"
```

```bash
./servergram-server
```

Open `https://localhost:9788`.

## Configuration

| Variable          | Required | Default                                     |
| ----------------- | -------- | ------------------------------------------- |
| `DATABASE_URL`    | No       | SQLite in the user config directory         |
| `SERVER_PORT`     | No       | `9788`                                      |
| `ALLOWED_ORIGINS` | Yes\*    | `*` when `SERVER_PORT` is unset             |
| `PASSWORD`        | No       | Empty, which disables authentication        |
| `SHH_MASTER_KEY`  | Yes      | None                                        |
| `SSL_PUBLIC_KEY`  | Yes      | None                                        |
| `SSL_PRIVATE_KEY` | Yes      | None                                        |

\* Required whenever `SERVER_PORT` is set.

When set, `PASSWORD` must be 16 to 128 printable characters with an uppercase letter, a lowercase letter, a number, and a symbol. `SHH_MASTER_KEY` is required and must satisfy the same rules. They must not be equal. Changing `SHH_MASTER_KEY` after SSH keys are saved makes those keys unreadable.

See [`../docs/config.mdx`](../docs/config.mdx) for the full reference.

## Endpoints

- `https://<host>:<port>/` serves the embedded web client.
- `/api/*` serves the REST API.
- `/ws/` serves the WebSocket used for live updates.
- `/docs` serves the generated API reference.
- `/openapi.json` serves the OpenAPI document.

The server also writes a copy of `openapi.json` to `../openapi.json` relative to the working directory at startup. From the `scripts/` directory, run `./generate-openapi-ts.sh` to regenerate the client schema after API changes.

## Layout

| Path                       | Contents                                                        |
| -------------------------- | --------------------------------------------------------------- |
| `main.go`, `lib.go`        | Startup, TLS, routing, database setup, frontend embedding.      |
| `agents/`                  | Embedded agent binaries and OS/architecture selection.          |
| `db/`                      | GORM models and migrations.                                     |
| `internal/env.go`          | Environment loading and validation.                             |
| `internal/sshmanager/`     | SSH connections, agent upload, and gRPC client setup.           |
| `internal/server.go`       | Route groups, middleware, SSE, and WebSocket plumbing.          |
| `routes/`                  | REST handlers.                                                  |
| `websockets/`              | WebSocket handlers and shared models.                           |
| `dist/`                    | Embedded web client build output.                               |

## Development

```bash
cd server
go run .
```

`air` works with the included `.air.toml`; it rebuilds the server on changes. Run `go test ./...` for the test suite.
