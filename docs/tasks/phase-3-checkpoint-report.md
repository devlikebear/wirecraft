# WireCraft — Phase 3 Checkpoint Report

_Last updated: 2026-04-18_

GitHub issue: [#38 Phase 3 Checkpoint: Verify physical actuators](https://github.com/devlikebear/wirecraft/issues/38)

## Summary

Phase 3 runtime verification is complete. The server-authoritative circuit loop can drive a kinematic piston actuator from button signal state, export actuator transforms in snapshots, and render interpolated actuator entities in the Three.js client. The browser UI can place actuator blocks, and the inspect panel now shows beginner-facing motor/driver warnings.

Phase 4 has not started yet. User approval is required before moving to Phase 4.

## Browser Smoke Scenario

Open the app at `http://127.0.0.1:5173/`, then run this scenario from Playwright or the browser console:

1. Place `Button`, `Wire`, `Piston`, and `Motor` blocks in adjacent positions.
2. Wait for the released button path to report `low`.
3. Confirm piston and motor entities are present in snapshots and rendered by the client.
4. Send `set_button` with `buttonPressed: true`.
5. Wait for the wire node to report `high` and the piston entity to extend.
6. Send `set_button` with `buttonPressed: false`.
7. Wait for the wire node to report `low` and the piston entity to retract.
8. Select `Motor` in the toolbar and confirm the inspect panel warning is visible.
9. Open a second browser tab and confirm it receives the same piston entity state.

## Verified Result

- WebSocket status: `open`
- Button released -> Wire -> Piston input: `low`
- Piston released position: `x = 2`
- Button pressed -> Wire -> Piston input: `high`
- Piston extended position observed: `x = 2.8000000000000007`
- Button released again -> Wire -> Piston input: `low`
- Piston returned position: `x = 2`
- Client rendered entities: `3` (`debug_mover`, `piston`, `motor`)
- Motor inspect card warning: `MCU GPIO pins cannot drive a motor directly.`
- Two browser tabs both received `piston:2:0:6` at `x = 2`.

## Verification Commands

- `go test ./...`
- `cd web && npm test`
- `cd web && npm run build`

## Notes

- `cd web && npm run build` still emits the existing Vite chunk size warning.
- The browser smoke placed runtime-only test blocks around `z = 6`; restarting the local server clears that in-memory world state.
