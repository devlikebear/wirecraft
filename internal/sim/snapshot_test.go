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
	if len(snapshot.Entities) != 0 {
		t.Fatalf("len(snapshot.Entities) = %d, want 0", len(snapshot.Entities))
	}
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
