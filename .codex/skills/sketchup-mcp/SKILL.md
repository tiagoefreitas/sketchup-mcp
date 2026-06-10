---
name: sketchup-mcp
description: Use the live SketchUp MCP bridge to inspect, create, transform, validate, snapshot, select, and undo geometry in SketchUp. Trigger when the user asks Codex to work with the current SketchUp model, drive SketchUp, create or modify model geometry, inspect selections/components/groups, render snapshots, validate placement/contact/clearance, or run safe Ruby through the local SketchUp MCP server.
---

# SketchUp MCP

Use the project-local SketchUp MCP server to operate the live SketchUp model through tools, not by editing extension source. The bridge normally listens on `127.0.0.1:9876`, and Codex is configured to launch:

```sh
/Users/coolkcah/Documents/sketchup/bin/sketchup-mcp --host localhost --port 9876
```

If the MCP tools are unavailable in the current session, first verify the bridge:

```sh
python3 /Users/coolkcah/Documents/sketchup/scripts/smoke_test.py --timeout 10
```

If the smoke test fails, run `/Users/coolkcah/Documents/sketchup/scripts/install_sketchup_extension.sh`, restart SketchUp 2025, and retry. The direct install autostarts the Ruby bridge unless `SKETCHUP_MCP_AUTOSTART=0` is set.

## Operating Loop

1. Start with `ping` and `units_info`.
2. Inspect existing model state with `find_groups`, `list_instances`, `list_definitions`, `get_selection`, `measure`, or `inspect_geometry`.
3. Prefer structured tools over `eval_ruby`.
4. Batch related mutations with `batch_create` so the user's action is one undo step.
5. Verify after every meaningful change with `validate_geometry`, `closest_points`, `intersect_ray`, `measure`, and/or `export_scene(format="png")`.
6. Use `undo_last` when a mutation went wrong and the user has not asked to keep it.

Always report exact entity names/IDs used for destructive actions.

## Tool Selection

Use `create_component` for simple cubes/cylinders/spheres/cones.

Use `create_extrusion` for real modeling work: wall plates, studs, rafters, sheathing, panels, brackets, and any 2D profile pushed through a depth. Prefer this over custom Ruby for prism-like geometry. Inputs are in SketchUp internal units: inches.

Use `batch_create` when creating or mutating multiple pieces. It supports primitive create ops, extrusion create ops, translate/move_to/delete/replace ops, `pattern_linear`, and `mirror` in one transaction.

Use `find_groups` to locate named groups by prefix, regex, bounds intersection, parent, recursive traversal, and component inclusion.

Use `list_instances` and `list_definitions` for component inventory and model discovery. Use `select` to set the SketchUp UI selection by entity IDs.

Use `measure` for one entity's bounds, origin, material, name, and definition. Use `inspect_geometry` for per-face normals, areas, and loops.

Use `pattern_linear` for repeated groups along a vector, and `mirror_component` for symmetry across axis-aligned or arbitrary planes.

Use `replace_geometry` when the name/material/layer of a group should survive but its geometry must change.

Use `boolean_op` only when SketchUp Pro Solid Tools are available and both inputs are manifold solid groups.

Use `set_material` for named colors, hex colors, or existing SketchUp materials.

Use `export_scene` with `format="png"` as the visual feedback loop. Pass `width`, `height`, and `camera` when the current viewport is not the right framing.

Use `eval_ruby` only for capabilities not represented by structured tools. See `references/sketchup-rules.md` before mutating geometry through Ruby.

## Verification Tools

Use `validate_geometry` for repeatable positional assertions:

- `bounds`: entity bounds min/max equals expected coordinates.
- `contact`: one side of group A touches the opposing side of group B.
- `aligned`: multiple targets share min/max/center on an axis.
- `no_overlap`: targets do not penetrate; use `mode="obb"` for rotated local-frame solids, otherwise default AABB.

Use `closest_points` to check clearance/contact/overlap between two groups. This is better than bounding boxes for cavity fits, mortise/tenon-like seating, notches, and sloped/irregular pieces.

Use `intersect_ray` to ask SketchUp where model geometry actually is. This is the right tool for "what is the roof/rafter height at this X/Y?", "where does this sloped face sit?", or "does this vertical line hit the named group?"

## SketchUp Safety Rules

Follow these when using MCP tools or `eval_ruby`:

- One user-visible operation should be one undo step. Prefer `batch_create`; in Ruby, wrap mutations in `model.start_operation(..., true)`, `commit_operation`, and `abort_operation` on error.
- Validate geometry before creating it. Reject zero-length vectors, duplicate points, colinear face profiles, zero/negative dimensions where invalid, and non-numeric coordinates.
- SketchUp stores lengths internally in inches. Use `units_info` when interpreting user-facing model units; do not silently mix cm/meters with tool inputs.
- Do not use file export or temporary files for calculations. Use `measure`, `inspect_geometry`, `find_groups`, `intersect_ray`, `closest_points`, and `validate_geometry`.
- Do not naively traverse huge models in Ruby. Prefer `find_groups`/`list_instances`; if Ruby traversal is necessary, collect entities before mutating, use `Set`/`Hash` for lookups, and impose scope/depth/limit controls.
- Do not modify locked entities unless explicitly asked.
- When editing a group through Ruby, call `make_unique` first if shared instances could be affected.
- Check `entity.deleted?`/`entity.valid?` before accessing entities saved from earlier operations.
- Avoid global variables, monkey-patching, and long opaque `eval_ruby` snippets. Keep Ruby namespaced or local.

For detailed reminders, read `references/sketchup-rules.md`. For the complete MCP tool catalog and examples, read `references/tools.md`.
