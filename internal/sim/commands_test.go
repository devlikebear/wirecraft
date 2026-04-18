package sim

import (
	"reflect"
	"testing"

	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/world"
)

func TestOrderQueuedCommandsUsesStableTickOrderingKey(t *testing.T) {
	commands := []QueuedCommand{
		{Command: commandForClient("client-b", "cmd-b", 2, world.Position{X: 0, Y: 0, Z: 0}, world.BlockWire), ReceivedSequence: 4},
		{Command: commandForClient("client-a", "cmd-c", 1, world.Position{X: 0, Y: 0, Z: 0}, world.BlockPower), ReceivedSequence: 2},
		{Command: commandForClient("client-a", "cmd-a", 1, world.Position{X: 0, Y: 0, Z: 0}, world.BlockButton), ReceivedSequence: 2},
		{Command: commandForClient("client-b", "cmd-a", 1, world.Position{X: 0, Y: 0, Z: 0}, world.BlockSolid), ReceivedSequence: 1},
	}

	ordered := OrderQueuedCommands(commands)
	got := commandIDs(ordered)
	want := []string{
		"client-b/cmd-a",
		"client-a/cmd-a",
		"client-a/cmd-c",
		"client-b/cmd-b",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("OrderQueuedCommands IDs = %+v, want %+v", got, want)
	}
}

func TestSimulationApplyCommandsOrdersSameTickConflicts(t *testing.T) {
	simulation := NewSimulation()
	pos := world.Position{X: 1, Y: 1, Z: 1}

	errors := simulation.ApplyCommands([]QueuedCommand{
		{Command: commandForClient("client-2", "cmd-2", 10, pos, world.BlockPower), ReceivedSequence: 2},
		{Command: commandForClient("client-1", "cmd-1", 10, pos, world.BlockSolid), ReceivedSequence: 1},
	})

	if len(errors) != 0 {
		t.Fatalf("ApplyCommands errors = %+v, want none", errors)
	}
	snapshot := simulation.Step(StepInput{})
	assertSnapshotBlocks(t, snapshot, []netproto.BlockSnapshot{
		{Position: pos, BlockType: world.BlockPower},
	})
}

func TestSimulationIgnoresDuplicateCommandIDFromSameClient(t *testing.T) {
	simulation := NewSimulation()
	pos := world.Position{X: 1, Y: 1, Z: 1}

	errors := simulation.ApplyCommands([]QueuedCommand{
		{Command: commandForClient("client-1", "cmd-1", 10, pos, world.BlockSolid), ReceivedSequence: 1},
		{Command: commandForClient("client-1", "cmd-1", 10, pos, world.BlockPower), ReceivedSequence: 2},
	})

	if len(errors) != 0 {
		t.Fatalf("ApplyCommands errors = %+v, want none", errors)
	}
	snapshot := simulation.Step(StepInput{})
	assertSnapshotBlocks(t, snapshot, []netproto.BlockSnapshot{
		{Position: pos, BlockType: world.BlockSolid},
	})
}

func TestSimulationAllowsSameCommandIDFromDifferentClients(t *testing.T) {
	simulation := NewSimulation()
	pos := world.Position{X: 1, Y: 1, Z: 1}

	errors := simulation.ApplyCommands([]QueuedCommand{
		{Command: commandForClient("client-1", "cmd-1", 10, pos, world.BlockSolid), ReceivedSequence: 1},
		{Command: commandForClient("client-2", "cmd-1", 10, pos, world.BlockPower), ReceivedSequence: 2},
	})

	if len(errors) != 0 {
		t.Fatalf("ApplyCommands errors = %+v, want none", errors)
	}
	snapshot := simulation.Step(StepInput{})
	assertSnapshotBlocks(t, snapshot, []netproto.BlockSnapshot{
		{Position: pos, BlockType: world.BlockPower},
	})
}

func commandForClient(clientID string, commandID string, tickHint uint64, pos world.Position, blockType world.BlockType) netproto.Command {
	return netproto.Command{
		Type:      netproto.CommandPlaceBlock,
		ClientID:  clientID,
		CommandID: commandID,
		TickHint:  tickHint,
		Position:  pos,
		BlockType: blockType,
	}
}

func commandIDs(commands []QueuedCommand) []string {
	ids := make([]string, 0, len(commands))
	for _, command := range commands {
		ids = append(ids, command.Command.ClientID+"/"+command.Command.CommandID)
	}
	return ids
}
