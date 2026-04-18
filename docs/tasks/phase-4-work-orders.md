# Phase 4 Work Orders — Multiplayer Physics Sync

_Last updated: 2026-04-18_

Parent issue: [#39 Phase 4: Multiplayer Physics Sync](https://github.com/devlikebear/wirecraft/issues/39)

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

Status: Next. GitHub issue: [#43](https://github.com/devlikebear/wirecraft/issues/43).

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

- [ ] Define compact command acknowledgement snapshot fields.
- [ ] Record accepted/rejected command IDs during command application.
- [ ] Parse acknowledgement metadata in the TypeScript protocol.
- [ ] Preserve existing WebSocket command flow.
- [ ] Update local task docs after verification.

## Acceptance Criteria

- [ ] `go test ./...` passes.
- [ ] `cd web && npm test` passes.
- [ ] `cd web && npm run build` passes.

## Verification Commands

- `go test ./...`
- `cd web && npm test`
- `cd web && npm run build`

---

## Planned Work Orders

- [ ] **WO-37: Add changed-set snapshot primitives** — represent full vs delta snapshot payloads.
- [ ] **WO-38: Apply client delta snapshots** — update the browser state store for delta application.
- [ ] **WO-39: Add basic actuator collision constraints** — clamp actuator motion against occupied solids.
- [ ] **WO-40: Add server metrics logging** — expose tick duration, queue length, bytes, and client count.
- [ ] **Phase 4 checkpoint** — verify 2-4 client collaboration, conflict handling, collision, and observability.

## Session Handoff

Start the next session from [current-status.md](./current-status.md), then continue with [#43 WO-36](https://github.com/devlikebear/wirecraft/issues/43).
