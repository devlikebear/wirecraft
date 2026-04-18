package sim

import (
	"encoding/json"

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
