# SketchUp MCP Project Setup

## Chosen Base

This project is forked from `lumberbarons/sketchup-mcp`.

Compared with the reviewed alternatives:

- `mhyrr/sketchup-mcp@aa096f0` is the Python FastMCP baseline with a useful compatibility fix and response parsing hardening.
- `rezaakm/sketchup-mcp@reliability-and-capability-overhaul` improves the Python transport and adds operational tools, but it has less test coverage and no release pipeline.
- `lumberbarons/sketchup-mcp` keeps the hardened response parsing idea, replaces the Python MCP layer with a static Go binary, has CI, release assets, an RBZ build path, and a broad Ruby-side test suite for geometry behavior.

## Local Runtime

The project-local MCP binary is installed at:

```sh
/Users/coolkcah/Documents/sketchup/bin/sketchup-mcp
```

The rebuilt SketchUp extension package is installed in the repo at:

```sh
/Users/coolkcah/Documents/sketchup/dist/su_mcp_vlocal.rbz
```

Install the RBZ in SketchUp with:

1. Open SketchUp.
2. Window > Extension Manager > Install Extension.
3. Select `dist/su_mcp_vlocal.rbz`.
4. Restart SketchUp.
5. Start the bridge with Extensions > MCP Server > Start Server.

## Codex MCP Config

The active Codex config contains:

```toml
[mcp_servers.sketchup]
command = "/Users/coolkcah/Documents/sketchup/bin/sketchup-mcp"
args = ["--host", "localhost", "--port", "9876"]
enabled = true
startup_timeout_sec = 10
tool_timeout_sec = 300
```

Restart Codex after changing MCP configuration.

## Verification

Ruby tests that do not require SketchUp:

```sh
make test-ruby
```

Direct TCP smoke test after SketchUp is running and the extension server is started:

```sh
python3 scripts/smoke_test.py --verbose
```

Runtime binary sanity check:

```sh
bin/sketchup-mcp --help
```

Go source tests require a Go toolchain:

```sh
go test ./...
```
