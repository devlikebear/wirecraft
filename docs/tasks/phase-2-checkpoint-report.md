# WireCraft — Phase 2 Checkpoint Report

_Last updated: 2026-04-18_

GitHub issue: [#27 Phase 2 Checkpoint: Verify circuit runtime](https://github.com/devlikebear/wirecraft/issues/27)

## Summary

Phase 2 runtime verification is complete. The server-authoritative circuit loop can place circuit blocks from the browser, evaluate deterministic signal state on the Go server, stream snapshots over WebSocket, and visualize/read that state in the Three.js client.

The user approved moving to Phase 3 on 2026-04-18.

## Browser Smoke Scenario

Open the app at `http://127.0.0.1:5173/`, then run this scenario from the browser console or Playwright:

1. Place `Power`, `Wire`, and `MCU Output` in adjacent positions.
2. Wait for the MCU output circuit node to report `high`.
3. Place `Button`, `Wire`, and `MCU Output` in adjacent positions.
4. Wait for the button-gated output to report `low` while released.
5. Send `set_button` with `buttonPressed: true`.
6. Wait for the button-gated output to report `high`.
7. Send `set_button` with `buttonPressed: false`.
8. Wait for the button-gated output to report `low`.
9. Confirm the inspect panel displays a starter component card.

## Verified Result

- WebSocket status: `open`
- Power -> Wire -> MCU Output: `high`
- Button released -> Wire -> MCU Output: `low`
- Button pressed -> Wire -> MCU Output: `high`
- Button released again -> Wire -> MCU Output: `low`
- Power + Power -> AND Gate -> MCU Output: gate `high`, output `high`
- Inspect panel includes the starter `Power` card.

## Verification Commands

- `go test ./...`
- `cd web && npm test`
- `cd web && npm run build`
