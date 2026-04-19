# MVP Reset Work Orders — Workbench Device Builder

_Created: 2026-04-19_

Parent issue: [#49 MVP Reset: Workbench Device Builder](https://github.com/devlikebear/wirecraft/issues/49)

## Phase Goal

Turn WireCraft from a broad technical prototype into a small usable workbench where a user can build, inspect, edit, and run one interactive device: a button-controlled sliding door.

## Work Order R0: Reset MVP Planning

Status: Completed. GitHub issue: [#49](https://github.com/devlikebear/wirecraft/issues/49).

## Goal

Rewrite the active planning documents so the next development work is centered on Workbench Device Builder rather than continuing generic Phase 4 hardening.

## Steps

- [x] Update the PRD around the reset MVP.
- [x] Replace the roadmap with reset phases R0-R5.
- [x] Add a reset rationale document.
- [x] Add reset work orders.
- [x] Update current status and Phase 4 handoff docs.

## Acceptance Criteria

- [x] New sessions can start from `docs/tasks/current-status.md` and see #49 as the active parent.
- [x] #47 is no longer listed as the next implementation task.
- [x] The next implementation task is small, testable, and tied to block instance state.

---

## Work Order R1: Add Block Instance Properties To Snapshots

Status: Next. GitHub issue: [#50](https://github.com/devlikebear/wirecraft/issues/50).

## Goal

Add a minimal server-authoritative property container for placed blocks so later work can expose selected block state, actuator blocked/open state, and module attachment metadata without replacing the protocol again.

## Non-goals

- Do not add a property editing UI yet.
- Do not implement port-aware circuit extraction yet.
- Do not add module slots yet.
- Do not change block placement UX beyond preserving existing behavior.

## Touch points

- `internal/world/world.go`
- `internal/world/world_test.go`
- `internal/netproto/snapshot.go`
- `internal/sim/snapshot_builder.go`
- `web/src/net/protocol.ts`

## Steps

- [ ] Add a `Properties` field to placed world blocks with default empty state.
- [ ] Include block properties in snapshot block payloads.
- [ ] Parse optional block properties in the TypeScript protocol.
- [ ] Preserve compatibility with snapshots that omit properties.
- [ ] Add focused Go and TypeScript tests.

## Acceptance Criteria

- [ ] `go test ./internal/world/...` passes.
- [ ] `go test ./internal/netproto/...` passes.
- [ ] `go test ./internal/sim/...` passes.
- [ ] `cd web && npm test` passes.
- [ ] Existing place/remove flows keep working.

## Verification Commands

- `go test ./internal/world/...`
- `go test ./internal/netproto/...`
- `go test ./internal/sim/...`
- `cd web && npm test`

---

## Work Order R2: Select Placed Block And Show Instance State

Status: Planned.

## Goal

Let users select an already placed block and inspect that specific block instance rather than only seeing the toolbar-selected component card.

## Non-goals

- Do not add property editing yet.
- Do not add multi-select.
- Do not add blueprint selection.

## Touch points

- `web/src/input/EditController.ts`
- `web/src/render/VoxelRenderer.ts`
- `web/src/ui/InspectPanel.ts`
- `web/src/state/snapshotStore.ts`
- `web/src/main.ts`

## Steps

- [ ] Add inspect/select mode or click behavior for placed blocks.
- [ ] Track the selected block position in client state.
- [ ] Resolve selected block snapshot from the latest snapshot store.
- [ ] Show position, block type, facing, properties, and component card details.
- [ ] Add focused UI/unit tests where practical.

## Acceptance Criteria

- [ ] Selecting a placed wire/button/piston updates the inspect panel.
- [ ] The panel distinguishes selected instance state from toolbar selection.
- [ ] `cd web && npm test` passes.
- [ ] `cd web && npm run build` passes.

---

## Work Order R3: Rotate Existing Directional Blocks

Status: Planned.

## Goal

Support rotating an already placed direction-sensitive block so users can correct wire, board, and actuator orientation after placement.

## Non-goals

- Do not implement arbitrary transforms.
- Do not add multi-block rotation.
- Do not implement port-aware behavior yet.

## Touch points

- `internal/netproto/command.go`
- `internal/sim/simulation.go`
- `internal/world/world.go`
- `web/src/input/EditController.ts`
- `web/src/main.ts`

## Steps

- [ ] Add a small rotate/update-facing command.
- [ ] Validate that the target block exists and facing is valid.
- [ ] Preserve existing block type and properties.
- [ ] Add keyboard/UI path for selected block rotation.
- [ ] Add Go and TypeScript tests.

## Acceptance Criteria

- [ ] Selected block rotation updates server state.
- [ ] Facing survives command -> world -> snapshot -> render.
- [ ] `go test ./...` passes.
- [ ] `cd web && npm test` passes.

---

## Work Order R4: Define Port And Slot Metadata

Status: Planned.

## Goal

Introduce static metadata for block ports and attachment slots before changing circuit behavior.

## Non-goals

- Do not replace circuit extraction yet.
- Do not add attachment commands yet.
- Do not build arbitrary nested assemblies.

## Touch points

- `internal/component/card.go`
- `internal/circuit/types.go`
- `web/src/state/componentCards.ts`
- `web/src/ui/InspectPanel.ts`
- `docs/reference/component-cards.md`

## Steps

- [ ] Define port metadata with name, direction, signal kind, and face.
- [ ] Define slot metadata with name and accepted module ids.
- [ ] Add metadata for power, wire, button, control/driver, piston.
- [ ] Show ports and slots in inspect panel.
- [ ] Add metadata validation tests.

## Acceptance Criteria

- [ ] Component metadata distinguishes pins/ports/slots clearly.
- [ ] Directional components expose face-specific ports.
- [ ] `go test ./internal/component/...` passes.
- [ ] `cd web && npm test` passes.

---

## Work Order R5: Port-Aware Starter Circuit Extraction

Status: Planned.

## Goal

Move starter circuit graph extraction from body adjacency to compatible port adjacency for the reset MVP components.

## Non-goals

- Do not model analog voltage.
- Do not add schematic view.
- Do not support arbitrary custom boards.

## Touch points

- `internal/circuit/extract.go`
- `internal/circuit/extract_test.go`
- `internal/circuit/evaluate.go`
- `internal/sim/snapshot_builder.go`
- `docs/reference/component-cards.md`

## Steps

- [ ] Build edges only between compatible adjacent ports.
- [ ] Use block facing when resolving port faces.
- [ ] Preserve existing simple circuits where the new ports are compatible.
- [ ] Add tests for valid and invalid adjacent placements.
- [ ] Include invalid/warning state in a follow-up work order if needed.

## Acceptance Criteria

- [ ] Power -> wire -> output still evaluates High when ports align.
- [ ] Misaligned wire or gate port does not connect.
- [ ] `go test ./internal/circuit/...` passes.
- [ ] `go test ./...` passes.

---

## Work Order R6: Sliding Door Device Slice

Status: Planned.

## Goal

Build the first coherent device slice: a button-controlled sliding door that uses structure, signal, driver/control, actuator movement, and blocked feedback.

## Non-goals

- Do not add elevators, cars, or multi-stage controllers.
- Do not add full rigid body physics.
- Do not add device save/load yet.

## Touch points

- `internal/world/block.go`
- `internal/actuator/piston.go`
- `internal/sim/simulation.go`
- `web/src/render/VoxelRenderer.ts`
- `web/src/ui/InspectPanel.ts`

## Steps

- [ ] Add minimal door/frame/panel block types if needed.
- [ ] Drive actuator axis from block facing.
- [ ] Attach actuator behavior to the door panel through a constrained rule.
- [ ] Add blocked/open/closed properties to snapshots.
- [ ] Add a browser smoke scenario for the sliding door.

## Acceptance Criteria

- [ ] Button released keeps door closed.
- [ ] Button pressed opens door.
- [ ] Solid obstruction produces blocked state.
- [ ] `go test ./...` passes.
- [ ] `cd web && npm test` passes.
- [ ] `cd web && npm run build` passes.

---

## Session Handoff

Start the next implementation from **Work Order R1: Add Block Instance Properties To Snapshots** ([#50](https://github.com/devlikebear/wirecraft/issues/50)).
