package circuit

import (
	"testing"

	"github.com/devlikebear/wirecraft/internal/world"
)

func TestCircuitBlockTypes(t *testing.T) {
	want := []world.BlockType{
		world.BlockPower,
		world.BlockWire,
		world.BlockButton,
		world.BlockAndGate,
		world.BlockMCUOutput,
	}

	got := BlockTypes()
	if len(got) != len(want) {
		t.Fatalf("len(BlockTypes()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("BlockTypes()[%d] = %s, want %s", i, got[i], want[i])
		}
	}
}

func TestMetadataForBlock(t *testing.T) {
	cases := []struct {
		block       world.BlockType
		role        BlockRole
		inputPins   int
		outputPins  int
		conductive  bool
		interactive bool
	}{
		{world.BlockPower, RolePowerSource, 0, 1, false, false},
		{world.BlockWire, RoleConductor, 0, 0, true, false},
		{world.BlockButton, RoleSwitch, 1, 1, false, true},
		{world.BlockAndGate, RoleLogicGate, 2, 1, false, false},
		{world.BlockMCUOutput, RoleOutput, 1, 0, false, false},
	}

	for _, tc := range cases {
		metadata, ok := MetadataForBlock(tc.block)
		if !ok {
			t.Fatalf("MetadataForBlock(%s) ok = false, want true", tc.block)
		}
		if metadata.BlockType != tc.block {
			t.Fatalf("MetadataForBlock(%s).BlockType = %s, want %s", tc.block, metadata.BlockType, tc.block)
		}
		if metadata.Role != tc.role {
			t.Fatalf("MetadataForBlock(%s).Role = %s, want %s", tc.block, metadata.Role, tc.role)
		}
		if metadata.InputPins != tc.inputPins {
			t.Fatalf("MetadataForBlock(%s).InputPins = %d, want %d", tc.block, metadata.InputPins, tc.inputPins)
		}
		if metadata.OutputPins != tc.outputPins {
			t.Fatalf("MetadataForBlock(%s).OutputPins = %d, want %d", tc.block, metadata.OutputPins, tc.outputPins)
		}
		if metadata.Conductive != tc.conductive {
			t.Fatalf("MetadataForBlock(%s).Conductive = %t, want %t", tc.block, metadata.Conductive, tc.conductive)
		}
		if metadata.Interactive != tc.interactive {
			t.Fatalf("MetadataForBlock(%s).Interactive = %t, want %t", tc.block, metadata.Interactive, tc.interactive)
		}
		if metadata.DisplayName == "" {
			t.Fatalf("MetadataForBlock(%s).DisplayName is empty", tc.block)
		}
	}
}

func TestMetadataRejectsNonCircuitBlocks(t *testing.T) {
	for _, block := range []world.BlockType{world.BlockAir, world.BlockSolid, world.BlockDebugMover} {
		if IsCircuitBlock(block) {
			t.Fatalf("IsCircuitBlock(%s) = true, want false", block)
		}
		if metadata, ok := MetadataForBlock(block); ok {
			t.Fatalf("MetadataForBlock(%s) = %+v, true; want zero, false", block, metadata)
		}
	}
}
