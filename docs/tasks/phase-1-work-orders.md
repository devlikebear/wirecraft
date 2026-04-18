# WireCraft — Phase 1 Work Orders

_Last updated: 2026-04-18_

Phase source: [`../plans/wire-craft-phase-1-authoritative-voxel-loop.md`](../plans/wire-craft-phase-1-authoritative-voxel-loop.md)

GitHub phase issue: [#1 Phase 1: Authoritative Voxel Loop](https://github.com/devlikebear/wirecraft/issues/1)

## Phase 1 Goal

Build the smallest server-authoritative voxel loop: Go owns world state and ticks, TypeScript/Three.js renders the world, and clients send edit commands without directly mutating authoritative state.

## Work Order 1: Scaffold Go Project And Embedded Web Boundary

Status: Completed. GitHub issue: [#2](https://github.com/devlikebear/wirecraft/issues/2).

## Goal

Create the repository's executable skeleton: Go module, server entrypoint, package layout, and a clear boundary for embedding the built TypeScript UI into the Go binary later.

## Non-goals

- Do not implement WebSocket protocol yet.
- Do not initialize Vite yet.
- Do not implement voxel simulation yet.

## Touch points (<=5)

- `go.mod`
- `cmd/wirecraft-server/main.go`
- `internal/server/server.go`
- `internal/webui/embed.go`
- `README.md`

## Steps

- [x] Create `go.mod` with module `github.com/devlikebear/wirecraft`.
- [x] Add `cmd/wirecraft-server/main.go` with flag-based host/port config and graceful startup.
- [x] Add `internal/server/server.go` with `New()` returning an `http.Handler` or server struct.
- [x] Add `internal/webui/embed.go` documenting the future `go:embed dist` boundary.
- [x] Update `README.md` with local development and future embedded release shape.

## Acceptance Criteria

- [x] `go test ./...` passes.
- [x] `go run ./cmd/wirecraft-server` starts an HTTP server.
- [x] `GET /healthz` returns 200 with a minimal response.
- [x] README explains development mode versus embedded release mode.

## Verification Commands

- `go test ./...`
- `go run ./cmd/wirecraft-server`

---

## Work Order 2: Add Tick Clock And Voxel World Core

Status: Completed. GitHub issue: [#3](https://github.com/devlikebear/wirecraft/issues/3).

## Goal

Implement test-first server simulation primitives: fixed tick timing and deterministic bounded voxel world operations.

## Non-goals

- Do not add network transport.
- Do not add persistence.
- Do not add circuit logic.

## Touch points (<=5)

- `internal/sim/tick.go`
- `internal/sim/tick_test.go`
- `internal/world/world.go`
- `internal/world/world_test.go`
- `internal/world/block.go`

## Steps

- [x] Write tests for a 20Hz clock target duration and monotonically increasing tick IDs.
- [x] Implement `TickID`, `Clock`, and tick duration helpers.
- [x] Write tests for bounded `Position`, `BlockType`, `World.Get`, `World.Set`, and `World.Remove`.
- [x] Implement a 32x32x16 default world with deterministic set/remove behavior.
- [x] Keep block types minimal: `air`, `solid`, `debug_mover`.

## Acceptance Criteria

- [x] Tick tests prove target duration and tick ID behavior.
- [x] World tests reject out-of-bounds positions.
- [x] World tests prove repeated set/remove operations are deterministic.
- [x] No frontend files are touched.

## Verification Commands

- `go test ./internal/sim/...`
- `go test ./internal/world/...`
- `go test ./...`

---

## Work Order 3: Initialize Vite Three.js Client Skeleton

Status: Completed. GitHub issue: [#4](https://github.com/devlikebear/wirecraft/issues/4).

## Goal

Create the TypeScript frontend scaffold with a minimal Three.js scene and package scripts that will later connect to the Go server.

## Non-goals

- Do not implement WebSocket client yet.
- Do not implement raycasting yet.
- Do not add advanced UI panels.

## Touch points (<=5)

- `web/package.json`
- `web/vite.config.ts`
- `web/src/main.ts`
- `web/src/styles.css`
- `README.md`

## Steps

- [x] Initialize `web/` as a Vite + TypeScript project using npm.
- [x] Add Three.js dependency.
- [x] Render a minimal scene with camera, renderer, grid helper, and basic lighting.
- [x] Add scripts for `dev`, `build`, and `test` or a placeholder test command if the test runner is introduced later.
- [x] Update `README.md` with Go server and Vite client commands.

## Acceptance Criteria

- [x] `cd web && npm install` completes.
- [x] `cd web && npm run build` succeeds.
- [x] Opening the Vite dev server shows a nonblank Three.js scene.
- [x] README gives exact commands for local development.

## Verification Commands

- `cd web && npm install`
- `cd web && npm run build`

---

## Work Order 4: Add Command And Snapshot Protocol Types

Status: Completed. GitHub issue: [#5](https://github.com/devlikebear/wirecraft/issues/5).

## Goal

Define the server-side protocol data types for client edit commands and world snapshots before adding WebSocket transport.

## Non-goals

- Do not implement WebSocket transport yet.
- Do not connect the TypeScript client yet.
- Do not implement raycasting or interpolation.

## Touch points (<=5)

- `internal/netproto/command.go`
- `internal/netproto/command_test.go`
- `internal/netproto/snapshot.go`
- `internal/netproto/snapshot_test.go`
- `docs/tasks/phase-1-work-orders.md`

## Steps

- [x] Define command types for `place_block` and `remove_block`.
- [x] Define command fields: `clientId`, `commandId`, `tickHint`, `position`, `blockType`.
- [x] Add command validation tests for valid edits, out-of-bounds positions, invalid block types, and unknown command types.
- [x] Define snapshot type with `tick`, `serverTimeMs`, `blocks`, `entities`, and `stats` fields.
- [x] Add JSON round-trip tests for command and snapshot payloads.

## Acceptance Criteria

- [x] `go test ./internal/netproto/...` passes.
- [x] `go test ./...` passes.
- [x] Protocol structs reuse `internal/world.Position` and `internal/world.BlockType` where appropriate.
- [x] No server transport or frontend files are touched.

## Verification Commands

- `go test ./internal/netproto/...`
- `go test ./...`

---

## Work Order 5: Build Simulation Snapshot From World State

Status: Completed. GitHub issue: [#6](https://github.com/devlikebear/wirecraft/issues/6).

## Goal

Create the first simulation assembly layer that can build a full snapshot from the authoritative world state using the protocol types.

## Non-goals

- Do not implement WebSocket transport yet.
- Do not connect the TypeScript client yet.
- Do not implement delta snapshots yet.

## Touch points (<=5)

- `internal/world/world.go`
- `internal/world/world_test.go`
- `internal/sim/snapshot.go`
- `internal/sim/snapshot_test.go`
- `docs/tasks/phase-1-work-orders.md`

## Steps

- [x] Add a deterministic way to list occupied blocks from `world.World`.
- [x] Add tests proving occupied blocks are listed in stable coordinate order.
- [x] Add a snapshot builder that emits a full `netproto.Snapshot` from tick, server time, world, and stats input.
- [x] Add tests proving snapshot blocks match world state and air blocks are omitted.
- [x] Keep dynamic entities empty for this work order.

## Acceptance Criteria

- [x] `go test ./internal/world/...` passes.
- [x] `go test ./internal/sim/...` passes.
- [x] `go test ./...` passes.
- [x] No WebSocket transport or frontend files are touched.

## Verification Commands

- `go test ./internal/world/...`
- `go test ./internal/sim/...`
- `go test ./...`

---

## Work Order 6: Add In-Memory Simulation Runner

Status: Completed. GitHub issue: [#7](https://github.com/devlikebear/wirecraft/issues/7).

## Goal

Add a small in-memory simulation runner that owns a world, applies validated commands, advances ticks, and produces snapshots without network transport.

## Non-goals

- Do not implement WebSocket transport yet.
- Do not connect the TypeScript client yet.
- Do not implement interpolation.

## Touch points (<=5)

- `internal/sim/simulation.go`
- `internal/sim/simulation_test.go`
- `internal/netproto/command.go`
- `internal/world/world.go`
- `docs/tasks/phase-1-work-orders.md`

## Steps

- [x] Add `Simulation` with default world bounds and current tick.
- [x] Add `ApplyCommand` for `place_block` and `remove_block` using command validation.
- [x] Add `Step` or `Snapshot` method that advances tick and returns a full snapshot.
- [x] Add tests for valid place/remove commands changing world state.
- [x] Add tests for invalid commands being rejected without changing world state.

## Acceptance Criteria

- [x] `go test ./internal/sim/...` passes.
- [x] `go test ./...` passes.
- [x] No WebSocket transport or frontend files are touched.

## Verification Commands

- `go test ./internal/sim/...`
- `go test ./...`

---

## Work Order 7: Add WebSocket Simulation Stream

Status: Completed. GitHub issue: [#8](https://github.com/devlikebear/wirecraft/issues/8).

## Goal

Expose the in-memory simulation through a server-owned WebSocket stream so clients can send edit commands and receive authoritative snapshots.

## Non-goals

- Do not connect the TypeScript client yet.
- Do not implement interpolation yet.
- Do not implement persistence or delta snapshots.

## Touch points (<=5)

- `go.mod`
- `internal/server/server.go`
- `internal/server/ws.go`
- `internal/server/ws_test.go`
- `docs/tasks/phase-1-work-orders.md`

## Steps

- [x] Choose and add a small maintained WebSocket package.
- [x] Register `GET /ws` on the server handler.
- [x] Create a server-owned simulation loop that sends full snapshots at a fixed tick rate.
- [x] Accept JSON edit commands from the WebSocket and apply them through `sim.Simulation`.
- [x] Add server tests for WebSocket connection, command receive, and snapshot response.

## Acceptance Criteria

- [x] `go test ./internal/server/...` passes.
- [x] `go test ./...` passes.
- [x] `/ws` rejects non-WebSocket requests cleanly.
- [x] No frontend files are touched.

## Verification Commands

- `go test ./internal/server/...`
- `go test ./...`

---

## Work Order 8: Add TypeScript WebSocket Client And Snapshot Store

Status: Completed. GitHub issue: [#9](https://github.com/devlikebear/wirecraft/issues/9).

## Goal

Connect the Vite client to the Go WebSocket endpoint and keep an in-memory snapshot buffer that later renderers can consume.

## Non-goals

- Do not implement voxel InstancedMesh rendering yet.
- Do not implement raycast editing yet.
- Do not implement interpolation yet.

## Touch points (<=5)

- `web/src/net/protocol.ts`
- `web/src/net/socket.ts`
- `web/src/state/snapshotStore.ts`
- `web/src/main.ts`
- `docs/tasks/phase-1-work-orders.md`

## Steps

- [x] Define TypeScript protocol types matching Go snapshot and command JSON.
- [x] Add a WebSocket client that connects to `/ws` and parses snapshots.
- [x] Add a snapshot store with append/latest/buffer length helpers.
- [x] Wire the client into `main.ts` for connection lifecycle and basic status logging or overlay-safe state.
- [x] Add build-time validation through the existing Vite build.

## Acceptance Criteria

- [x] `cd web && npm run build` passes.
- [x] Client code derives the WebSocket URL from the current page origin.
- [x] Snapshot parsing does not mutate authoritative world state locally.
- [x] No Go server files are touched.

## Verification Commands

- `cd web && npm run build`

---

## Work Order 9: Render Authoritative Snapshots As Voxels

GitHub issue: [#10 WO-9: Render authoritative snapshots as voxels](https://github.com/devlikebear/wirecraft/issues/10)

## Goal

Render server-authoritative snapshot blocks in the Three.js scene using an efficient voxel renderer path.

## Non-goals

- Do not implement raycast editing yet.
- Do not implement interpolation yet.
- Do not add delta snapshots.

## Touch points (<=5)

- `web/src/render/VoxelRenderer.ts`
- `web/src/render/VoxelRenderer.test.ts`
- `web/src/main.ts`
- `web/src/styles.css`
- `docs/tasks/phase-1-work-orders.md`

## Steps

- [ ] Add a `VoxelRenderer` that maps full snapshot blocks into Three.js mesh instances or stable meshes.
- [ ] Define block materials for solid and debug mover block types.
- [ ] Replace the static demo footprint with server snapshot-driven voxel updates.
- [ ] Keep the existing scene, camera, lights, and resize behavior intact.
- [ ] Add focused tests for snapshot-to-render data mapping where practical.

## Acceptance Criteria

- [ ] `cd web && npm test` passes.
- [ ] `cd web && npm run build` passes.
- [ ] With Go server and Vite dev server running, the browser scene reflects blocks received from `/ws` snapshots.
- [ ] No Go server files are touched.

## Verification Commands

- `cd web && npm test`
- `cd web && npm run build`
