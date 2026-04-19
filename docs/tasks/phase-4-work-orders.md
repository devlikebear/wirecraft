# Phase 4 Work Orders — Multiplayer Physics Sync

_Last updated: 2026-04-19_

Parent issue: [#39 Phase 4: Multiplayer Physics Sync](https://github.com/devlikebear/wirecraft/issues/39)

## MVP Reset Notice

Status: Paused by [#49 MVP Reset: Workbench Device Builder](https://github.com/devlikebear/wirecraft/issues/49).

The completed Phase 4 work remains part of the foundation: room/session model, presence, deterministic command ordering, command acknowledgements, changed-set snapshots, client delta application, viewport navigation, and placement facing.

Do not continue the remaining generic collision/metrics/checkpoint work as the immediate next step. Resume it only when it is tied to the reset MVP device slice, especially the button-controlled sliding door.

## Phase Goal

Strengthen WireCraft for 2-4 collaborative users in the same room by separating room/session state, defining deterministic command conflict handling, improving snapshot efficiency, adding basic actuator collision constraints, and exposing tick/snapshot/client observability.

## Work Order 33: Add Server Room Model

Status: Completed. GitHub issue: [#40](https://github.com/devlikebear/wirecraft/issues/40).

## Goal

Introduce a small server-side room model that owns a simulation instance, tracks joined clients, and exposes apply/step/subscribe behavior for the default room without changing the existing `/ws` route behavior.

## Non-goals

- Do not add public lobby or dynamic room selection yet.
- Do not change WebSocket protocol shape.
- Do not implement command conflict ordering yet.
- Do not implement delta snapshots yet.

## Touch points (<=5)

- `internal/server/room.go`
- `internal/server/room_test.go`
- `internal/server/server.go`
- `internal/server/ws.go`
- `docs/tasks/phase-4-work-orders.md`

## Steps

- [x] Add `RoomID` and `Room` types with one simulation per room.
- [x] Move subscriber/client-count ownership behind the room model.
- [x] Keep the server on a single default room for now.
- [x] Add tests for join/leave, apply command, and snapshot stats per room.
- [x] Keep existing WebSocket tests compatible.

## Acceptance Criteria

- [x] `go test ./internal/server/...` passes.
- [x] `go test ./...` passes.
- [x] Existing `/ws` route still streams snapshots and applies commands.

## Verification Commands

- `go test ./internal/server/...`
- `go test ./...`

---

## Work Order 34: Add Client Presence Metadata

Status: Completed. GitHub issue: [#41](https://github.com/devlikebear/wirecraft/issues/41).

## Goal

Expose lightweight client presence metadata so snapshots/debug UI can show who is connected to the default room and how many clients are present.

## Non-goals

- Do not add login, accounts, or public profile data.
- Do not add lobby or room selection UI yet.
- Do not change command conflict handling yet.
- Do not add delta snapshots yet.

## Touch points (<=5)

- `internal/server/client.go`
- `internal/netproto/snapshot.go`
- `internal/sim/snapshot.go`
- `web/src/net/protocol.ts`
- `web/src/debug/DebugOverlay.ts`

## Steps

- [x] Add a minimal server-side client identity model.
- [x] Include client count and basic presence metadata in snapshots.
- [x] Parse the presence metadata in the TypeScript protocol.
- [x] Show the presence count in the debug overlay.
- [x] Keep existing snapshot consumers compatible.

## Acceptance Criteria

- [x] `go test ./...` passes.
- [x] `cd web && npm test` passes.
- [x] `cd web && npm run build` passes.

## Verification Commands

- `go test ./...`
- `cd web && npm test`
- `cd web && npm run build`

---

## Work Order 35: Add Deterministic Command Ordering

Status: Completed. GitHub issue: [#42](https://github.com/devlikebear/wirecraft/issues/42).

## Goal

Define and implement deterministic command ordering so same-tick edits are processed consistently across clients and duplicate commands are ignored.

## Non-goals

- Do not add command acknowledgement fields yet.
- Do not add snapshot delta support yet.
- Do not add UI conflict resolution controls yet.
- Do not change room selection or lobby behavior.

## Touch points (<=5)

- `internal/sim/commands.go`
- `internal/sim/commands_test.go`
- `internal/sim/simulation.go`
- `internal/server/room.go`
- `docs/tasks/phase-4-work-orders.md`

## Steps

- [x] Document and encode the command ordering key.
- [x] Add tests for same-coordinate same-tick command ordering.
- [x] Ignore duplicate command IDs from the same client.
- [x] Keep room command application deterministic.
- [x] Update local task docs after verification.

## Acceptance Criteria

- [x] `go test ./internal/sim/...` passes.
- [x] `go test ./...` passes.
- [x] Existing WebSocket command flow remains compatible.

## Verification Commands

- `go test ./internal/sim/...`
- `go test ./...`

---

## Work Order 36: Add Command Acknowledgement Snapshot Fields

Status: Completed. GitHub issue: [#43](https://github.com/devlikebear/wirecraft/issues/43).

## Goal

Expose lightweight command acknowledgement metadata in snapshots so clients can tell which submitted commands were accepted or rejected.

## Non-goals

- Do not build full conflict resolution UI yet.
- Do not add delta snapshots yet.
- Do not add lobby or account identity.
- Do not change room selection behavior.

## Touch points (<=5)

- `internal/netproto/snapshot.go`
- `internal/sim/commands.go`
- `internal/sim/simulation.go`
- `web/src/net/protocol.ts`
- `web/src/net/socket.ts`

## Steps

- [x] Define compact command acknowledgement snapshot fields.
- [x] Record accepted/rejected command IDs during command application.
- [x] Parse acknowledgement metadata in the TypeScript protocol.
- [x] Preserve existing WebSocket command flow.
- [x] Update local task docs after verification.

## Acceptance Criteria

- [x] `go test ./...` passes.
- [x] `cd web && npm test` passes.
- [x] `cd web && npm run build` passes.

## Verification Commands

- `go test ./...`
- `cd web && npm test`
- `cd web && npm run build`

---

## Work Order 36A: Improve Circuit Block Visual Readability

Status: Completed. GitHub issue: [#45](https://github.com/devlikebear/wirecraft/issues/45).

## Goal

Make starter circuit blocks visually distinguishable enough that users can actually compose and inspect simple circuits in the scene.

## Non-goals

- Do not change block protocol values or server simulation behavior.
- Do not add full UX polish, tutorials, or blueprint editing yet.
- Do not implement browser-side delta snapshots.
- Do not add complex imported 3D assets.

## Touch points (<=5)

- `web/src/render/VoxelRenderer.ts`
- `web/src/render/VoxelRenderer.test.ts`
- `docs/tasks/phase-4-work-orders.md`
- `docs/tasks/current-status.md`

## Steps

- [x] Add visual profiles for starter circuit block types.
- [x] Render Wire as a low thin conductor instead of a cube.
- [x] Render Power/Button/AND/MCU Output with distinct silhouettes.
- [x] Keep raycast/edit behavior compatible with existing block positions.
- [x] Verify with unit tests, build, and a browser smoke screenshot.

## Acceptance Criteria

- [x] `cd web && npm test` passes.
- [x] `cd web && npm run build` passes.
- [x] Browser smoke shows circuit blocks with visibly distinct shapes.
- [x] `go test ./...` remains green.

## Verification Commands

- `cd web && npm test`
- `cd web && npm run build`
- `go test ./...`

---

## Work Order 37: Add Changed-Set Snapshot Primitives

Status: Completed. GitHub issue: [#44](https://github.com/devlikebear/wirecraft/issues/44).

## Goal

Introduce snapshot primitives that can represent full snapshots and changed-set snapshots with changed blocks, removed blocks, and changed entities while preserving full snapshot fallback.

## Non-goals

- Do not implement browser-side delta application yet.
- Do not remove full snapshots.
- Do not add compression or binary protocol.
- Do not change room selection or lobby behavior.

## Touch points (<=5)

- `internal/netproto/snapshot.go`
- `internal/sim/snapshot_builder.go`
- `internal/sim/snapshot_builder_test.go`
- `internal/sim/snapshot.go`
- `docs/tasks/phase-4-work-orders.md`

## Steps

- [x] Define full vs changed-set snapshot schema fields.
- [x] Add changed block and removed block primitives.
- [x] Add changed entity primitives.
- [x] Preserve periodic/full snapshot compatibility.
- [x] Update local task docs after verification.

## Acceptance Criteria

- [x] `go test ./internal/sim/...` passes.
- [x] `go test ./...` passes.
- [x] Existing full snapshot consumers remain compatible.

## Verification Commands

- `go test ./internal/sim/...`
- `go test ./...`

---

## Work Order 38: Apply Client Delta Snapshots

Status: Completed. GitHub issue: [#46](https://github.com/devlikebear/wirecraft/issues/46).

## Goal

Teach the browser state layer to accept full snapshots and apply changed-set snapshot payloads produced by the server snapshot primitives.

## Non-goals

- Do not remove full snapshot support.
- Do not add compression or binary protocol.
- Do not implement server-side delta streaming policy yet beyond consuming the schema.
- Do not redesign the UI.

## Touch points (<=5)

- `web/src/net/protocol.ts`
- `web/src/net/protocol.test.ts`
- `web/src/state/snapshotStore.ts`
- `web/src/state/snapshotStore.test.ts`
- `docs/tasks/phase-4-work-orders.md`

## Steps

- [x] Parse full vs changed-set snapshot fields in TypeScript.
- [x] Add state-store logic for changed blocks and removed blocks.
- [x] Add entity changed-set application.
- [x] Preserve full snapshot reset behavior.
- [x] Update local task docs after verification.

## Acceptance Criteria

- [x] `cd web && npm test` passes.
- [x] `cd web && npm run build` passes.
- [x] Existing full snapshot rendering remains compatible.

## Verification Commands

- `cd web && npm test`
- `cd web && npm run build`

---

## Work Order 38A: Improve Viewport Usability and Directional Placement

Status: Completed. GitHub issue: [#48](https://github.com/devlikebear/wirecraft/issues/48).

## Goal

Remove confusing debug-only motion from the default view, add usable camera navigation, and add the first server-authoritative block facing path for direction-sensitive placement.

## Non-goals

- Do not implement full pin-aware circuit graph rules in this work order.
- Do not redesign the full toolbar or mode system.
- Do not replace the server-authoritative command/snapshot model.
- Do not implement actuator collision here; keep Work Order 39 as the next physics work order.

## Touch points (<=5 groups)

- `internal/world/*`
- `internal/netproto/*`
- `internal/sim/*`
- `web/src/input/*`
- `web/src/render/*`

## Steps

- [x] Hide the Phase 1 debug mover from the default client render path.
- [x] Add orbit, zoom, pan, and keyboard planar camera navigation.
- [x] Add placement facing metadata to commands, world blocks, snapshots, and the TypeScript protocol.
- [x] Add `R` key facing rotation for block placement.
- [x] Rotate rendered voxel instances from snapshot facing data.
- [x] Update local task docs after verification.

## Acceptance Criteria

- [x] `go test ./...` passes.
- [x] `cd web && npm test` passes.
- [x] `cd web && npm run build` passes.
- [x] Default runtime view no longer renders the moving debug mover.
- [x] Block facing survives command -> world -> snapshot -> client parse/render flow.

## Verification Commands

- `go test ./...`
- `cd web && npm test`
- `cd web && npm run build`

---

## Work Order 39: Add Basic Actuator Collision Constraints

Status: Paused / superseded as immediate next work by [#49](https://github.com/devlikebear/wirecraft/issues/49). GitHub issue: [#47](https://github.com/devlikebear/wirecraft/issues/47).

## Goal

Add a small physics collision policy that prevents actuator motion from entering occupied solid block cells, using a simple stop/clamp response.

## Non-goals

- Do not implement a full rigid body solver.
- Do not add joints, friction, or continuous collision detection.
- Do not redesign actuator rendering or placement UI.
- Do not add multiplayer-specific metrics in this work order.

## Touch points (<=5)

- `internal/physics/collision.go`
- `internal/physics/collision_test.go`
- `internal/actuator/piston.go`
- `internal/sim/simulation.go`
- `docs/tasks/phase-4-work-orders.md`

## Steps

- [ ] Add a collision/occupancy helper for actuator target cells.
- [ ] Clamp or stop actuator motion when the target cell is occupied by a solid block.
- [ ] Integrate the collision policy into piston update flow.
- [ ] Preserve existing actuator behavior when the path is clear.
- [ ] Update local task docs after verification.

## Acceptance Criteria

- [ ] `go test ./internal/physics/...` passes.
- [ ] `go test ./...` passes.
- [ ] Existing actuator snapshots remain compatible.

## Verification Commands

- `go test ./internal/physics/...`
- `go test ./...`

---

## Planned Work Orders

- [ ] **WO-40: Add server metrics logging** — paused until the reset MVP has a useful device loop to measure.
- [ ] **Phase 4 checkpoint** — paused until multiplayer is reintroduced after the single-user workbench/device loop.

## Session Handoff

Start the next session from [current-status.md](./current-status.md), then continue with [mvp-reset-work-orders.md](./mvp-reset-work-orders.md).
