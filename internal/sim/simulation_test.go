package sim

import (
	"errors"
	"testing"

	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/physics"
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

func TestSimulationStepIncludesCircuitStateAfterCommands(t *testing.T) {
	simulation := NewSimulationWithDimensions(world.Dimensions{X: 8, Y: 8, Z: 8})

	if err := simulation.ApplyCommand(placeCommand("cmd-1", world.Position{X: 0, Y: 0, Z: 0}, world.BlockPower)); err != nil {
		t.Fatalf("ApplyCommand(power) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-2", world.Position{X: 1, Y: 0, Z: 0}, world.BlockWire)); err != nil {
		t.Fatalf("ApplyCommand(wire) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-3", world.Position{X: 2, Y: 0, Z: 0}, world.BlockMCUOutput)); err != nil {
		t.Fatalf("ApplyCommand(output) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-4", world.Position{X: 3, Y: 0, Z: 0}, world.BlockSolid)); err != nil {
		t.Fatalf("ApplyCommand(solid) error = %v, want nil", err)
	}

	snapshot := simulation.Step(StepInput{})

	assertCircuitNodes(t, snapshot.Circuit.Nodes, []netproto.CircuitNodeSnapshot{
		{
			Position:    world.Position{X: 0, Y: 0, Z: 0},
			NodeID:      "0:0:0",
			NodeType:    "power_source",
			SignalState: "high",
		},
		{
			Position:    world.Position{X: 1, Y: 0, Z: 0},
			NodeID:      "1:0:0",
			NodeType:    "wire",
			SignalState: "high",
		},
		{
			Position:    world.Position{X: 2, Y: 0, Z: 0},
			NodeID:      "2:0:0",
			NodeType:    "mcu_output",
			SignalState: "high",
		},
	})
}

func TestSimulationSetButtonCommandAffectsCircuitState(t *testing.T) {
	simulation := NewSimulationWithDimensions(world.Dimensions{X: 8, Y: 8, Z: 8})
	buttonPos := world.Position{X: 0, Y: 0, Z: 0}
	wirePos := world.Position{X: 1, Y: 0, Z: 0}
	outputPos := world.Position{X: 2, Y: 0, Z: 0}

	if err := simulation.ApplyCommand(placeCommand("cmd-1", buttonPos, world.BlockButton)); err != nil {
		t.Fatalf("ApplyCommand(button) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-2", wirePos, world.BlockWire)); err != nil {
		t.Fatalf("ApplyCommand(wire) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-3", outputPos, world.BlockMCUOutput)); err != nil {
		t.Fatalf("ApplyCommand(output) error = %v, want nil", err)
	}

	released := simulation.Step(StepInput{})
	assertCircuitSignal(t, released, "0:0:0", "low")
	assertCircuitSignal(t, released, "1:0:0", "low")
	assertCircuitSignal(t, released, "2:0:0", "low")

	if err := simulation.ApplyCommand(setButtonCommand("cmd-4", buttonPos, true)); err != nil {
		t.Fatalf("ApplyCommand(press) error = %v, want nil", err)
	}

	pressed := simulation.Step(StepInput{})
	assertCircuitSignal(t, pressed, "0:0:0", "high")
	assertCircuitSignal(t, pressed, "1:0:0", "high")
	assertCircuitSignal(t, pressed, "2:0:0", "high")

	if err := simulation.ApplyCommand(setButtonCommand("cmd-5", buttonPos, false)); err != nil {
		t.Fatalf("ApplyCommand(release) error = %v, want nil", err)
	}

	releasedAgain := simulation.Step(StepInput{})
	assertCircuitSignal(t, releasedAgain, "0:0:0", "low")
	assertCircuitSignal(t, releasedAgain, "1:0:0", "low")
	assertCircuitSignal(t, releasedAgain, "2:0:0", "low")
}

func TestSimulationUpdatesPistonAfterCircuitEvaluation(t *testing.T) {
	simulation := NewSimulationWithDimensions(world.Dimensions{X: 8, Y: 8, Z: 8})
	pistonPos := world.Position{X: 2, Y: 0, Z: 0}

	if err := simulation.ApplyCommand(placeCommand("cmd-1", world.Position{X: 0, Y: 0, Z: 0}, world.BlockPower)); err != nil {
		t.Fatalf("ApplyCommand(power) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-2", world.Position{X: 1, Y: 0, Z: 0}, world.BlockWire)); err != nil {
		t.Fatalf("ApplyCommand(wire) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-3", pistonPos, world.BlockPiston)); err != nil {
		t.Fatalf("ApplyCommand(piston) error = %v, want nil", err)
	}

	snapshot := simulation.Step(StepInput{DeltaSeconds: 0.25})

	assertActuatorEntities(t, simulation.ActuatorEntities(), []physics.DynamicEntity{
		{
			ID:   physics.EntityID("piston:2:0:0"),
			Type: physics.EntityTypePiston,
			Transform: physics.Transform{
				Position: physics.Vec3{X: 3, Y: 0.5, Z: 0},
				Rotation: physics.IdentityQuat(),
				Scale:    physics.UnitVec3(),
			},
			Target: physics.Transform{
				Position: physics.Vec3{X: 3, Y: 0.5, Z: 0},
				Rotation: physics.IdentityQuat(),
				Scale:    physics.UnitVec3(),
			},
		},
	})
	if len(snapshot.Entities) != 1 || snapshot.Entities[0].ID != netproto.EntityIDDebugMover {
		t.Fatalf("snapshot.Entities = %+v, want existing debug mover only", snapshot.Entities)
	}
}

func TestSimulationRetractsPistonWhenButtonSignalFallsLow(t *testing.T) {
	simulation := NewSimulationWithDimensions(world.Dimensions{X: 8, Y: 8, Z: 8})
	buttonPos := world.Position{X: 0, Y: 0, Z: 0}
	pistonPos := world.Position{X: 2, Y: 0, Z: 0}

	if err := simulation.ApplyCommand(placeCommand("cmd-1", buttonPos, world.BlockButton)); err != nil {
		t.Fatalf("ApplyCommand(button) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-2", world.Position{X: 1, Y: 0, Z: 0}, world.BlockWire)); err != nil {
		t.Fatalf("ApplyCommand(wire) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-3", pistonPos, world.BlockPiston)); err != nil {
		t.Fatalf("ApplyCommand(piston) error = %v, want nil", err)
	}

	simulation.Step(StepInput{DeltaSeconds: 0.25})
	assertActuatorPosition(t, simulation.ActuatorEntities(), "piston:2:0:0", physics.Vec3{X: 2, Y: 0.5, Z: 0})

	if err := simulation.ApplyCommand(setButtonCommand("cmd-4", buttonPos, true)); err != nil {
		t.Fatalf("ApplyCommand(press) error = %v, want nil", err)
	}
	simulation.Step(StepInput{DeltaSeconds: 0.25})
	assertActuatorPosition(t, simulation.ActuatorEntities(), "piston:2:0:0", physics.Vec3{X: 3, Y: 0.5, Z: 0})

	if err := simulation.ApplyCommand(setButtonCommand("cmd-5", buttonPos, false)); err != nil {
		t.Fatalf("ApplyCommand(release) error = %v, want nil", err)
	}
	simulation.Step(StepInput{DeltaSeconds: 0.25})
	assertActuatorPosition(t, simulation.ActuatorEntities(), "piston:2:0:0", physics.Vec3{X: 2, Y: 0.5, Z: 0})
}

func TestSimulationUsesServerTimeDeltaForActuatorStep(t *testing.T) {
	simulation := NewSimulationWithDimensions(world.Dimensions{X: 8, Y: 8, Z: 8})

	if err := simulation.ApplyCommand(placeCommand("cmd-1", world.Position{X: 0, Y: 0, Z: 0}, world.BlockPower)); err != nil {
		t.Fatalf("ApplyCommand(power) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-2", world.Position{X: 1, Y: 0, Z: 0}, world.BlockWire)); err != nil {
		t.Fatalf("ApplyCommand(wire) error = %v, want nil", err)
	}
	if err := simulation.ApplyCommand(placeCommand("cmd-3", world.Position{X: 2, Y: 0, Z: 0}, world.BlockPiston)); err != nil {
		t.Fatalf("ApplyCommand(piston) error = %v, want nil", err)
	}

	simulation.Step(StepInput{ServerTimeMS: 1000})
	simulation.Step(StepInput{ServerTimeMS: 1250})

	assertActuatorPosition(t, simulation.ActuatorEntities(), "piston:2:0:0", physics.Vec3{X: 3, Y: 0.5, Z: 0})
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

func setButtonCommand(commandID string, pos world.Position, pressed bool) netproto.Command {
	return netproto.Command{
		Type:          netproto.CommandSetButton,
		ClientID:      "client-1",
		CommandID:     commandID,
		Position:      pos,
		ButtonPressed: pressed,
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

func assertCircuitSignal(t *testing.T, snapshot netproto.Snapshot, nodeID string, want string) {
	t.Helper()

	for _, node := range snapshot.Circuit.Nodes {
		if node.NodeID == nodeID {
			if node.SignalState != want {
				t.Fatalf("node %s signal = %s, want %s", nodeID, node.SignalState, want)
			}
			return
		}
	}
	t.Fatalf("node %s not found in circuit snapshot: %+v", nodeID, snapshot.Circuit.Nodes)
}

func assertCircuitNodes(t *testing.T, got []netproto.CircuitNodeSnapshot, want []netproto.CircuitNodeSnapshot) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(circuit nodes) = %d, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("circuit nodes[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func assertActuatorEntities(t *testing.T, got []physics.DynamicEntity, want []physics.DynamicEntity) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(actuator entities) = %d, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i].ID != want[i].ID {
			t.Fatalf("actuator[%d].ID = %q, want %q", i, got[i].ID, want[i].ID)
		}
		if got[i].Type != want[i].Type {
			t.Fatalf("actuator[%d].Type = %q, want %q", i, got[i].Type, want[i].Type)
		}
		if got[i].Transform != want[i].Transform {
			t.Fatalf("actuator[%d].Transform = %+v, want %+v", i, got[i].Transform, want[i].Transform)
		}
		if got[i].Target != want[i].Target {
			t.Fatalf("actuator[%d].Target = %+v, want %+v", i, got[i].Target, want[i].Target)
		}
	}
}

func assertActuatorPosition(t *testing.T, got []physics.DynamicEntity, id string, want physics.Vec3) {
	t.Helper()

	for _, entity := range got {
		if entity.ID == physics.EntityID(id) {
			if entity.Transform.Position != want {
				t.Fatalf("actuator %s position = %+v, want %+v", id, entity.Transform.Position, want)
			}
			return
		}
	}
	t.Fatalf("actuator %s not found: %+v", id, got)
}
