package sim

import (
	"reflect"
	"testing"

	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/world"
)

func TestBuildChangedSetSnapshotIncludesChangedAndRemovedBlocks(t *testing.T) {
	base := netproto.Snapshot{
		Mode: netproto.SnapshotModeFull,
		Tick: 4,
		Blocks: []netproto.BlockSnapshot{
			{Position: world.Position{X: 0, Y: 0, Z: 0}, BlockType: world.BlockSolid},
			{Position: world.Position{X: 1, Y: 0, Z: 0}, BlockType: world.BlockWire},
			{Position: world.Position{X: 2, Y: 0, Z: 0}, BlockType: world.BlockPower},
		},
		Entities: []netproto.EntitySnapshot{
			entitySnapshot("debug-mover-1", netproto.EntityTypeDebugMover, netproto.Vec3{X: 1, Y: 1, Z: 1}),
			entitySnapshot("piston:2:0:0", netproto.EntityTypePiston, netproto.Vec3{X: 2, Y: 0.5, Z: 0}),
		},
	}
	next := netproto.Snapshot{
		Mode:         netproto.SnapshotModeFull,
		Tick:         5,
		ServerTimeMS: 1700000000500,
		Blocks: []netproto.BlockSnapshot{
			{Position: world.Position{X: 0, Y: 0, Z: 0}, BlockType: world.BlockSolid},
			{Position: world.Position{X: 2, Y: 0, Z: 0}, BlockType: world.BlockButton},
			{Position: world.Position{X: 3, Y: 0, Z: 0}, BlockType: world.BlockMCUOutput},
		},
		Entities: []netproto.EntitySnapshot{
			entitySnapshot("debug-mover-1", netproto.EntityTypeDebugMover, netproto.Vec3{X: 1, Y: 1, Z: 1}),
			entitySnapshot("piston:2:0:0", netproto.EntityTypePiston, netproto.Vec3{X: 3, Y: 0.5, Z: 0}),
			entitySnapshot("motor:4:0:0", netproto.EntityTypeMotor, netproto.Vec3{X: 4, Y: 0.5, Z: 0}),
		},
	}

	changedSet := BuildChangedSetSnapshot(base, next)

	if changedSet.Mode != netproto.SnapshotModeChangedSet {
		t.Fatalf("changedSet.Mode = %q, want %q", changedSet.Mode, netproto.SnapshotModeChangedSet)
	}
	if changedSet.BaseTick != 4 {
		t.Fatalf("changedSet.BaseTick = %d, want 4", changedSet.BaseTick)
	}
	if len(changedSet.Blocks) != 0 {
		t.Fatalf("changedSet.Blocks = %+v, want empty full block payload", changedSet.Blocks)
	}
	if len(changedSet.Entities) != 0 {
		t.Fatalf("changedSet.Entities = %+v, want empty full entity payload", changedSet.Entities)
	}

	wantChangedBlocks := []netproto.BlockSnapshot{
		{Position: world.Position{X: 2, Y: 0, Z: 0}, BlockType: world.BlockButton},
		{Position: world.Position{X: 3, Y: 0, Z: 0}, BlockType: world.BlockMCUOutput},
	}
	if !reflect.DeepEqual(changedSet.ChangedBlocks, wantChangedBlocks) {
		t.Fatalf("changedSet.ChangedBlocks = %+v, want %+v", changedSet.ChangedBlocks, wantChangedBlocks)
	}

	wantRemovedBlocks := []world.Position{{X: 1, Y: 0, Z: 0}}
	if !reflect.DeepEqual(changedSet.RemovedBlocks, wantRemovedBlocks) {
		t.Fatalf("changedSet.RemovedBlocks = %+v, want %+v", changedSet.RemovedBlocks, wantRemovedBlocks)
	}

	wantChangedEntities := []netproto.EntitySnapshot{
		entitySnapshot("motor:4:0:0", netproto.EntityTypeMotor, netproto.Vec3{X: 4, Y: 0.5, Z: 0}),
		entitySnapshot("piston:2:0:0", netproto.EntityTypePiston, netproto.Vec3{X: 3, Y: 0.5, Z: 0}),
	}
	if !reflect.DeepEqual(changedSet.ChangedEntities, wantChangedEntities) {
		t.Fatalf("changedSet.ChangedEntities = %+v, want %+v", changedSet.ChangedEntities, wantChangedEntities)
	}
	if changedSet.Stats.SnapshotBytes <= 0 {
		t.Fatalf("changedSet.Stats.SnapshotBytes = %d, want positive", changedSet.Stats.SnapshotBytes)
	}
}

func entitySnapshot(id string, entityType string, position netproto.Vec3) netproto.EntitySnapshot {
	return netproto.EntitySnapshot{
		ID:   id,
		Type: entityType,
		Transform: netproto.TransformSnapshot{
			Position: position,
			Rotation: netproto.Quat{W: 1},
			Scale:    netproto.Vec3{X: 1, Y: 1, Z: 1},
		},
	}
}
