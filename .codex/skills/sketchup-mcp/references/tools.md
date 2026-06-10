# SketchUp MCP Tool Reference

## Connection

- `ping({})`: health check. Returns `{pong, version, time}`.
- `units_info({})`: model unit settings and `inches_per_centimeter` / `centimeters_per_inch`.

## Discovery

- `get_selection({})`: current SketchUp selection.
- `find_groups({name_prefix?, name_pattern?, in_bounds?, parent_id?, limit?, recursive?, include_components?})`: group search. Bounds are inches and intersection-based.
- `list_instances({definition_name?, name_pattern?, bounds?, limit?, recursive?, include_components?})`: groups and optionally component instances.
- `list_definitions({name_pattern?, include_bounds?})`: component definition inventory.
- `measure({id})`: one entity's type, name, definition, bounds, origin, material.
- `inspect_geometry({id? name?, include_vertices?})`: face normals, areas, loops, and optional vertices for a group.

## Creation and Mutation

- `create_component({type, name?, position?, dimensions?})`: cube/cylinder/sphere/cone. Prefer only for primitives.
- `create_extrusion({name, profile, extrude_axis+extrude_from+extrude_to OR plane+extrude_depth, holes?, material?})`: primary prism/solid creation tool.
- `batch_create({transaction_name?, operations})`: atomic multi-op create/mutate/delete batch.
- `transform_component({id? name?, move_to?, position?, rotation?, scale?})`: `move_to` is absolute bounds.min; `position` is relative delta.
- `delete_component({id? name?})`: strict one-target delete.
- `replace_geometry({id? name?, geometry, recursive?})`: swap group geometry while preserving identity metadata.
- `pattern_linear({id? name?, vector, count, include_source?, name_template?})`: additional copies along vector.
- `mirror_component({id? name?, axis+offset OR plane, include_source?, name_template?})`: mirror copy or flip in place.
- `set_material({id, material})`: named material, known color, or `#RRGGBB`.
- `select({ids})`: replace SketchUp UI selection.
- `undo_last({steps?})`: undo recent SketchUp operations.

## Validation and Visual Feedback

- `validate_geometry({assertions})`: read-only assertion batch.
- `closest_points({a, b, tolerance?})`: clearance/contact/overlap between two named/ID groups.
- `intersect_ray({origin, direction, target?, max_distance?, include_back_faces?})`: raycast against model or a target group.
- `export_scene({format, width?, height?, camera?})`: use `format="png"` for snapshots. `camera` accepts `eye`, `target`, optional `up`, `perspective`, `fov`.

## Ruby Escape Hatch

- `eval_ruby({code})`: use only when structured tools cannot express the task. Mutating Ruby must be undo-wrapped and defensive.

## Batch Operation Shapes

Primitive create:

```json
{"op":"cube","name":"Block 1","position":[0,0,0],"dimensions":[16,16,8],"material":"#8B4513"}
```

Extrusion create:

```json
{"op":"extrusion","name":"Rafter W 1","profile":[[0,0],[72,36],[72,41.5],[0,5.5]],"extrude_axis":"y","extrude_from":0,"extrude_to":1.5}
```

Mutation:

```json
{"op":"move_to","name":"Ridge","target":[0,0,96]}
{"op":"translate","id":12345,"position":[0,0,1.5]}
{"op":"delete","name":"Temporary Brace"}
{"op":"replace","name":"Panel A","geometry":{"type":"extrusion","profile":[[0,0],[48,0],[48,96],[0,96]],"extrude_axis":"y","extrude_from":0,"extrude_to":0.5}}
```

Replication:

```json
{"op":"pattern_linear","name":"Stud 1","vector":[16,0,0],"count":5,"name_template":"Stud {n}"}
{"op":"mirror","name":"Rafter W 1","axis":"x","offset":120,"name_template":"Rafter E {n}"}
```
