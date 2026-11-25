# Servergram web client

The web client is the SvelteKit UI served by the Servergram server. It is a static single-page app: the server embeds the production build and serves it from `/`.

A desktop GUI client is planned but not available yet. Use this web UI for now.

## Stack

- SvelteKit 3 with Svelte 5 runes
- Tailwind CSS 4 and shadcn-svelte / bits-ui components
- `openapi-fetch` with types generated from the server's OpenAPI document
- LayerChart and D3 for resource charts
- A WebSocket for live service updates and server-sent events for process streams

## Features

- Server list, SSH keys, and groups
- Server overview with CPU, memory, disk, network, and top processes
- File browser with upload, rename, view, and delete
- Process list with stop and signal actions
- Cron job creation, editing, manual runs, and run history
- Activity logs and notifications
- Password login when the server sets `PASSWORD`

## Requirements

- Bun (recommended) or Node.js 20 or newer
- A running Servergram server on `https://localhost:9788`
- TLS certificates at `../certs/localhost+2.pem` and `../certs/localhost+2-key.pem`

The Vite config reads those certificate paths for both development and production builds. Create them first:

```bash
cd ..
mkdir -p certs
mkcert -install
mkcert \
  -cert-file certs/localhost+2.pem \
  -key-file certs/localhost+2-key.pem \
  localhost 127.0.0.1 ::1
```

## Development

```bash
bun install
bun run dev
```

Open `https://localhost:5173`. Vite proxies `/api` and `/ws` to `https://localhost:9788` and accepts the server's self-signed certificate.

## Build

```bash
bun run build
```

The static adapter writes the output to `../server/dist`. The server embeds that directory at compile time, so rebuild the server after building the client.

## Generated types

The client uses generated API and WebSocket types:

| File                                  | Generated from                              | Command (run from `scripts/`)       |
| ------------------------------------- | ------------------------------------------- | ----------------------------------- |
| `src/lib/types/request-schema.d.ts`   | `openapi.json` at the repository root       | `./generate-openapi-ts.sh`          |
| `src/lib/types/websockets/`           | `server/websockets/models` Go types         | `./generate-websockets-ts.sh`       |

The server writes `openapi.json` when it starts. Regenerate the client types after changing the API.

## Layout

| Path                | Contents                                                     |
| ------------------- | ------------------------------------------------------------ |
| `src/routes/`       | Pages for servers, files, processes, cron, logs, and login.  |
| `src/lib/api/`      | API client, auth handling, and WebSocket helpers.            |
| `src/lib/components/` | Shared UI components.                                      |
| `src/lib/types/`    | Generated and hand-written TypeScript types.                 |
| `src/service-worker/` | Service worker registration and caching.                   |

## Checks

```bash
bun run check
```
