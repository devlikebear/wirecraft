# WireCraft — Current Status

_Last updated: 2026-04-19_

## Purpose

This file is the handoff point for new Codex sessions. Start here, then read the reset roadmap and active reset work orders.

## Canonical Documents

- Product and MVP reset PRD: [`../plans/wire-craft-prd.md`](../plans/wire-craft-prd.md)
- Reset roadmap: [`../plans/wire-craft-roadmap.md`](../plans/wire-craft-roadmap.md)
- Reset rationale: [`../plans/wire-craft-mvp-reset.md`](../plans/wire-craft-mvp-reset.md)
- Active reset work orders: [`mvp-reset-work-orders.md`](./mvp-reset-work-orders.md)
- Research notes: [`../plans/wire-craft-research-notes.md`](../plans/wire-craft-research-notes.md)

Historical phase plans remain useful as implementation context:

- Phase 1: [`../plans/wire-craft-phase-1-authoritative-voxel-loop.md`](../plans/wire-craft-phase-1-authoritative-voxel-loop.md)
- Phase 2: [`../plans/wire-craft-phase-2-circuit-runtime.md`](../plans/wire-craft-phase-2-circuit-runtime.md)
- Phase 3: [`../plans/wire-craft-phase-3-physical-actuators.md`](../plans/wire-craft-phase-3-physical-actuators.md)
- Phase 4 paused context: [`../plans/wire-craft-phase-4-multiplayer-physics-sync.md`](../plans/wire-craft-phase-4-multiplayer-physics-sync.md), [`phase-4-work-orders.md`](./phase-4-work-orders.md)

## Repository State

- GitHub repo: `https://github.com/devlikebear/wirecraft`
- Default branch: `main`
- Go module target: `github.com/devlikebear/wirecraft`
- Frontend package manager: `npm`
- Deployment shape: development uses separate Go/Vite servers; release build embeds Vite output into the Go binary with `go:embed`.

## Current Product Direction

**MVP Reset: Workbench Device Builder**

Tracking issue: [#49 MVP Reset: Workbench Device Builder](https://github.com/devlikebear/wirecraft/issues/49)

WireCraft is now scoped around a small workbench where a user can build, inspect, edit, and run one interactive device before broader multiplayer/export work resumes. The first target device is a **button-controlled sliding door**.

## Current Phase

**Phase R0: MVP Reset Planning**

Status: completed locally by the reset documentation update.

The next implementation phase is **Phase R1: Workbench Block Instance UX**.

## Next Work Order

Start with **Work Order R1: Add Block Instance Properties To Snapshots** in [`mvp-reset-work-orders.md`](./mvp-reset-work-orders.md).

GitHub issue: [#50 WO-R1: Add block instance properties to snapshots](https://github.com/devlikebear/wirecraft/issues/50)

Recommended first implementation behavior:

- Add minimal block instance properties to server world/snapshot models.
- Preserve existing place/remove behavior.
- Parse optional properties in TypeScript snapshots.
- Keep this as a small TDD change before UI property editing.

## Completed Foundation

The following completed work remains part of the foundation:

- Phase 1: authoritative voxel loop, WebSocket command/snapshot protocol, voxel rendering, raycast editing, interpolation.
- Phase 2: starter circuit block types, circuit graph/evaluation, button input, circuit snapshot, toolbar, signal overlay, component cards.
- Phase 3: dynamic entity/transform primitives, actuator model, piston/motor snapshots, actuator rendering, sensor input store, motor/driver cards.
- Phase 4 partial: room model, presence metadata, deterministic command ordering, command acknowledgements, changed-set snapshots, client delta application, viewport navigation, placement facing.

## Paused / Superseded Work

- [#39 Phase 4: Multiplayer Physics Sync](https://github.com/devlikebear/wirecraft/issues/39) is paused by the MVP reset. Completed foundation work remains valid.
- [#47 WO-39: Add basic actuator collision constraints](https://github.com/devlikebear/wirecraft/issues/47) is no longer the next task. Collision should return inside the sliding door vertical slice when blocked state has a product meaning.
- Phase 5 blueprint/export work is deferred until block instance metadata, port/slot rules, and the first device slice exist.

## Current Technical Reality

- Server world blocks currently store `Position`, `BlockType`, and `Facing`.
- Circuit extraction currently connects adjacent circuit blocks through a default `body` pin.
- Inspect UI currently shows toolbar-selected component cards, not selected placed block instances.
- Recent work added placement facing and viewport controls, but existing placed block rotation and editable properties are not implemented yet.

These gaps are now the intended next focus.

## GitHub Issue Index

- [#39 Phase 4: Multiplayer Physics Sync](https://github.com/devlikebear/wirecraft/issues/39) — paused context.
- [#47 WO-39: Add basic actuator collision constraints](https://github.com/devlikebear/wirecraft/issues/47) — superseded as immediate next work.
- [#49 MVP Reset: Workbench Device Builder](https://github.com/devlikebear/wirecraft/issues/49) — active reset parent.
- [#50 WO-R1: Add block instance properties to snapshots](https://github.com/devlikebear/wirecraft/issues/50) — next implementation task.

Historical completed issues are listed in GitHub and older phase work-order documents.

## Session Rules

- Before editing, run `git status --short --branch`.
- Prefer small commits using `feat/fix/chore` prefixes.
- Keep work test-first where possible.
- Update this file when a work order is completed or the next task changes.
- Update the matching GitHub issue when a work order is completed.
- Do not resume paused Phase 4/5 work until it is reconnected to the reset MVP.

## Open Decisions

- First device remains sliding door unless the user explicitly chooses a different vertical slice.
- Block instance properties should start minimal and backwards-compatible.
- Port/slot/attachment generalization should be added only as needed for the sliding door MVP.
