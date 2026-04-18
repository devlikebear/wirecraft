# WireCraft — Phase 2 Work Orders

_Last updated: 2026-04-18_

Parent phase: [#17 Phase 2: Circuit Runtime](https://github.com/devlikebear/wirecraft/issues/17)

## Phase Goal

Add the smallest deterministic server-side digital circuit runtime. The server should understand circuit-capable blocks in the voxel world, extract a circuit graph, evaluate High/Low/Unknown state every tick, and expose enough snapshot data for the client to visualize starter circuits.

## Phase Non-goals

- Do not add physical actuators or kinematic movement. That starts in Phase 3.
- Do not simulate analog voltage, current, resistance, or voltage drop.
- Do not implement Arduino/C++ runtime compatibility.
- Do not add database-backed persistence.

---

## Work Order 15: Add Circuit Block Types and Metadata

Status: Completed. GitHub issue: [#18](https://github.com/devlikebear/wirecraft/issues/18).

## Goal

Add the first Phase 2 circuit domain foundation: circuit-capable block types and minimal metadata helpers that identify how each block participates in digital circuits.

## Non-goals

- Do not implement graph extraction yet.
- Do not evaluate signal state yet.
- Do not change client toolbar/rendering yet.
- Do not add button input commands yet.

## Touch points (<=5)

- `internal/world/block.go`
- `internal/world/world_test.go`
- `internal/circuit/types.go`
- `internal/circuit/types_test.go`
- `docs/tasks/phase-2-work-orders.md`

## Steps

- [x] Add block types for power, wire, button, AND gate, and MCU output.
- [x] Add circuit-facing metadata helpers for block role and pin requirements.
- [x] Keep existing block numeric values stable for current protocol compatibility.
- [x] Add tests for circuit block validity and metadata completeness.
- [x] Keep Phase 1 block behavior unchanged.

## Acceptance Criteria

- [x] `go test ./internal/world/...` passes.
- [x] `go test ./internal/circuit/...` passes.
- [x] `go test ./...` passes.
- [x] No frontend files are touched.

## Verification Commands

- `go test ./internal/world/...`
- `go test ./internal/circuit/...`
- `go test ./...`

---

## Work Order 16: Add Circuit Graph Primitives

Status: Completed. GitHub issue: [#19](https://github.com/devlikebear/wirecraft/issues/19).

## Goal

Implement a deterministic in-memory graph model for circuit nodes, pins, edges, and signal states.

## Non-goals

- Do not extract from voxel world yet.
- Do not evaluate signal propagation yet.
- Do not add snapshots or client UI.

## Touch points (<=5)

- `internal/circuit/graph.go`
- `internal/circuit/graph_test.go`
- `internal/circuit/types.go`
- `docs/tasks/phase-2-work-orders.md`

## Steps

- [x] Add `NodeID`, `NodeType`, `SignalState`, `PinID`, and edge structures.
- [x] Add deterministic node/edge sorting helpers.
- [x] Add validation for duplicate nodes and missing edge endpoints.
- [x] Add tests for graph construction and deterministic ordering.

## Acceptance Criteria

- [x] `go test ./internal/circuit/...` passes.
- [x] No world, sim, server, or frontend files are touched.

## Verification Commands

- `go test ./internal/circuit/...`

---

## Work Order 17: Extract Circuit Graph from World Blocks

Status: Completed. GitHub issue: [#20](https://github.com/devlikebear/wirecraft/issues/20).

## Goal

Convert circuit-capable voxel blocks into a circuit graph using adjacency rules.

## Non-goals

- Do not evaluate signal state yet.
- Do not implement oriented gate pin routing beyond the minimal metadata from WO-15.
- Do not expose snapshot state yet.

## Touch points (<=5)

- `internal/circuit/extract.go`
- `internal/circuit/extract_test.go`
- `internal/circuit/graph.go`
- `internal/world/world.go`
- `docs/tasks/phase-2-work-orders.md`

## Steps

- [x] Iterate occupied world blocks and include only circuit-capable blocks.
- [x] Create graph nodes with stable IDs derived from block position.
- [x] Connect adjacent wire/power/button/output blocks by deterministic rules.
- [x] Add tests for power-wire-output and disconnected circuits.

## Acceptance Criteria

- [x] `go test ./internal/circuit/...` passes.
- [x] `go test ./internal/world/...` passes.

## Verification Commands

- `go test ./internal/circuit/...`
- `go test ./internal/world/...`

---

## Work Order 18: Evaluate Digital Signal State

Status: Completed. GitHub issue: [#21](https://github.com/devlikebear/wirecraft/issues/21).

## Goal

Evaluate High/Low/Unknown digital signal state from a circuit graph deterministically.

## Non-goals

- Do not integrate button commands yet.
- Do not integrate simulation snapshots yet.
- Do not implement analog behavior.

## Touch points (<=5)

- `internal/circuit/evaluate.go`
- `internal/circuit/evaluate_test.go`
- `internal/circuit/graph.go`
- `internal/circuit/types.go`
- `docs/tasks/phase-2-work-orders.md`

## Steps

- [x] Add propagation from power sources through wires.
- [x] Add button default/off state behavior.
- [x] Add AND gate truth table behavior.
- [x] Handle cycles without panic by resolving to stable state or Unknown.
- [x] Add deterministic truth table and cycle tests.

## Acceptance Criteria

- [x] `go test ./internal/circuit/...` passes.

## Verification Commands

- `go test ./internal/circuit/...`

---

## Work Order 19: Add Circuit State to Simulation Snapshots

Status: Completed. GitHub issue: [#22](https://github.com/devlikebear/wirecraft/issues/22).

## Goal

Run circuit extraction/evaluation during simulation steps and include circuit debug state in authoritative snapshots.

## Non-goals

- Do not add frontend visualization yet.
- Do not add button input commands yet.
- Do not optimize snapshot size.

## Touch points (<=5)

- `internal/sim/simulation.go`
- `internal/sim/snapshot.go`
- `internal/sim/simulation_test.go`
- `internal/netproto/snapshot.go`
- `docs/tasks/phase-2-work-orders.md`

## Steps

- [x] Add a circuit evaluation step after world command application and before snapshot build.
- [x] Add snapshot payload for block position signal state.
- [x] Add tests proving snapshots include deterministic circuit state.
- [x] Keep Phase 1 blocks/entities unchanged.

## Acceptance Criteria

- [x] `go test ./internal/sim/...` passes.
- [x] `go test ./internal/netproto/...` passes.
- [x] `go test ./...` passes.

## Verification Commands

- `go test ./internal/sim/...`
- `go test ./internal/netproto/...`
- `go test ./...`

---

## Work Order 20: Add Button Input Command

Status: Completed. GitHub issue: [#23](https://github.com/devlikebear/wirecraft/issues/23).

## Goal

Add a server-authoritative command for button press/release state and feed that state into circuit evaluation.

## Non-goals

- Do not add toolbar UI yet.
- Do not add actuator behavior.
- Do not persist button state across server restarts.

## Touch points (<=5)

- `internal/netproto/command.go`
- `internal/netproto/command_test.go`
- `internal/sim/simulation.go`
- `internal/sim/simulation_test.go`
- `docs/tasks/phase-2-work-orders.md`

## Steps

- [x] Add `set_button` command type and validation.
- [x] Store button state by block position inside simulation state.
- [x] Feed button state into circuit evaluation.
- [x] Add tests for press/release affecting output signal.

## Acceptance Criteria

- [x] `go test ./internal/netproto/...` passes.
- [x] `go test ./internal/sim/...` passes.
- [x] `go test ./...` passes.

## Verification Commands

- `go test ./internal/netproto/...`
- `go test ./internal/sim/...`
- `go test ./...`

---

## Work Order 21: Add Circuit Block Toolbar

Status: Completed. GitHub issue: [#24](https://github.com/devlikebear/wirecraft/issues/24).

## Goal

Let the browser place Phase 2 circuit block types through a compact toolbar without breaking existing voxel editing.

## Non-goals

- Do not implement circuit state visualization yet.
- Do not add inspect/component cards yet.
- Do not change server command validation except as required by existing block types.

## Touch points (<=5)

- `web/src/ui/Toolbar.ts`
- `web/src/ui/Toolbar.test.ts`
- `web/src/input/EditController.ts`
- `web/src/main.ts`
- `web/src/styles.css`

## Steps

- [x] Add toolbar buttons for solid, power, wire, button, AND gate, and MCU output.
- [x] Wire selected block type into existing place command flow.
- [x] Keep remove behavior unchanged.
- [x] Add focused tests for selected block type and command payload.

## Acceptance Criteria

- [x] `cd web && npm test` passes.
- [x] `cd web && npm run build` passes.
- [x] Browser can place each circuit block type.

## Verification Commands

- `cd web && npm test`
- `cd web && npm run build`

---

## Work Order 22: Visualize Circuit Signal State

Status: Completed. GitHub issue: [#25](https://github.com/devlikebear/wirecraft/issues/25).

## Goal

Render High/Low/Unknown circuit state in the Three.js scene using authoritative snapshot data.

## Non-goals

- Do not add inspect/component cards yet.
- Do not add actuator movement.
- Do not hand-roll a complex electrical visualization.

## Touch points (<=5)

- `web/src/net/protocol.ts`
- `web/src/render/CircuitOverlay.ts`
- `web/src/render/CircuitOverlay.test.ts`
- `web/src/main.ts`
- `web/src/styles.css`

## Steps

- [x] Parse circuit state from snapshots.
- [x] Add a lightweight visual overlay or material cue for High/Low/Unknown states.
- [x] Keep static voxel rendering behavior unchanged.
- [x] Add tests for snapshot parsing and render item creation.

## Acceptance Criteria

- [x] `cd web && npm test` passes.
- [x] `cd web && npm run build` passes.
- [x] Browser shows signal state changes from authoritative snapshots.

## Verification Commands

- `cd web && npm test`
- `cd web && npm run build`

---

## Work Order 23: Add Starter Component Cards

Status: Completed. GitHub issue: [#26](https://github.com/devlikebear/wirecraft/issues/26).

## Goal

Add beginner-friendly component card data and a simple inspect panel for starter circuit parts.

## Non-goals

- Do not add a full documentation system.
- Do not add AI assistant or code generation.
- Do not add advanced component search.

## Touch points (<=5)

- `internal/component/card.go`
- `internal/component/card_test.go`
- `web/src/state/componentCards.ts`
- `web/src/ui/InspectPanel.ts`
- `docs/reference/component-cards.md`

## Steps

- [x] Define component card schema with pins, role, wiring notes, warnings, and simplification notes.
- [x] Add starter cards for power, wire, button, AND gate, MCU output, LED, and resistor.
- [x] Add a compact inspect panel that can show selected component text.
- [x] Add validation tests for required card fields.

## Acceptance Criteria

- [x] `go test ./internal/component/...` passes.
- [x] `cd web && npm test` passes.
- [x] `cd web && npm run build` passes.
- [x] Browser can display at least one starter component card.

## Verification Commands

- `go test ./internal/component/...`
- `cd web && npm test`
- `cd web && npm run build`

---

## Phase 2 Checkpoint: Verify Circuit Runtime

GitHub issue: [#27](https://github.com/devlikebear/wirecraft/issues/27)

## Goal

Verify the smallest end-to-end digital circuit loop before moving to Phase 3.

## Steps

- [ ] Place power, wire, button, AND gate, and MCU output blocks.
- [ ] Verify server evaluates signal state deterministically each tick.
- [ ] Verify button input changes snapshot signal state.
- [ ] Verify browser visualizes High/Low/Unknown state.
- [ ] Verify starter component cards are readable in the inspect panel.

## Acceptance Criteria

- [ ] `go test ./...` passes.
- [ ] `cd web && npm test` passes.
- [ ] `cd web && npm run build` passes.
- [ ] Browser smoke test passes for power -> wire -> output.
- [ ] Browser smoke test passes for button-gated output.
- [ ] User approves moving to Phase 3.

## Verification Commands

- `go test ./...`
- `cd web && npm test`
- `cd web && npm run build`
