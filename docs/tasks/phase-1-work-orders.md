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

GitHub issue: [#5 WO-4: Add command and snapshot protocol types](https://github.com/devlikebear/wirecraft/issues/5)

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

- [ ] Define command types for `place_block` and `remove_block`.
- [ ] Define command fields: `clientId`, `commandId`, `tickHint`, `position`, `blockType`.
- [ ] Add command validation tests for valid edits, out-of-bounds positions, invalid block types, and unknown command types.
- [ ] Define snapshot type with `tick`, `serverTimeMs`, `blocks`, `entities`, and `stats` fields.
- [ ] Add JSON round-trip tests for command and snapshot payloads.

## Acceptance Criteria

- [ ] `go test ./internal/netproto/...` passes.
- [ ] `go test ./...` passes.
- [ ] Protocol structs reuse `internal/world.Position` and `internal/world.BlockType` where appropriate.
- [ ] No server transport or frontend files are touched.

## Verification Commands

- `go test ./internal/netproto/...`
- `go test ./...`
