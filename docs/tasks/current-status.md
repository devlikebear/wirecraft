# WireCraft — Current Status

_Last updated: 2026-04-18_

## Purpose

This file is the handoff point for new Codex sessions. Start here, then read the roadmap and the active work order file.

## Canonical Documents

- Product and MVP: [`../plans/wire-craft-prd.md`](../plans/wire-craft-prd.md)
- Full roadmap: [`../plans/wire-craft-roadmap.md`](../plans/wire-craft-roadmap.md)
- Research notes: [`../plans/wire-craft-research-notes.md`](../plans/wire-craft-research-notes.md)
- Active phase plan: [`../plans/wire-craft-phase-3-physical-actuators.md`](../plans/wire-craft-phase-3-physical-actuators.md)
- Active work orders: [`phase-3-work-orders.md`](./phase-3-work-orders.md)
- Recent checkpoint report: [`phase-2-checkpoint-report.md`](./phase-2-checkpoint-report.md)

## Repository State

- GitHub repo: `https://github.com/devlikebear/wirecraft`
- Default branch: `main`
- Go module target: `github.com/devlikebear/wirecraft`
- Frontend package manager: `npm`
- Deployment shape: development uses separate Go/Vite servers; release build embeds Vite output into the Go binary with `go:embed`.

## Current Phase

**Phase 3: Physical Actuators**

Goal: connect server-authoritative circuit signal state to deterministic physical actuator motion. The server should turn High/Low circuit outputs into kinematic piston or motor transforms, and the client should render those transforms through the existing interpolation path.

GitHub issue: [#28 Phase 3: Physical Actuators](https://github.com/devlikebear/wirecraft/issues/28)

## Completed Work

- [x] [#2 WO-1: Scaffold Go project and embedded web boundary](https://github.com/devlikebear/wirecraft/issues/2)
- [x] [#3 WO-2: Add tick clock and voxel world core](https://github.com/devlikebear/wirecraft/issues/3)
- [x] [#4 WO-3: Initialize Vite Three.js client skeleton](https://github.com/devlikebear/wirecraft/issues/4)
- [x] [#5 WO-4: Add command and snapshot protocol types](https://github.com/devlikebear/wirecraft/issues/5)
- [x] [#6 WO-5: Build simulation snapshot from world state](https://github.com/devlikebear/wirecraft/issues/6)
- [x] [#7 WO-6: Add in-memory simulation runner](https://github.com/devlikebear/wirecraft/issues/7)
- [x] [#8 WO-7: Add WebSocket simulation stream](https://github.com/devlikebear/wirecraft/issues/8)
- [x] [#9 WO-8: Add TypeScript WebSocket client and snapshot store](https://github.com/devlikebear/wirecraft/issues/9)
- [x] [#10 WO-9: Render authoritative snapshots as voxels](https://github.com/devlikebear/wirecraft/issues/10)
- [x] [#11 WO-10: Add raycast block edit commands](https://github.com/devlikebear/wirecraft/issues/11)
- [x] [#12 WO-11: Add snapshot interpolation primitives](https://github.com/devlikebear/wirecraft/issues/12)
- [x] [#13 WO-12: Add dynamic debug entity snapshots](https://github.com/devlikebear/wirecraft/issues/13)
- [x] [#14 WO-13: Render dynamic debug entity with interpolation](https://github.com/devlikebear/wirecraft/issues/14)
- [x] [#15 WO-14: Add client debug overlay](https://github.com/devlikebear/wirecraft/issues/15)
- [x] [#16 Phase 1 checkpoint: Verify authoritative voxel loop](https://github.com/devlikebear/wirecraft/issues/16)
- [x] [#18 WO-15: Add circuit block types and metadata](https://github.com/devlikebear/wirecraft/issues/18)
- [x] [#19 WO-16: Add circuit graph primitives](https://github.com/devlikebear/wirecraft/issues/19)
- [x] [#20 WO-17: Extract circuit graph from world blocks](https://github.com/devlikebear/wirecraft/issues/20)
- [x] [#21 WO-18: Evaluate digital signal state](https://github.com/devlikebear/wirecraft/issues/21)
- [x] [#22 WO-19: Add circuit state to simulation snapshots](https://github.com/devlikebear/wirecraft/issues/22)
- [x] [#23 WO-20: Add button input command](https://github.com/devlikebear/wirecraft/issues/23)
- [x] [#24 WO-21: Add circuit block toolbar](https://github.com/devlikebear/wirecraft/issues/24)
- [x] [#25 WO-22: Visualize circuit signal state](https://github.com/devlikebear/wirecraft/issues/25)
- [x] [#26 WO-23: Add starter component cards](https://github.com/devlikebear/wirecraft/issues/26)
- [x] [#27 Phase 2 Checkpoint: Verify circuit runtime](https://github.com/devlikebear/wirecraft/issues/27)
- [x] [#29 WO-24: Add dynamic entity and transform primitives](https://github.com/devlikebear/wirecraft/issues/29)
- [x] [#30 WO-25: Add actuator component model](https://github.com/devlikebear/wirecraft/issues/30)

## Next Work Order

Start with **Work Order 26: Add Actuator Block Types** in [`phase-3-work-orders.md`](./phase-3-work-orders.md).

GitHub issue: [#31 WO-26: Add actuator block types](https://github.com/devlikebear/wirecraft/issues/31)

Phase 2 is complete and approved. Phase 3 has started; the physics transform/entity foundation and actuator domain model are complete. The next step is adding actuator-facing block/protocol types.

## GitHub Issue Index

- [#1 Phase 1: Authoritative Voxel Loop](https://github.com/devlikebear/wirecraft/issues/1)
- [#2 WO-1: Scaffold Go project and embedded web boundary](https://github.com/devlikebear/wirecraft/issues/2)
- [#3 WO-2: Add tick clock and voxel world core](https://github.com/devlikebear/wirecraft/issues/3)
- [#4 WO-3: Initialize Vite Three.js client skeleton](https://github.com/devlikebear/wirecraft/issues/4)
- [#5 WO-4: Add command and snapshot protocol types](https://github.com/devlikebear/wirecraft/issues/5)
- [#6 WO-5: Build simulation snapshot from world state](https://github.com/devlikebear/wirecraft/issues/6)
- [#7 WO-6: Add in-memory simulation runner](https://github.com/devlikebear/wirecraft/issues/7)
- [#8 WO-7: Add WebSocket simulation stream](https://github.com/devlikebear/wirecraft/issues/8)
- [#9 WO-8: Add TypeScript WebSocket client and snapshot store](https://github.com/devlikebear/wirecraft/issues/9)
- [#10 WO-9: Render authoritative snapshots as voxels](https://github.com/devlikebear/wirecraft/issues/10)
- [#11 WO-10: Add raycast block edit commands](https://github.com/devlikebear/wirecraft/issues/11)
- [#12 WO-11: Add snapshot interpolation primitives](https://github.com/devlikebear/wirecraft/issues/12)
- [#13 WO-12: Add dynamic debug entity snapshots](https://github.com/devlikebear/wirecraft/issues/13)
- [#14 WO-13: Render dynamic debug entity with interpolation](https://github.com/devlikebear/wirecraft/issues/14)
- [#15 WO-14: Add client debug overlay](https://github.com/devlikebear/wirecraft/issues/15)
- [#16 Phase 1 checkpoint: Verify authoritative voxel loop](https://github.com/devlikebear/wirecraft/issues/16)
- [#17 Phase 2: Circuit Runtime](https://github.com/devlikebear/wirecraft/issues/17)
- [#18 WO-15: Add circuit block types and metadata](https://github.com/devlikebear/wirecraft/issues/18)
- [#19 WO-16: Add circuit graph primitives](https://github.com/devlikebear/wirecraft/issues/19)
- [#20 WO-17: Extract circuit graph from world blocks](https://github.com/devlikebear/wirecraft/issues/20)
- [#21 WO-18: Evaluate digital signal state](https://github.com/devlikebear/wirecraft/issues/21)
- [#22 WO-19: Add circuit state to simulation snapshots](https://github.com/devlikebear/wirecraft/issues/22)
- [#23 WO-20: Add button input command](https://github.com/devlikebear/wirecraft/issues/23)
- [#24 WO-21: Add circuit block toolbar](https://github.com/devlikebear/wirecraft/issues/24)
- [#25 WO-22: Visualize circuit signal state](https://github.com/devlikebear/wirecraft/issues/25)
- [#26 WO-23: Add starter component cards](https://github.com/devlikebear/wirecraft/issues/26)
- [#27 Phase 2 Checkpoint: Verify circuit runtime](https://github.com/devlikebear/wirecraft/issues/27)
- [#28 Phase 3: Physical Actuators](https://github.com/devlikebear/wirecraft/issues/28)
- [#29 WO-24: Add dynamic entity and transform primitives](https://github.com/devlikebear/wirecraft/issues/29)
- [#30 WO-25: Add actuator component model](https://github.com/devlikebear/wirecraft/issues/30)
- [#31 WO-26: Add actuator block types](https://github.com/devlikebear/wirecraft/issues/31)

## Session Rules

- Before editing, run `git status --short --branch`.
- Prefer small commits using `feat/fix/chore` prefixes.
- Keep work test-first where possible.
- Update this file when a work order is completed or the next task changes.
- Update the matching GitHub issue when a work order is completed.
- Do not advance to the next phase until the phase checkpoint is complete and the user approves.

## Open Decisions

- WebSocket library: chosen for WO-7 as `github.com/coder/websocket v1.8.12` to preserve the Go 1.22 module target.
- Exact frontend test runner: choose during WO-3. Default preference is Vitest.
