#!/usr/bin/env bash
# Install the SketchUp Ruby bridge directly into the user's SketchUp Plugins
# folder. This avoids needing to click through Extension Manager during local
# development and ensures the TCP bridge autoloads on SketchUp startup.
set -euo pipefail

version="${SKETCHUP_VERSION:-2025}"
plugins_dir="${SKETCHUP_PLUGINS_DIR:-$HOME/Library/Application Support/SketchUp ${version}/SketchUp/Plugins}"
repo_root="$(cd "$(dirname "$0")/.." && pwd)"

if [ ! -d "$plugins_dir" ]; then
  echo "SketchUp Plugins folder not found: $plugins_dir" >&2
  echo "Set SKETCHUP_VERSION or SKETCHUP_PLUGINS_DIR and retry." >&2
  exit 1
fi

rm -rf "$plugins_dir/su_mcp" "$plugins_dir/su_mcp.rb"
mkdir -p "$plugins_dir/su_mcp"

cp "$repo_root/su_mcp/su_mcp.rb" "$plugins_dir/su_mcp.rb"
cp "$repo_root/su_mcp/extension.json" "$plugins_dir/su_mcp/extension.json"
cp -R "$repo_root/su_mcp/su_mcp/." "$plugins_dir/su_mcp/"

cat > "$plugins_dir/zz_su_mcp_autostart.rb" <<'RUBY'
# Autoload SketchUp MCP without requiring Extension Manager enablement.
begin
  require File.join(__dir__, 'su_mcp', 'main')
rescue Exception => e
  begin
    SKETCHUP_CONSOLE.write("MCP autostart failed: #{e.class}: #{e.message}\n#{Array(e.backtrace).first(10).join("\n")}\n")
  rescue Exception
    puts "MCP autostart failed: #{e.class}: #{e.message}"
  end
end
RUBY

echo "Installed SketchUp MCP extension to:"
echo "  $plugins_dir/su_mcp.rb"
echo "  $plugins_dir/su_mcp/"
echo "  $plugins_dir/zz_su_mcp_autostart.rb"
