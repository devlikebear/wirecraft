# WireCraft — Phase 1 Work Orders

_Last updated: 2026-04-18_

Phase source: [`../plans/wire-craft-phase-1-authoritative-voxel-loop.md`](../plans/wire-craft-phase-1-authoritative-voxel-loop.md)

## Phase 1 Goal

Build the smallest server-authoritative voxel loop: Go owns world state and ticks, TypeScript/Three.js renders the world, and clients send edit commands without directly mutating authoritative state.

## Work Order 1: Scaffold Go Project And Embedded Web Boundary

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

- [ ] Create `go.mod` with module `github.com/devlikebear/wirecraft`.
- [ ] Add `cmd/wirecraft-server/main.go` with flag-based host/port config and graceful startup.
- [ ] Add `internal/server/server.go` with `New()` returning an `http.Handler` or server struct.
- [ ] Add `internal/webui/embed.go` documenting the future `go:embed dist` boundary.
- [ ] Update `README.md` with local development and future embedded release shape.

## Acceptance Criteria

- [ ] `go test ./...` passes.
- [ ] `go run ./cmd/wirecraft-server` starts an HTTP server.
- [ ] `GET /healthz` returns 200 with a minimal response.
- [ ] README explains development mode versus embedded release mode.

## Verification Commands

- `go test ./...`
- `go run ./cmd/wirecraft-server`

---

## Work Order 2: Add Tick Clock And Voxel World Core

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

- [ ] Write tests for a 20Hz clock target duration and monotonically increasing tick IDs.
- [ ] Implement `TickID`, `Clock`, and tick duration helpers.
- [ ] Write tests for bounded `Position`, `BlockType`, `World.Get`, `World.Set`, and `World.Remove`.
- [ ] Implement a 32x32x16 default world with deterministic set/remove behavior.
- [ ] Keep block types minimal: `air`, `solid`, `debug_mover`.

## Acceptance Criteria

- [ ] Tick tests prove target duration and tick ID behavior.
- [ ] World tests reject out-of-bounds positions.
- [ ] World tests prove repeated set/remove operations are deterministic.
- [ ] No frontend files are touched.

## Verification Commands

- `go test ./internal/sim/...`
- `go test ./internal/world/...`
- `go test ./...`

---

## Work Order 3: Initialize Vite Three.js Client Skeleton

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

- [ ] Initialize `web/` as a Vite + TypeScript project using npm.
- [ ] Add Three.js dependency.
- [ ] Render a minimal scene with camera, renderer, grid helper, and basic lighting.
- [ ] Add scripts for `dev`, `build`, and `test` or a placeholder test command if the test runner is introduced later.
- [ ] Update `README.md` with Go server and Vite client commands.

## Acceptance Criteria

- [ ] `cd web && npm install` completes.
- [ ] `cd web && npm run build` succeeds.
- [ ] Opening the Vite dev server shows a nonblank Three.js scene.
- [ ] README gives exact commands for local development.

## Verification Commands

- `cd web && npm install`
- `cd web && npm run build`

