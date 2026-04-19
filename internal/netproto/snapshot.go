package netproto

import "github.com/devlikebear/wirecraft/internal/world"

const (
	EntityIDDebugMover   = "debug-mover-1"
	EntityTypeDebugMover = "debug_mover"
	EntityTypePiston     = "piston"
	EntityTypeMotor      = "motor"
)

type SnapshotMode string

const (
	SnapshotModeFull       SnapshotMode = "full"
	SnapshotModeChangedSet SnapshotMode = "changed_set"
)

type CommandAckStatus string

const (
	CommandAckAccepted CommandAckStatus = "accepted"
	CommandAckRejected CommandAckStatus = "rejected"
)

type Snapshot struct {
	Mode            SnapshotMode         `json:"mode"`
	Tick            uint64               `json:"tick"`
	BaseTick        uint64               `json:"baseTick,omitempty"`
	ServerTimeMS    int64                `json:"serverTimeMs"`
	Blocks          []BlockSnapshot      `json:"blocks"`
	ChangedBlocks   []BlockSnapshot      `json:"changedBlocks"`
	RemovedBlocks   []world.Position     `json:"removedBlocks"`
	Entities        []EntitySnapshot     `json:"entities"`
	ChangedEntities []EntitySnapshot     `json:"changedEntities"`
	Circuit         CircuitSnapshot      `json:"circuit"`
	Presence        PresenceSnapshot     `json:"presence"`
	CommandAcks     []CommandAckSnapshot `json:"commandAcks"`
	Stats           SnapshotStats        `json:"stats"`
}

type BlockSnapshot struct {
	Position   world.Position        `json:"position"`
	BlockType  world.BlockType       `json:"blockType"`
	Facing     world.Facing          `json:"facing,omitempty"`
	Properties world.BlockProperties `json:"properties,omitempty"`
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

type PresenceSnapshot struct {
	Clients []ClientPresenceSnapshot `json:"clients"`
}

type ClientPresenceSnapshot struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type CommandAckSnapshot struct {
	ClientID  string           `json:"clientId"`
	CommandID string           `json:"commandId"`
	Status    CommandAckStatus `json:"status"`
	Reason    string           `json:"reason,omitempty"`
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
