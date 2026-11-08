# go-my-mcp

MCP server provides SearXNG search and generic web page fetching.

## Tools

| Tool | Description | Input | Output |
| --- | --- | --- | --- |
| `web_search` | Web search using SearXNG backend | `query` (string) | `results[]` — `url`, `title`, `content`, `engine` |
| `web_fetch` | Fetch a web page by URL, following redirects | `url` (string, HTTP/HTTPS) | `content`, `url` (final), `status_code`, `content_type`, `redirect_count` |

Notes:

- `web_fetch` follows up to `RedirectMax` redirects (default 10) and times out after 15s. Every hop must be a plain `http`/`https` URL without embedded credentials; `redirect_count` reports how many hops were followed, and `url` is the final URL. `RedirectMax: 0` disables following — the 3xx response is returned as-is.
- Response body limits are configurable per tool: `FetchBodySizeMax` for `web_fetch`, `SearchBodySizeMax` for `web_search` (both default to 2 MiB).
- `web_search` requires a SearXNG instance with the JSON output format enabled.

## Requirements

- Go 1.27.2 or newer
- A reachable SearXNG instance (only for `web_search`)

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

## Configuration

### Precedence

Lowest to highest:

1. Built-in defaults
2. Config file
3. Environment variables
4. Command line flags

### Config File

`~/.config/go-my-mcp.json`

```json
{
  "ApiKey": "",
  "ListenAddress": ":9000",
  "SearXngUrl": "http://localhost:8080",
  "RedirectMax": 10,
  "FetchBodySizeMax": 2097152,
  "SearchBodySizeMax": 2097152,
  "Debug": false
}
```

| Key | Type | Default | Description |
| --- | --- | --- | --- |
| `ApiKey` | string | *(empty)* | Shared secret clients must present; empty disables access control |
| `ListenAddress` | string | `:9000` | Address the HTTP server listens on |
| `SearXngUrl` | string | `http://localhost:8080` | Base URL of the SearXNG instance |
| `RedirectMax` | int | `10` | Max redirects followed by `web_fetch`; `0` disables following |
| `FetchBodySizeMax` | int | `2097152` (2 MiB) | `web_fetch` response body limit in bytes |
| `SearchBodySizeMax` | int | `2097152` (2 MiB) | `web_search` response body limit in bytes |

`ApiKey` has no flag — it is read from the config file or the environment only, so it stays out of shell history and `ps` output. `config` always shows it as `<redacted>`.

`RedirectMax` must be `>= 0`; the body limits must be `> 0`. Out-of-range values are rejected at startup instead of silently misbehaving.

### Environment Variables

Prefix `GO_MY_MCP_`, upper-cased key name:

```sh
GO_MY_MCP_LISTENADDRESS=127.0.0.1:8080 GO_MY_MCP_SEARXNGURL=http://searxng.local go-my-mcp start
```

### Flags

These are persistent (root) flags, so they work with every sub command:

| Flag | Key | Default | Description |
| --- | --- | --- | --- |
| `-c`, `--FileConf` | `FileConf` | `~/.config/go-my-mcp.json` | Config file path |
| `-l`, `--ListenAddress` | `ListenAddress` | `:9000` | HTTP listen address |
| `-u`, `--SearXngUrl` | `SearXngUrl` | `http://localhost:8080` | SearXNG base URL |
| | `--RedirectMax` | `10` | Max redirects followed by `web_fetch` |
| | `--FetchBodySizeMax` | `2097152` | `web_fetch` body limit in bytes |
| | `--SearchBodySizeMax` | `2097152` | `web_search` body limit in bytes |
| `-d`, `--Debug` | `Debug` | `false` | Debug logging |

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

## VS Code MCP Server Setup

Pick HTTP endpoint. Or create `$HOME/.copilot/mcp-config.json`:

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

With API key

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

The server must be running before the client connects.

## Project layout

| Path | Purpose |
| --- | --- |
| [main.go](main.go) | Entry point, calls `cmd.Execute()` |
| [cmd/](cmd) | CLI commands: `root`, `start`, `config` |
| [global/](global) | Configuration type, viper wiring and shared helpers |
| [server/](server) | MCP server construction and HTTP transport |
| [tool/](tool) | Tool interface and implementations (`searxng`, `webfetch`) |

## Adding Tool

1. New tool file should go in `tool` directory
2. Implement the `tool.Tool` interface
3. Create an `init()` function to add an instance to the `tool.Tools` array

    ```go
    func init() { Tools.Addnew(New_Tool_Type) }
    ```

4. Add options in `conf.Conf.go` and `cmd.start.go`

## License

MIT — see [LICENSE](LICENSE).
