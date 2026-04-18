# WireCraft — Current Status

_Last updated: 2026-04-18_

## Purpose

This file is the handoff point for new Codex sessions. Start here, then read the roadmap and the active work order file.

## Canonical Documents

- Product and MVP: [`../plans/wire-craft-prd.md`](../plans/wire-craft-prd.md)
- Full roadmap: [`../plans/wire-craft-roadmap.md`](../plans/wire-craft-roadmap.md)
- Research notes: [`../plans/wire-craft-research-notes.md`](../plans/wire-craft-research-notes.md)
- Active phase plan: [`../plans/wire-craft-phase-1-authoritative-voxel-loop.md`](../plans/wire-craft-phase-1-authoritative-voxel-loop.md)
- Active work orders: [`phase-1-work-orders.md`](./phase-1-work-orders.md)

## Repository State

- GitHub repo: `https://github.com/devlikebear/wirecraft`
- Default branch: `main`
- Go module target: `github.com/devlikebear/wirecraft`
- Frontend package manager: `npm`
- Deployment shape: development uses separate Go/Vite servers; release build embeds Vite output into the Go binary with `go:embed`.

## Current Phase

**Phase 1: Authoritative Voxel Loop**

Goal: create the smallest server-authoritative 3D editing loop. The Go server owns world state and snapshots; the TypeScript/Three.js client renders snapshots and sends commands only.

GitHub issue: [#1 Phase 1: Authoritative Voxel Loop](https://github.com/devlikebear/wirecraft/issues/1)

## Completed Work

- [x] [#2 WO-1: Scaffold Go project and embedded web boundary](https://github.com/devlikebear/wirecraft/issues/2)
- [x] [#3 WO-2: Add tick clock and voxel world core](https://github.com/devlikebear/wirecraft/issues/3)
- [x] [#4 WO-3: Initialize Vite Three.js client skeleton](https://github.com/devlikebear/wirecraft/issues/4)
- [x] [#5 WO-4: Add command and snapshot protocol types](https://github.com/devlikebear/wirecraft/issues/5)
- [x] [#6 WO-5: Build simulation snapshot from world state](https://github.com/devlikebear/wirecraft/issues/6)
- [x] [#7 WO-6: Add in-memory simulation runner](https://github.com/devlikebear/wirecraft/issues/7)

## Next Work Order

Start with **WO-7** in [`phase-1-work-orders.md`](./phase-1-work-orders.md).

GitHub issue: [#8 WO-7: Add WebSocket simulation stream](https://github.com/devlikebear/wirecraft/issues/8)

## GitHub Issue Index

- [#1 Phase 1: Authoritative Voxel Loop](https://github.com/devlikebear/wirecraft/issues/1)
- [#2 WO-1: Scaffold Go project and embedded web boundary](https://github.com/devlikebear/wirecraft/issues/2)
- [#3 WO-2: Add tick clock and voxel world core](https://github.com/devlikebear/wirecraft/issues/3)
- [#4 WO-3: Initialize Vite Three.js client skeleton](https://github.com/devlikebear/wirecraft/issues/4)
- [#5 WO-4: Add command and snapshot protocol types](https://github.com/devlikebear/wirecraft/issues/5)
- [#6 WO-5: Build simulation snapshot from world state](https://github.com/devlikebear/wirecraft/issues/6)
- [#7 WO-6: Add in-memory simulation runner](https://github.com/devlikebear/wirecraft/issues/7)
- [#8 WO-7: Add WebSocket simulation stream](https://github.com/devlikebear/wirecraft/issues/8)

## Session Rules

- Before editing, run `git status --short --branch`.
- Prefer small commits using `feat/fix/chore` prefixes.
- Keep work test-first where possible.
- Update this file when a work order is completed or the next task changes.
- Update the matching GitHub issue when a work order is completed.
- Do not advance to the next phase until the phase checkpoint is complete and the user approves.

## Open Decisions

- WebSocket library: choose before WebSocket protocol implementation. Default preference is a small maintained package such as `github.com/coder/websocket`.
- Exact frontend test runner: choose during WO-3. Default preference is Vitest.
