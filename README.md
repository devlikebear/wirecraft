# WireCraft

WireCraft is an online maker sandbox for building circuits, mechanisms, and small robotic systems in a shared 3D workspace. The project starts with a server-authoritative voxel and circuit simulation, then grows toward a Reality Bridge that can turn virtual builds into parts lists, wiring guides, controller code, and 3D-printable model candidates.

## Planning

- [PRD](docs/plans/wire-craft-prd.md)
- [Development Roadmap](docs/plans/wire-craft-roadmap.md)
- [Research Notes](docs/plans/wire-craft-research-notes.md)

## Task Tracking

- [Current Status](docs/tasks/current-status.md)
- [Phase 1 Work Orders](docs/tasks/phase-1-work-orders.md)
- [GitHub Issue: Phase 1](https://github.com/devlikebear/wirecraft/issues/1)
- [GitHub Issues](https://github.com/devlikebear/wirecraft/issues)

## Current Status

Phase 1 implementation is in progress.

- Completed: WO-1 scaffolded the Go server and embedded web boundary.
- Next: WO-2 adds the fixed tick clock and voxel world core.

## Development

Run the Go server:

```sh
go run ./cmd/wirecraft-server
```

The server listens on `127.0.0.1:8080` by default.

Health check:

```sh
curl http://127.0.0.1:8080/healthz
```

Run tests:

```sh
go test ./...
```

## Web UI Packaging

During development, the TypeScript/Three.js app will run through Vite. For release builds, Vite output will be copied into `internal/webui/dist` and embedded in the Go binary with `go:embed`.
