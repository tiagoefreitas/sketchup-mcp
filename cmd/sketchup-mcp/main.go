// Command sketchup-mcp is the MCP server that talks to the SketchUp Ruby
// extension on localhost:9876. It serves MCP over stdio.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lumberbarons/sketchup-mcp/internal/skpclient"
	"github.com/lumberbarons/sketchup-mcp/internal/tools"
)

// version is stamped at link time by GoReleaser; "dev" in local builds.
var version = "dev"

func main() {
	host := flag.String("host", envString("SKETCHUP_MCP_HOST", "localhost"), "SketchUp Ruby extension host")
	port := flag.Int("port", envInt("SKETCHUP_MCP_PORT", 9876), "SketchUp Ruby extension port")
	dialTimeout := flag.Duration("dial-timeout", envDuration("SKETCHUP_MCP_DIAL_TIMEOUT", 3*time.Second), "SketchUp TCP dial timeout")
	callTimeout := flag.Duration("call-timeout", envDuration("SKETCHUP_MCP_CALL_TIMEOUT", 120*time.Second), "SketchUp request/response timeout")
	flag.Parse()

	// Logs must go to stderr — stdout is the MCP transport.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	slog.Info("SketchupMCP server starting", "version", version)

	client := skpclient.New(*host, *port)
	client.Timeout = *dialTimeout
	client.CallTimeout = *callTimeout
	reachable := client.Probe()
	if !reachable {
		slog.Warn("SketchUp not reachable; make sure the extension is running and Start Server has been clicked")
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "SketchupMCP",
		Version: version,
	}, &mcp.ServerOptions{
		Instructions: buildInstructions(*host, *port, reachable),
	})
	tools.RegisterAll(server, client)

	ctx := context.Background()
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		fmt.Fprintln(os.Stderr, "server exited with error:", err)
		os.Exit(1)
	}
}

func envString(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func envInt(name string, fallback int) int {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid %s=%q, using %d\n", name, value, fallback)
		return fallback
	}
	return n
}

func envDuration(name string, fallback time.Duration) time.Duration {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}
	seconds, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid %s=%q, using %s\n", name, value, fallback)
		return fallback
	}
	return time.Duration(seconds * float64(time.Second))
}

// buildInstructions returns the MCP Instructions string, appending a
// reachability advisory when the startup probe failed so the MCP client
// surfaces the diagnostic rather than only learning about it on the first
// failing tool call. Advisory-only — the extension can come up after the
// MCP server starts.
func buildInstructions(host string, port int, reachable bool) string {
	base := "Sketchup integration through the Model Context Protocol"
	if reachable {
		return base
	}
	return fmt.Sprintf(
		"%s\n\nNOTE: the SketchUp extension was not reachable on %s:%d at startup. Tool calls will fail until you open SketchUp, install the su_mcp extension, and click Extensions → MCP Server → Start Server.",
		base, host, port,
	)
}
