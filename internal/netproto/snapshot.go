package netproto

import "github.com/devlikebear/wirecraft/internal/world"

const (
	EntityIDDebugMover   = "debug-mover-1"
	EntityTypeDebugMover = "debug_mover"
	EntityTypePiston     = "piston"
	EntityTypeMotor      = "motor"
)

type Snapshot struct {
	Tick         uint64           `json:"tick"`
	ServerTimeMS int64            `json:"serverTimeMs"`
	Blocks       []BlockSnapshot  `json:"blocks"`
	Entities     []EntitySnapshot `json:"entities"`
	Circuit      CircuitSnapshot  `json:"circuit"`
	Stats        SnapshotStats    `json:"stats"`
}

type BlockSnapshot struct {
	Position  world.Position  `json:"position"`
	BlockType world.BlockType `json:"blockType"`
}

type EntitySnapshot struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Transform TransformSnapshot `json:"transform"`
}

type CircuitSnapshot struct {
	Nodes []CircuitNodeSnapshot `json:"nodes"`
}

type CircuitNodeSnapshot struct {
	Position    world.Position `json:"position"`
	NodeID      string         `json:"nodeId"`
	NodeType    string         `json:"nodeType"`
	SignalState string         `json:"signalState"`
}

type TransformSnapshot struct {
	Position Vec3 `json:"position"`
	Rotation Quat `json:"rotation"`
	Scale    Vec3 `json:"scale"`
}

type Vec3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Quat struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

type SnapshotStats struct {
	ClientCount        int `json:"clientCount"`
	CommandQueueLength int `json:"commandQueueLength"`
	SnapshotBytes      int `json:"snapshotBytes"`
}
