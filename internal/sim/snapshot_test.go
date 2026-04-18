package sim

import (
	"testing"

	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/world"
)

func TestBuildSnapshotFromWorld(t *testing.T) {
	w := world.NewDefault()
	if err := w.Set(world.Position{X: 2, Y: 0, Z: 0}, world.BlockSolid); err != nil {
		t.Fatalf("Set(solid) error = %v, want nil", err)
	}
	if err := w.Set(world.Position{X: 1, Y: 0, Z: 0}, world.BlockDebugMover); err != nil {
		t.Fatalf("Set(debug mover) error = %v, want nil", err)
	}

	snapshot := BuildSnapshot(SnapshotInput{
		Tick:         TickID(9),
		ServerTimeMS: 1700000000999,
		World:        w,
		Stats: SnapshotStatsInput{
			ClientCount:        3,
			CommandQueueLength: 5,
		},
	})

	if snapshot.Tick != 9 {
		t.Fatalf("snapshot.Tick = %d, want 9", snapshot.Tick)
	}
	if snapshot.ServerTimeMS != 1700000000999 {
		t.Fatalf("snapshot.ServerTimeMS = %d, want 1700000000999", snapshot.ServerTimeMS)
	}

	wantBlocks := []netproto.BlockSnapshot{
		{Position: world.Position{X: 1, Y: 0, Z: 0}, BlockType: world.BlockDebugMover},
		{Position: world.Position{X: 2, Y: 0, Z: 0}, BlockType: world.BlockSolid},
	}
	if len(snapshot.Blocks) != len(wantBlocks) {
		t.Fatalf("len(snapshot.Blocks) = %d, want %d", len(snapshot.Blocks), len(wantBlocks))
	}
	for i := range wantBlocks {
		if snapshot.Blocks[i] != wantBlocks[i] {
			t.Fatalf("snapshot.Blocks[%d] = %+v, want %+v", i, snapshot.Blocks[i], wantBlocks[i])
		}
	}
	if len(snapshot.Entities) != 1 {
		t.Fatalf("len(snapshot.Entities) = %d, want 1", len(snapshot.Entities))
	}
	assertDebugMoverEntity(t, snapshot.Entities[0], netproto.Vec3{X: 2.125, Y: 1.25, Z: 1.5625})
	if snapshot.Stats.ClientCount != 3 {
		t.Fatalf("snapshot.Stats.ClientCount = %d, want 3", snapshot.Stats.ClientCount)
	}
	if snapshot.Stats.CommandQueueLength != 5 {
		t.Fatalf("snapshot.Stats.CommandQueueLength = %d, want 5", snapshot.Stats.CommandQueueLength)
	}
	if snapshot.Stats.SnapshotBytes <= 0 {
		t.Fatalf("snapshot.Stats.SnapshotBytes = %d, want positive", snapshot.Stats.SnapshotBytes)
	}
}

func TestBuildSnapshotIncludesDeterministicDebugEntity(t *testing.T) {
	first := BuildSnapshot(SnapshotInput{Tick: TickID(1)})
	second := BuildSnapshot(SnapshotInput{Tick: TickID(2)})

	if len(first.Entities) != 1 {
		t.Fatalf("len(first.Entities) = %d, want 1", len(first.Entities))
	}
	if len(second.Entities) != 1 {
		t.Fatalf("len(second.Entities) = %d, want 1", len(second.Entities))
	}

	assertDebugMoverEntity(t, first.Entities[0], netproto.Vec3{X: 1.125, Y: 1.25, Z: 1.0625})
	assertDebugMoverEntity(t, second.Entities[0], netproto.Vec3{X: 1.25, Y: 1.25, Z: 1.125})
	if first.Entities[0].Transform.Position == second.Entities[0].Transform.Position {
		t.Fatalf("debug entity position did not change: %+v", first.Entities[0].Transform.Position)
	}
}

func TestBuildSnapshotOmitsAirBlocks(t *testing.T) {
	w := world.NewDefault()
	pos := world.Position{X: 1, Y: 2, Z: 3}

	if err := w.Set(pos, world.BlockSolid); err != nil {
		t.Fatalf("Set(solid) error = %v, want nil", err)
	}
	if err := w.Remove(pos); err != nil {
		t.Fatalf("Remove() error = %v, want nil", err)
	}

	snapshot := BuildSnapshot(SnapshotInput{
		Tick:  TickID(1),
		World: w,
	})

	if len(snapshot.Blocks) != 0 {
		t.Fatalf("snapshot.Blocks = %+v, want empty", snapshot.Blocks)
	}
}

func assertDebugMoverEntity(t *testing.T, entity netproto.EntitySnapshot, wantPosition netproto.Vec3) {
	t.Helper()

	if entity.ID != netproto.EntityIDDebugMover {
		t.Fatalf("entity.ID = %q, want %q", entity.ID, netproto.EntityIDDebugMover)
	}
	if entity.Type != netproto.EntityTypeDebugMover {
		t.Fatalf("entity.Type = %q, want %q", entity.Type, netproto.EntityTypeDebugMover)
	}
	if entity.Transform.Position != wantPosition {
		t.Fatalf("entity.Transform.Position = %+v, want %+v", entity.Transform.Position, wantPosition)
	}
	if entity.Transform.Rotation != (netproto.Quat{X: 0, Y: 0, Z: 0, W: 1}) {
		t.Fatalf("entity.Transform.Rotation = %+v, want identity", entity.Transform.Rotation)
	}
	if entity.Transform.Scale != (netproto.Vec3{X: 0.5, Y: 0.5, Z: 0.5}) {
		t.Fatalf("entity.Transform.Scale = %+v, want half scale", entity.Transform.Scale)
	}
}
