package sim

import (
	"errors"
	"testing"

	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/world"
)

func TestSimulationApplyCommandPlacesAndRemovesBlocks(t *testing.T) {
	simulation := NewSimulation()
	pos := world.Position{X: 1, Y: 2, Z: 3}

	if err := simulation.ApplyCommand(placeCommand("cmd-1", pos, world.BlockSolid)); err != nil {
		t.Fatalf("ApplyCommand(place) error = %v, want nil", err)
	}

	placed := simulation.Step(StepInput{ServerTimeMS: 1000})
	assertSnapshotBlocks(t, placed, []netproto.BlockSnapshot{
		{Position: pos, BlockType: world.BlockSolid},
	})

	if err := simulation.ApplyCommand(removeCommand("cmd-2", pos)); err != nil {
		t.Fatalf("ApplyCommand(remove) error = %v, want nil", err)
	}

	removed := simulation.Step(StepInput{ServerTimeMS: 1050})
	assertSnapshotBlocks(t, removed, nil)
}

func TestSimulationStepAdvancesTickAndIncludesStats(t *testing.T) {
	simulation := NewSimulation()

	first := simulation.Step(StepInput{
		ServerTimeMS: 1700000000001,
		Stats: SnapshotStatsInput{
			ClientCount:        2,
			CommandQueueLength: 0,
		},
	})
	second := simulation.Step(StepInput{ServerTimeMS: 1700000000051})

	if first.Tick != 1 {
		t.Fatalf("first.Tick = %d, want 1", first.Tick)
	}
	if second.Tick != 2 {
		t.Fatalf("second.Tick = %d, want 2", second.Tick)
	}
	if first.ServerTimeMS != 1700000000001 {
		t.Fatalf("first.ServerTimeMS = %d, want 1700000000001", first.ServerTimeMS)
	}
	if first.Stats.ClientCount != 2 {
		t.Fatalf("first.Stats.ClientCount = %d, want 2", first.Stats.ClientCount)
	}
	if first.Stats.SnapshotBytes <= 0 {
		t.Fatalf("first.Stats.SnapshotBytes = %d, want positive", first.Stats.SnapshotBytes)
	}
}

func TestSimulationRejectsInvalidCommandWithoutChangingWorld(t *testing.T) {
	simulation := NewSimulation()
	validPos := world.Position{X: 1, Y: 1, Z: 1}
	invalid := placeCommand("cmd-2", world.Position{X: 99, Y: 1, Z: 1}, world.BlockSolid)

	if err := simulation.ApplyCommand(placeCommand("cmd-1", validPos, world.BlockDebugMover)); err != nil {
		t.Fatalf("ApplyCommand(valid) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(invalid); !errors.Is(err, world.ErrOutOfBounds) {
		t.Fatalf("ApplyCommand(invalid) error = %v, want %v", err, world.ErrOutOfBounds)
	}

	snapshot := simulation.Step(StepInput{})
	assertSnapshotBlocks(t, snapshot, []netproto.BlockSnapshot{
		{Position: validPos, BlockType: world.BlockDebugMover},
	})
}

func TestSimulationRejectsUnknownCommandWithoutChangingWorld(t *testing.T) {
	simulation := NewSimulation()
	pos := world.Position{X: 2, Y: 2, Z: 2}

	if err := simulation.ApplyCommand(placeCommand("cmd-1", pos, world.BlockSolid)); err != nil {
		t.Fatalf("ApplyCommand(valid) error = %v, want nil", err)
	}

	unknown := placeCommand("cmd-2", pos, world.BlockDebugMover)
	unknown.Type = netproto.CommandType("warp_block")
	if err := simulation.ApplyCommand(unknown); !errors.Is(err, netproto.ErrUnknownCommandType) {
		t.Fatalf("ApplyCommand(unknown) error = %v, want %v", err, netproto.ErrUnknownCommandType)
	}

	snapshot := simulation.Step(StepInput{})
	assertSnapshotBlocks(t, snapshot, []netproto.BlockSnapshot{
		{Position: pos, BlockType: world.BlockSolid},
	})
}

func placeCommand(commandID string, pos world.Position, blockType world.BlockType) netproto.Command {
	return netproto.Command{
		Type:      netproto.CommandPlaceBlock,
		ClientID:  "client-1",
		CommandID: commandID,
		Position:  pos,
		BlockType: blockType,
	}
}

func removeCommand(commandID string, pos world.Position) netproto.Command {
	return netproto.Command{
		Type:      netproto.CommandRemoveBlock,
		ClientID:  "client-1",
		CommandID: commandID,
		Position:  pos,
		BlockType: world.BlockAir,
	}
}

func assertSnapshotBlocks(t *testing.T, snapshot netproto.Snapshot, want []netproto.BlockSnapshot) {
	t.Helper()

	if len(snapshot.Blocks) != len(want) {
		t.Fatalf("len(snapshot.Blocks) = %d, want %d: %+v", len(snapshot.Blocks), len(want), snapshot.Blocks)
	}
	for i := range want {
		if snapshot.Blocks[i] != want[i] {
			t.Fatalf("snapshot.Blocks[%d] = %+v, want %+v", i, snapshot.Blocks[i], want[i])
		}
	}
}
