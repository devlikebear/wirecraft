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

Status: Next. GitHub issue: [#30](https://github.com/devlikebear/wirecraft/issues/30).

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

- [ ] Add actuator type and input signal primitives.
- [ ] Implement piston target extension/retraction from High/Low input.
- [ ] Clamp movement by tick delta and movement speed.
- [ ] Add focused tests for target calculation and movement step behavior.
- [ ] Update local handoff docs when complete.

## Acceptance Criteria

- [ ] `go test ./internal/actuator/...` passes.
- [ ] `go test ./...` passes.
- [ ] No frontend files are touched.

## Verification Commands

- `go test ./internal/actuator/...`
- `go test ./...`

---

## Planned Work Orders

- [ ] **WO-26: Add actuator block types** — add piston, motor, motor driver, and transistor switch block/protocol types.
- [ ] **WO-27: Integrate actuator update order** — run actuator input mapping and physics update after circuit evaluation.
- [ ] **WO-28: Add transform snapshot schema** — include actuator transforms in authoritative snapshots.
- [ ] **WO-29: Render actuator meshes** — render piston or motor entities from interpolated server transforms.
- [ ] **WO-30: Add actuator placement UI** — expose actuator block placement and minimal orientation controls.
- [ ] **WO-31: Add sensor input extension point** — separate button/proximity-style sensor primitives for later inputs.
- [ ] **WO-32: Add motor and driver component cards** — explain motor driver/transistor constraints and simulation simplifications.
- [ ] **Phase 3 Checkpoint** — verify button-driven actuator motion across server tests, client build, and browser smoke.

## Session Handoff

Start the next session from [current-status.md](./current-status.md), then continue with [#30 WO-25](https://github.com/devlikebear/wirecraft/issues/30).
