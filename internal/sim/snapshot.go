package sim

import (
	"encoding/json"

	"github.com/devlikebear/wirecraft/internal/circuit"
	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/world"
)

type SnapshotInput struct {
	Tick         TickID
	ServerTimeMS int64
	World        *world.World
	Stats        SnapshotStatsInput
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
				Position:  block.Position,
				BlockType: block.BlockType,
			})
		}
	}

	snapshot := netproto.Snapshot{
		Tick:         uint64(input.Tick),
		ServerTimeMS: input.ServerTimeMS,
		Blocks:       blocks,
		Entities:     []netproto.EntitySnapshot{buildDebugMoverEntity(input.Tick)},
		Circuit:      buildCircuitSnapshot(input.World),
		Stats: netproto.SnapshotStats{
			ClientCount:        input.Stats.ClientCount,
			CommandQueueLength: input.Stats.CommandQueueLength,
		},
	}

	encoded, err := json.Marshal(snapshot)
	if err == nil {
		snapshot.Stats.SnapshotBytes = len(encoded)
	}

	return snapshot
}

func buildCircuitSnapshot(w *world.World) netproto.CircuitSnapshot {
	nodes := make([]netproto.CircuitNodeSnapshot, 0)
	if w == nil {
		return netproto.CircuitSnapshot{Nodes: nodes}
	}

	graph, err := circuit.ExtractGraphFromWorld(w)
	if err != nil {
		return netproto.CircuitSnapshot{Nodes: nodes}
	}

	positionsByNodeID := make(map[circuit.NodeID]world.Position)
	for _, block := range w.OccupiedBlocks() {
		if !circuit.IsCircuitBlock(block.BlockType) {
			continue
		}
		positionsByNodeID[circuit.NodeIDForPosition(block.Position)] = block.Position
	}

	evaluation := circuit.EvaluateGraph(graph)
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
