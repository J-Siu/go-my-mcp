# go-my-mcp

MCP server provides
— SearXNG search
- Generic web page fetching — over streamable HTTP.

## Tools

| Tool | Description | Input | Output |
| --- | --- | --- | --- |
| `searxng_search` | Web search via SearXNG | `query` (string) | `results[]` — `url`, `title`, `content`, `engine` |
| `web_fetch` | Fetch a web page by URL | `url` (string, HTTP/HTTPS) | `content`, `url`, `status_code`, `content_type` |

Notes:

- `web_fetch` follows no redirects (returned as-is) and times out after 15s.
- Both tools cap the response body at 2 MiB.
- `searxng_search` requires a SearXNG instance with the JSON output format enabled.

## Requirements

- Go 1.27.2 or newer
- A reachable SearXNG instance (only for `searxng_search`)

## Build

```sh
go build -o go-my-mcp .
```

Optionally install it into your `PATH`:

```sh
go install .
```

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

### Show the version

```sh
go-my-mcp version
```

```text
go-my-mcp version v0.0.0-20261102083237-cb4d9e817523
```

The version comes from the build info the Go toolchain stamps into the binary, so nothing has to be set by hand:

| Build | Reported version |
| --- | --- |
| Tagged, clean checkout | the tag, e.g. `v1.2.3` |
| Untagged commit | pseudo-version, e.g. `v0.0.0-20261102083237-cb4d9e817523` |
| Uncommitted changes | suffix `+dirty` |
| `go run .` | `(devel)` |

The same string is reported to MCP clients as `serverInfo.version` and shown by `config`.

### Start the server

```sh
go-my-mcp start
```

```text
MCP server listening on :9000
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
  "ListenAddress": ":9000",
  "SearXngUrl": "http://xng.vms.local",
  "Debug": false,
  "Version": "v0.0.0-20261102083237-cb4d9e817523",
  "Env": {
    "Prefix": "GO_MY_MCP",
    "Recognized": [
      "GO_MY_MCP_FILECONF",
      "GO_MY_MCP_LISTENADDRESS",
      "GO_MY_MCP_SEARXNGURL",
      "GO_MY_MCP_DEBUG"
    ],
    "Set": {}
  }
}
```

Use it to confirm exactly which values are in effect after combining the config file, environment variables and flags.

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
  "ListenAddress": "127.0.0.1:8080",
  "SearXngUrl": "http://xng.vms.local",
  "Debug": false,
  "Version": "v0.0.0-20261102083237-cb4d9e817523",
  "Env": {
    "Prefix": "GO_MY_MCP",
    "Recognized": [
      "GO_MY_MCP_FILECONF",
      "GO_MY_MCP_LISTENADDRESS",
      "GO_MY_MCP_SEARXNGURL",
      "GO_MY_MCP_DEBUG"
    ],
    "Set": {
      "GO_MY_MCP_LISTENADDRESS": "127.0.0.1:8080"
    }
  }
}
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
  "SearXngUrl": "http://xng.vms.local"
}
```

| Key | Type | Default | Description |
| --- | --- | --- | --- |
| `ListenAddress` | string | `:9000` | Address the HTTP server listens on |
| `SearXngUrl` | string | `http://xng.vms.local` | Base URL of the SearXNG instance |

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
| `-u`, `--SearXngUrl` | `SearXngUrl` | `http://xng.vms.local` | SearXNG base URL |
| `-d`, `--debug` | `Debug` | `false` | Debug logging |

## MCP client setup

Point your MCP client at the running server's HTTP endpoint. For example, in VS Code `mcp.json`:

```json
{
  "servers": {
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
| [mcpSrv/](mcpSrv) | MCP server construction and HTTP transport |
| [tool/](tool) | Tool interface and implementations (`searxng`, `webfetch`) |
| [lib/](lib) | Small shared helpers (URL validation) |

## Development

```sh
go build ./...   # build
go vet ./...     # static checks
gofmt -l .       # list files needing formatting
```

## License

MIT — see [LICENSE](LICENSE).
