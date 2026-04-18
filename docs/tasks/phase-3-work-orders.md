# WireCraft — Phase 3 Work Orders

_Last updated: 2026-04-18_

Parent phase: [#28 Phase 3: Physical Actuators](https://github.com/devlikebear/wirecraft/issues/28)

## Phase Goal

Connect server-authoritative circuit signal state to deterministic physical actuator motion. Phase 3 starts with pure physics and actuator domain primitives, then maps circuit output to kinematic piston or motor transforms and renders those transforms through the existing snapshot interpolation path.

## Phase Non-goals

- Do not add complex rigid body collision or joint solving.
- Do not implement analog current, torque, or PWM electrical accuracy.
- Do not add full robot templates or factory-line tooling.
- Do not add persistence beyond the current in-memory runtime.

---

## Work Order 24: Add Dynamic Entity and Transform Primitives

Status: Completed. GitHub issue: [#29](https://github.com/devlikebear/wirecraft/issues/29).

## Goal

Add the Phase 3 physics foundation: value types for transforms and dynamic entities, plus deterministic entity ordering.

## Non-goals

- Do not add actuator block types yet.
- Do not connect circuits to motion yet.
- Do not change snapshot schema or client rendering yet.
- Do not remove the existing debug mover snapshot behavior.

## Touch points (<=5)

- `internal/physics/transform.go`
- `internal/physics/entity.go`
- `internal/physics/entity_test.go`
- `docs/tasks/phase-3-work-orders.md`
- `docs/tasks/current-status.md`

## Steps

- [x] Add `Vec3`, `Quat`, and `Transform` value types.
- [x] Add `EntityID`, `EntityType`, and `DynamicEntity` primitives.
- [x] Add deterministic entity sorting/copy helpers.
- [x] Add focused tests for ordering and defensive copies.
- [x] Update local handoff docs when complete.

## Acceptance Criteria

- [x] `go test ./internal/physics/...` passes.
- [x] `go test ./...` passes.
- [x] No frontend files are touched.

## Verification Commands

- `go test ./internal/physics/...`
- `go test ./...`

---

## Work Order 25: Add Actuator Component Model

Status: Completed. GitHub issue: [#30](https://github.com/devlikebear/wirecraft/issues/30).

## Goal

Implement the first actuator domain behavior: piston and motor actuator primitives that turn a digital signal into deterministic movement targets.

## Non-goals

- Do not add new world block types yet.
- Do not wire actuators into simulation snapshots yet.
- Do not add client rendering or toolbar UI yet.
- Do not implement full rigid body physics.

## Touch points (<=5)

- `internal/actuator/types.go`
- `internal/actuator/piston.go`
- `internal/actuator/piston_test.go`
- `docs/tasks/phase-3-work-orders.md`
- `docs/tasks/current-status.md`

## Steps

- [x] Add actuator type and input signal primitives.
- [x] Implement piston target extension/retraction from High/Low input.
- [x] Clamp movement by tick delta and movement speed.
- [x] Add focused tests for target calculation and movement step behavior.
- [x] Update local handoff docs when complete.

## Acceptance Criteria

- [x] `go test ./internal/actuator/...` passes.
- [x] `go test ./...` passes.
- [x] No frontend files are touched.

## Verification Commands

- `go test ./internal/actuator/...`
- `go test ./...`

---

## Work Order 26: Add Actuator Block Types

Status: Completed. GitHub issue: [#31](https://github.com/devlikebear/wirecraft/issues/31).

## Goal

Add the world/protocol/circuit-facing block type foundation for Phase 3 actuator placement. Introduce piston, motor, motor driver, and transistor switch types without wiring them into simulation motion yet.

## Non-goals

- Do not implement actuator movement in simulation yet.
- Do not add transform snapshots yet.
- Do not add client actuator mesh rendering yet.
- Do not implement full wiring validation UI yet.

## Touch points (<=5)

- `internal/world/block.go`
- `internal/world/world_test.go`
- `internal/circuit/types.go`
- `internal/circuit/types_test.go`
- `web/src/net/protocol.ts`

## Steps

- [x] Add piston, motor, motor driver, and transistor switch block types while preserving existing numeric values.
- [x] Add actuator-facing metadata or role definitions needed by later wiring work.
- [x] Extend Go tests for valid block types and metadata completeness.
- [x] Update TypeScript protocol block type definitions.
- [x] Keep simulation behavior unchanged.

## Acceptance Criteria

- [x] `go test ./internal/world/...` passes.
- [x] `go test ./internal/circuit/...` passes.
- [x] `go test ./...` passes.
- [x] `cd web && npm test` passes.
- [x] `cd web && npm run build` passes.

## Verification Commands

- `go test ./internal/world/...`
- `go test ./internal/circuit/...`
- `go test ./...`
- `cd web && npm test`
- `cd web && npm run build`

---

## Work Order 27: Integrate Actuator Update Order

Status: Next. GitHub issue: [#32](https://github.com/devlikebear/wirecraft/issues/32).

## Goal

Integrate the actuator update order into the simulation loop without exposing final client rendering yet. The simulation should evaluate circuits, map actuator inputs, update actuator/physics state deterministically, then build snapshots.

## Non-goals

- Do not add client actuator rendering yet.
- Do not add toolbar placement UI beyond existing protocol support.
- Do not implement rigid body collision or joints.
- Do not add persistence.

## Touch points (<=5)

- `internal/sim/simulation.go`
- `internal/sim/simulation_test.go`
- `internal/actuator/types.go`
- `internal/actuator/piston.go`
- `docs/tasks/phase-3-work-orders.md`

## Steps

- [ ] Define how actuator blocks are discovered from world state.
- [ ] Map adjacent circuit High/Low state to actuator input.
- [ ] Update actuator/physics state after circuit evaluation and before snapshot build.
- [ ] Add tests for button or power signal updating actuator target deterministically.
- [ ] Keep existing snapshot behavior compatible until transform schema is expanded.

## Acceptance Criteria

- [ ] `go test ./internal/sim/...` passes.
- [ ] `go test ./internal/actuator/...` passes.
- [ ] `go test ./...` passes.

## Verification Commands

- `go test ./internal/sim/...`
- `go test ./internal/actuator/...`
- `go test ./...`

---

## Planned Work Orders

- [ ] **WO-28: Add transform snapshot schema** — include actuator transforms in authoritative snapshots.
- [ ] **WO-29: Render actuator meshes** — render piston or motor entities from interpolated server transforms.
- [ ] **WO-30: Add actuator placement UI** — expose actuator block placement and minimal orientation controls.
- [ ] **WO-31: Add sensor input extension point** — separate button/proximity-style sensor primitives for later inputs.
- [ ] **WO-32: Add motor and driver component cards** — explain motor driver/transistor constraints and simulation simplifications.
- [ ] **Phase 3 Checkpoint** — verify button-driven actuator motion across server tests, client build, and browser smoke.

## Session Handoff

Start the next session from [current-status.md](./current-status.md), then continue with [#32 WO-27](https://github.com/devlikebear/wirecraft/issues/32).
