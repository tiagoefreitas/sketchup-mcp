#!/usr/bin/env python3
"""Smoke-test the SketchUp Ruby extension over its TCP JSON-RPC socket.

Prerequisite: SketchUp is open, the su_mcp extension is installed, and
Extensions > MCP Server > Start Server has been clicked.
"""

from __future__ import annotations

import argparse
import json
import os
import socket
import sys
from typing import Any


DEFAULT_HOST = os.environ.get("SKETCHUP_MCP_HOST", "127.0.0.1")
DEFAULT_PORT = int(os.environ.get("SKETCHUP_MCP_PORT", "9876"))


def call_tool(host: str, port: int, name: str, arguments: dict[str, Any], timeout: float) -> dict[str, Any]:
    request = {
        "jsonrpc": "2.0",
        "method": "tools/call",
        "params": {"name": name, "arguments": arguments},
        "id": 1,
    }
    payload = (json.dumps(request) + "\n").encode("utf-8")

    with socket.create_connection((host, port), timeout=timeout) as sock:
        sock.settimeout(timeout)
        sock.sendall(payload)
        chunks: list[bytes] = []
        while True:
            chunk = sock.recv(65536)
            if not chunk:
                break
            chunks.append(chunk)
            if b"\n" in chunk:
                break

    raw = b"".join(chunks).strip()
    if not raw:
        raise RuntimeError("SketchUp extension closed the connection without a response")
    return json.loads(raw.decode("utf-8"))


def assert_success(response: dict[str, Any]) -> None:
    if response.get("error"):
        raise RuntimeError(f"JSON-RPC error: {response['error']}")

    result = response.get("result")
    if not isinstance(result, dict):
        raise RuntimeError(f"Expected object result, got: {result!r}")
    if result.get("success") is not True or result.get("isError"):
        raise RuntimeError(f"Tool returned failure: {result!r}")

    content = result.get("content")
    if not isinstance(content, list) or not content:
        raise RuntimeError(f"Missing response content: {result!r}")
    text = content[0].get("text") if isinstance(content[0], dict) else None
    if text != "3":
        raise RuntimeError(f"Expected eval_ruby result '3', got: {text!r}")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default=DEFAULT_HOST)
    parser.add_argument("--port", type=int, default=DEFAULT_PORT)
    parser.add_argument("--timeout", type=float, default=10.0)
    parser.add_argument("--verbose", action="store_true")
    args = parser.parse_args()

    try:
        response = call_tool(args.host, args.port, "eval_ruby", {"code": "1 + 2"}, args.timeout)
        if args.verbose:
            print(json.dumps(response, indent=2, sort_keys=True))
        assert_success(response)
    except Exception as exc:
        print(f"smoke test failed: {exc}", file=sys.stderr)
        return 1

    print(f"smoke test passed: SketchUp extension responded on {args.host}:{args.port}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
