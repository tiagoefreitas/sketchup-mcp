# SketchUp Rules for MCP Work

These rules are distilled from the SketchUp Ruby API guide for the narrow case of using MCP tools and occasional `eval_ruby`. They are not extension-authoring instructions.

## Undo and Transactions

Prefer `batch_create` for multiple changes because it is one SketchUp transaction and rolls back on failure.

When using `eval_ruby` for mutation:

```ruby
model = Sketchup.active_model
model.start_operation('MCP action', true)
begin
  # mutate model here
  model.commit_operation
rescue StandardError
  model.abort_operation
  raise
end
```

Do not nest operations unless you know SketchUp's transparent operation semantics.

## Geometry Hygiene

- Reject zero-length vectors before normalizing or using them as directions.
- Reject duplicate or colinear points before `add_face`.
- Check face creation return values; `add_face` can fail.
- Use tolerances; do not compare floating point lengths with `==`.
- Use `model.active_entities` when creating geometry in the user's current context unless intentionally operating at model root.
- Create groups first, then add geometry inside `group.entities`.

## Units

SketchUp stores geometry in inches. MCP tool coordinates and dimensions are inches. Use `units_info` before converting from user-facing metric dimensions.

For Ruby length parsing from user text, prefer SketchUp's length parser (`"1m".to_l`, `"3'6\"".to_l`) instead of hand parsing.

## Entity Safety

- Check `entity.valid?` and `entity.deleted?` before using an entity reference saved from an earlier tool call.
- Filter out locked entities unless the user explicitly asked to alter locked objects.
- If Ruby edits a group that may share a definition, call `group.make_unique` first to avoid changing sibling instances unexpectedly.
- Do not store long-lived Ruby entity references across independent user requests; resolve by ID/name again.

## Inspection and Calculation

Do not export temporary files to measure geometry. Use:

- `measure` for bounds and size.
- `inspect_geometry` for faces, normals, loops, and areas.
- `intersect_ray` for model-derived positions on surfaces.
- `closest_points` for clearance/contact.
- `validate_geometry` for repeatable checks.

When exact face area is critical and glued components may affect `Face#area` in some SketchUp versions, use a mesh-based calculation in Ruby instead of trusting a single area call.

## Traversal and Performance

Prefer `find_groups` and `list_instances` over model-wide Ruby scans.

If Ruby traversal is necessary:

- Scope traversal to `active_entities`, a parent group, or a known component definition.
- Collect entities into an array before mutating.
- Use `Set` or `Hash` for membership/deduplication.
- Add recursion depth and result limits for large models.
- Cache `model`, `entities`, `selection`, and repeated lookups outside tight loops.

## Visual and Numeric Verification

Do not rely on mental geometry after mutations. Verify with at least one of:

- `validate_geometry` for known dimensions/alignment/contact/no-overlap.
- `closest_points` for clearance/contact/penetration.
- `intersect_ray` for sloped/derived locations.
- `export_scene(format="png")` for visual inspection.

Use `undo_last` promptly when verification shows the change is wrong.
