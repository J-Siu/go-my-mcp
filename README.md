# go-my-mcp

MCP server provides SearXNG search and generic web page fetching.

## Tools

| Tool | Description | Input | Output |
| --- | --- | --- | --- |
| `searxng_search` | Web search via SearXNG | `query` (string) | `results[]` — `url`, `title`, `content`, `engine` |
| `web_fetch` | Fetch a web page by URL, following redirects | `url` (string, HTTP/HTTPS) | `content`, `url` (final), `status_code`, `content_type`, `redirect_count` |

Notes:

- `web_fetch` follows up to `RedirectMax` redirects (default 10) and times out after 15s. Every hop must be a plain `http`/`https` URL without embedded credentials; `redirect_count` reports how many hops were followed, and `url` is the final URL. `RedirectMax: 0` disables following — the 3xx response is returned as-is.
- Response body limits are configurable per tool: `FetchBodySizeMax` for `web_fetch`, `SearchBodySizeMax` for `searxng_search` (both default to 2 MiB).
- `searxng_search` requires a SearXNG instance with the JSON output format enabled.

## Requirements

- Go 1.27.2 or newer
- A reachable SearXNG instance (only for `searxng_search`)

## Build and Install

```sh
go build
```

```sh
go install
```

Docker: [docker_my_mcp](https://github.com/J-Siu/docker_my_mcp)

## Usage

```sh
go-my-mcp [command]
```

| Command | Description |
| --- | --- |
| `start`, `s` | Start the MCP server over HTTP |
| `config`, `c` | Show the effective configuration as JSON |
| `version`, `v` | Show the version |
| `help` | Help about any command |
| `completion` | Generate a shell completion script |

### Start the server

```sh
go-my-mcp start
```

The MCP endpoint is served at `/mcp` using the streamable HTTP transport:

```text
http://localhost:9000/mcp
```

### Show the effective configuration

```sh
go-my-mcp config
```

```json
{
  "FileConf": "~/.config/go-my-mcp.json",
  "ApiKey": "",
  "ListenAddress": ":9000",
  "SearXngUrl": "http://localhost:8080",
  "RedirectMax": 10,
  "FetchBodySizeMax": 2097152,
  "SearchBodySizeMax": 2097152,
  "Debug": false,
  "Version": "v1.0.0",
  "Source": {
    "ApiKey": "default",
    "Debug": "default",
    "FetchBodySizeMax": "default",
    "FileConf": "default",
    "ListenAddress": "default",
    "RedirectMax": "default",
    "SearXngUrl": "default",
    "SearchBodySizeMax": "default"
  },
  "Env": {
    "Prefix": "GO_MY_MCP",
    "Recognized": [
      "GO_MY_MCP_APIKEY",
      "GO_MY_MCP_FILECONF",
      "GO_MY_MCP_LISTENADDRESS",
      "GO_MY_MCP_SEARXNGURL",
      "GO_MY_MCP_REDIRECTMAX",
      "GO_MY_MCP_FETCHBODYSIZEMAX",
      "GO_MY_MCP_SEARCHBODYSIZEMAX",
      "GO_MY_MCP_DEBUG"
    ],
    "Set": {}
  }
}
```

Use it to confirm exactly which values are in effect after combining the config file, environment variables and flags.

The `Source` block reports where each value came from — `default`, the `config` file, `env`, or `flag` — so a value that looks wrong can be traced to its origin immediately.

The `Env` block reports the environment side of that resolution:

| Field | Description |
| --- | --- |
| `Prefix` | Prefix used to build the variable names |
| `Recognized` | Every variable name viper reads, whether set or not |
| `Set` | Variables actually present in the environment, with their values |

With overrides active:

```sh
GO_MY_MCP_LISTENADDRESS=127.0.0.1:8080 go-my-mcp config
```

```json
{
  "FileConf": "~/.config/go-my-mcp.json",
  "ApiKey": "",
  "ListenAddress": "127.0.0.1:8080",
  "SearXngUrl": "http://localhost:8080",
  "RedirectMax": 10,
  "FetchBodySizeMax": 2097152,
  "SearchBodySizeMax": 2097152,
  "Debug": false,
  "Version": "v1.0.0",
  "Source": {
    "ApiKey": "default",
    "Debug": "default",
    "FetchBodySizeMax": "default",
    "FileConf": "default",
    "ListenAddress": "env",
    "RedirectMax": "default",
    "SearXngUrl": "default",
    "SearchBodySizeMax": "default"
  },
  "Env": {
    "Prefix": "GO_MY_MCP",
    "Recognized": [
      "GO_MY_MCP_APIKEY",
      "GO_MY_MCP_FILECONF",
      "GO_MY_MCP_LISTENADDRESS",
      "GO_MY_MCP_SEARXNGURL",
      "GO_MY_MCP_REDIRECTMAX",
      "GO_MY_MCP_FETCHBODYSIZEMAX",
      "GO_MY_MCP_SEARCHBODYSIZEMAX",
      "GO_MY_MCP_DEBUG"
    ],
    "Set": {
      "GO_MY_MCP_LISTENADDRESS": "127.0.0.1:8080"
    }
  }
}
```

`Source` reports the highest-precedence origin, so in the output above `ListenAddress` is `env` even if the config file also set it. Adding a flag moves the entry to `flag`:

```sh
GO_MY_MCP_LISTENADDRESS=127.0.0.1:8080 go-my-mcp --listen :9999 config
```

```json
  "Source": {
    "ListenAddress": "flag",
    ...
  },
```

`config` writes JSON to stdout, so it can be piped directly:

```sh
go-my-mcp config | jq '.Source'
```

Any `GO_MY_MCP_`-prefixed variable appears under `Set`, so a typo such as `GO_MY_MCP_LISENADDRESS` shows up there while missing from `Recognized` — a quick way to spot variables that are being ignored.

## Configuration

### Precedence

Lowest to highest:

1. Built-in defaults
2. Config file
3. Environment variables
4. Command line flags

### Config file

The default path is `~/.config/go-my-mcp.json`. A missing file is **not** an error — the defaults are used instead. A file that exists but cannot be parsed is an error.

Override the path with the root `--config` flag:

```sh
go-my-mcp start --config ./local.json
```

`~` and environment variables in the path are expanded.

Example:

```json
{
  "ListenAddress": "127.0.0.1:9000",
  "SearXngUrl": "http://localhost:8080",
  "RedirectMax": 10,
  "FetchBodySizeMax": 2097152,
  "SearchBodySizeMax": 2097152
}
```

| Key | Type | Default | Description |
| --- | --- | --- | --- |
| `ApiKey` | string | *(empty)* | Shared secret clients must present; empty disables access control |
| `ListenAddress` | string | `:9000` | Address the HTTP server listens on |
| `SearXngUrl` | string | `http://localhost:8080` | Base URL of the SearXNG instance |
| `RedirectMax` | int | `10` | Max redirects followed by `web_fetch`; `0` disables following |
| `FetchBodySizeMax` | int | `2097152` (2 MiB) | `web_fetch` response body limit in bytes |
| `SearchBodySizeMax` | int | `2097152` (2 MiB) | `searxng_search` response body limit in bytes |

`ApiKey` has no flag — it is read from the config file or the environment only, so it stays out of shell history and `ps` output. `config` always shows it as `<redacted>`.

`RedirectMax` must be `>= 0`; the body limits must be `> 0`. Out-of-range values are rejected at startup instead of silently misbehaving.

### Environment variables

Prefix `GO_MY_MCP_`, upper-cased key name:

```sh
GO_MY_MCP_LISTENADDRESS=127.0.0.1:8080 GO_MY_MCP_SEARXNGURL=http://searxng.local go-my-mcp start
```

### Flags

These are persistent (root) flags, so they work with every sub command:

| Flag | Key | Default | Description |
| --- | --- | --- | --- |
| `-c`, `--config` | `FileConf` | `~/.config/go-my-mcp.json` | Config file path |
| `-l`, `--listen` | `ListenAddress` | `:9000` | HTTP listen address |
| `-u`, `--SearXngUrl` | `SearXngUrl` | `http://localhost:8080` | SearXNG base URL |
| | `--RedirectMax` | `10` | Max redirects followed by `web_fetch` |
| | `--FetchBodySizeMax` | `2097152` | `web_fetch` body limit in bytes |
| | `--SearchBodySizeMax` | `2097152` | `searxng_search` body limit in bytes |
| `-d`, `--debug` | `Debug` | `false` | Debug logging |

### Access control

Set `ApiKey` (via `GO_MY_MCP_APIKEY` or the config file) to require a key on **every** request — the MCP handshake included, since the server is usually reachable over the network:

```sh
GO_MY_MCP_APIKEY=s3cr3t go-my-mcp start
```

Clients may present it either way:

```sh
curl -H "Authorization: Bearer s3cr3t"   ...   # standard bearer token
curl -H "X-API-Key: s3cr3t"              ...   # simpler, common for machine clients
```

Requests without a valid key get `401 Unauthorized` plus a `WWW-Authenticate: Bearer` challenge, and are logged at `WARN` with the remote address. The comparison is constant time. If `Authorization` is present but malformed (e.g. `Basic ...`), it is rejected rather than falling back to `X-API-Key`.

When `ApiKey` is empty the server requires nothing, so an existing setup keeps working. That also means an unset key is *not* a safe default for an exposed host — set one there.

For VS Code `mcp.json`, add the header to the server entry:

```json
{
  "servers": {
    "go-my-mcp": {
      "type": "http",
      "url": "http://localhost:9000/mcp",
      "headers": { "Authorization": "Bearer s3cr3t" }
    }
  }
}
```

### Logging

Logs go to **stderr** as [slog](https://pkg.go.dev/log/slog) text, so they never mix with `config`/`version` output on stdout:

```text
time=... level=DEBUG msg=TypeWebFetch.Fetch url=http://example.com/start
time=... level=DEBUG msg=TypeWebFetch.Fetch redirect=1 to=http://example.com/mid
time=... level=DEBUG msg=TypeWebFetch.Fetch status=200 content_len=40 elapsed_ms=3 content_type="text/html; charset=utf-8"
time=... level=ERROR msg=TypeWebFetch.Fetch status=404 content_len=0 content_type="" err="HTTP 404 Not Found" elapsed_ms=1
```

`Debug: false` (the default) keeps the level at `INFO`, where the tools are silent — only errors reaching the MCP client are surfaced, and the client already receives them.

```sh
go-my-mcp start --debug        # or Debug: true / GO_MY_MCP_DEBUG=true
```

Debug adds one line per request (URL, each redirect hop, status, body length, and `elapsed_ms`), plus the resolved tool settings at startup. Response bodies are never logged — they are returned to the client.

Each request is logged exactly once at the tool boundary, on success or failure, so errors are not duplicated. The MCP SDK is wired to the same `slog` logger, so protocol events such as session connect/initialize are visible at `INFO`:

```text
time=... level=INFO msg="server connecting"
time=... level=INFO msg="server session connected" session_id=GBPUED452JZAOAPEXWC6FHFYVN
time=... level=INFO msg="session initialized"
```

## MCP client setup

Point your MCP client at the running server's HTTP endpoint. For example, in VS Code `$HOME/.copilot/mcp-config.json`:

```json
{
  "mcpServers": {
    "go-my-mcp": {
      "type": "http",
      "url": "http://localhost:9000/mcp"
    }
  }
}
```

The server must be running before the client connects.

## Project layout

| Path | Purpose |
| --- | --- |
| [main.go](main.go) | Entry point, calls `cmd.Execute()` |
| [cmd/](cmd) | CLI commands: `root`, `start`, `config` |
| [global/](global) | Configuration type, viper wiring and shared helpers |
| [server/](server) | MCP server construction and HTTP transport |
| [tool/](tool) | Tool interface and implementations (`searxng`, `webfetch`) |

## License

MIT — see [LICENSE](LICENSE).
