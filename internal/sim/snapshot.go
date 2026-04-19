package sim

import (
	"encoding/json"

	"github.com/devlikebear/wirecraft/internal/circuit"
	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/physics"
	"github.com/devlikebear/wirecraft/internal/world"
)

type SnapshotInput struct {
	Tick             TickID
	ServerTimeMS     int64
	World            *world.World
	ButtonStates     map[world.Position]bool
	ActuatorEntities []physics.DynamicEntity
	Presence         netproto.PresenceSnapshot
	CommandAcks      []netproto.CommandAckSnapshot
	Stats            SnapshotStatsInput
}

type SnapshotStatsInput struct {
	ClientCount        int
	CommandQueueLength int
}

func BuildSnapshot(input SnapshotInput) netproto.Snapshot {
	blocks := make([]netproto.BlockSnapshot, 0)
	if input.World != nil {
		occupied := input.World.OccupiedBlocks()
		blocks = make([]netproto.BlockSnapshot, 0, len(occupied))
		for _, block := range occupied {
			blocks = append(blocks, netproto.BlockSnapshot{
				Position:   block.Position,
				BlockType:  block.BlockType,
				Facing:     block.Facing,
				Properties: block.Properties.Clone(),
			})
		}
	}

	entities := append([]netproto.EntitySnapshot{buildDebugMoverEntity(input.Tick)}, buildActuatorEntities(input.ActuatorEntities)...)
	commandAcks := append([]netproto.CommandAckSnapshot(nil), input.CommandAcks...)
	if commandAcks == nil {
		commandAcks = []netproto.CommandAckSnapshot{}
	}

	snapshot := netproto.Snapshot{
		Mode:            netproto.SnapshotModeFull,
		Tick:            uint64(input.Tick),
		ServerTimeMS:    input.ServerTimeMS,
		Blocks:          blocks,
		ChangedBlocks:   []netproto.BlockSnapshot{},
		RemovedBlocks:   []world.Position{},
		Entities:        entities,
		ChangedEntities: []netproto.EntitySnapshot{},
		Circuit:         buildCircuitSnapshot(input.World, input.ButtonStates),
		Presence:        input.Presence,
		CommandAcks:     commandAcks,
		Stats: netproto.SnapshotStats{
			ClientCount:        input.Stats.ClientCount,
			CommandQueueLength: input.Stats.CommandQueueLength,
		},
	}

	return finalizeSnapshot(snapshot)
}

func finalizeSnapshot(snapshot netproto.Snapshot) netproto.Snapshot {
	encoded, err := json.Marshal(snapshot)
	if err == nil {
		snapshot.Stats.SnapshotBytes = len(encoded)
	}

	return snapshot
}

func buildCircuitSnapshot(w *world.World, buttonStates map[world.Position]bool) netproto.CircuitSnapshot {
	nodes := make([]netproto.CircuitNodeSnapshot, 0)
	if w == nil {
		return netproto.CircuitSnapshot{Nodes: nodes}
	}

	graph, err := circuit.ExtractGraphFromWorld(w)
	if err != nil {
		return netproto.CircuitSnapshot{Nodes: nodes}
	}

	positionsByNodeID := make(map[circuit.NodeID]world.Position)
	buttonStatesByNodeID := make(map[circuit.NodeID]circuit.SignalState)
	for _, block := range w.OccupiedBlocks() {
		if !circuit.IsCircuitBlock(block.BlockType) {
			continue
		}
		nodeID := circuit.NodeIDForPosition(block.Position)
		positionsByNodeID[nodeID] = block.Position
		if block.BlockType == world.BlockButton && buttonStates[block.Position] {
			buttonStatesByNodeID[nodeID] = circuit.SignalHigh
		}
	}

	evaluation := circuit.EvaluateGraphWithInput(graph, circuit.EvaluationInput{
		ButtonStates: buttonStatesByNodeID,
	})
	for _, node := range graph.Nodes() {
		position, ok := positionsByNodeID[node.ID]
		if !ok {
			continue
		}

		nodes = append(nodes, netproto.CircuitNodeSnapshot{
			Position:    position,
			NodeID:      string(node.ID),
			NodeType:    string(node.Type),
			SignalState: evaluation.State(node.ID).String(),
		})
	}

	return netproto.CircuitSnapshot{Nodes: nodes}
}

func buildDebugMoverEntity(tick TickID) netproto.EntitySnapshot {
	step := float64(tick % 16)

	return netproto.EntitySnapshot{
		ID:   netproto.EntityIDDebugMover,
		Type: netproto.EntityTypeDebugMover,
		Transform: netproto.TransformSnapshot{
			Position: netproto.Vec3{
				X: 1 + step/8,
				Y: 1.25,
				Z: 1 + step/16,
			},
			Rotation: netproto.Quat{X: 0, Y: 0, Z: 0, W: 1},
			Scale:    netproto.Vec3{X: 0.5, Y: 0.5, Z: 0.5},
		},
	}
}

func buildActuatorEntities(entities []physics.DynamicEntity) []netproto.EntitySnapshot {
	sorted := physics.SortEntities(entities)
	snapshots := make([]netproto.EntitySnapshot, 0, len(sorted))
	for _, entity := range sorted {
		snapshots = append(snapshots, netproto.EntitySnapshot{
			ID:        string(entity.ID),
			Type:      string(entity.Type),
			Transform: transformSnapshot(entity.Transform),
		})
	}
	return snapshots
}

func transformSnapshot(transform physics.Transform) netproto.TransformSnapshot {
	return netproto.TransformSnapshot{
		Position: netproto.Vec3{
			X: transform.Position.X,
			Y: transform.Position.Y,
			Z: transform.Position.Z,
		},
		Rotation: netproto.Quat{
			X: transform.Rotation.X,
			Y: transform.Rotation.Y,
			Z: transform.Rotation.Z,
			W: transform.Rotation.W,
		},
		Scale: netproto.Vec3{
			X: transform.Scale.X,
			Y: transform.Scale.Y,
			Z: transform.Scale.Z,
		},
	}
}
