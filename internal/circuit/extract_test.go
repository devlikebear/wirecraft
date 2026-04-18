package circuit

import (
	"testing"

	"github.com/devlikebear/wirecraft/internal/world"
)

func TestExtractGraphFromWorldBuildsPowerWireOutputGraph(t *testing.T) {
	w := world.New(world.Dimensions{X: 8, Y: 8, Z: 8})
	mustSetBlock(t, w, world.Position{X: 0, Y: 0, Z: 0}, world.BlockPower)
	mustSetBlock(t, w, world.Position{X: 1, Y: 0, Z: 0}, world.BlockWire)
	mustSetBlock(t, w, world.Position{X: 2, Y: 0, Z: 0}, world.BlockMCUOutput)
	mustSetBlock(t, w, world.Position{X: 1, Y: 1, Z: 0}, world.BlockSolid)
	mustSetBlock(t, w, world.Position{X: 2, Y: 1, Z: 0}, world.BlockDebugMover)

	graph, err := ExtractGraphFromWorld(w)
	if err != nil {
		t.Fatalf("ExtractGraphFromWorld() error = %v, want nil", err)
	}

	assertNodes(t, graph.Nodes(), []Node{
		{ID: NodeID("0:0:0"), Type: NodeTypePowerSource},
		{ID: NodeID("1:0:0"), Type: NodeTypeWire},
		{ID: NodeID("2:0:0"), Type: NodeTypeMCUOutput},
	})
	assertEdges(t, graph.Edges(), []Edge{
		{
			From: PinRef{NodeID: NodeID("0:0:0"), PinID: PinID("body")},
			To:   PinRef{NodeID: NodeID("1:0:0"), PinID: PinID("body")},
		},
		{
			From: PinRef{NodeID: NodeID("1:0:0"), PinID: PinID("body")},
			To:   PinRef{NodeID: NodeID("2:0:0"), PinID: PinID("body")},
		},
	})
}

func TestExtractGraphFromWorldKeepsDisconnectedCircuitsWithoutEdges(t *testing.T) {
	w := world.New(world.Dimensions{X: 8, Y: 8, Z: 8})
	mustSetBlock(t, w, world.Position{X: 0, Y: 0, Z: 0}, world.BlockPower)
	mustSetBlock(t, w, world.Position{X: 3, Y: 0, Z: 0}, world.BlockWire)
	mustSetBlock(t, w, world.Position{X: 0, Y: 3, Z: 0}, world.BlockButton)
	mustSetBlock(t, w, world.Position{X: 0, Y: 0, Z: 3}, world.BlockAndGate)

	graph, err := ExtractGraphFromWorld(w)
	if err != nil {
		t.Fatalf("ExtractGraphFromWorld() error = %v, want nil", err)
	}

	assertNodes(t, graph.Nodes(), []Node{
		{ID: NodeID("0:0:0"), Type: NodeTypePowerSource},
		{ID: NodeID("0:0:3"), Type: NodeTypeAndGate},
		{ID: NodeID("0:3:0"), Type: NodeTypeButton},
		{ID: NodeID("3:0:0"), Type: NodeTypeWire},
	})
	assertEdges(t, graph.Edges(), nil)
}

func TestExtractGraphFromWorldConnectsVerticalAdjacentBlocksDeterministically(t *testing.T) {
	w := world.New(world.Dimensions{X: 8, Y: 8, Z: 8})
	mustSetBlock(t, w, world.Position{X: 2, Y: 2, Z: 2}, world.BlockWire)
	mustSetBlock(t, w, world.Position{X: 2, Y: 2, Z: 3}, world.BlockButton)

	graph, err := ExtractGraphFromWorld(w)
	if err != nil {
		t.Fatalf("ExtractGraphFromWorld() error = %v, want nil", err)
	}

	assertEdges(t, graph.Edges(), []Edge{
		{
			From: PinRef{NodeID: NodeID("2:2:2"), PinID: PinID("body")},
			To:   PinRef{NodeID: NodeID("2:2:3"), PinID: PinID("body")},
		},
	})
}

func TestNodeIDForPosition(t *testing.T) {
	got := NodeIDForPosition(world.Position{X: 12, Y: 3, Z: 4})
	if got != NodeID("12:3:4") {
		t.Fatalf("NodeIDForPosition() = %q, want %q", got, NodeID("12:3:4"))
	}
}

func mustSetBlock(t *testing.T, w *world.World, pos world.Position, block world.BlockType) {
	t.Helper()

	if err := w.Set(pos, block); err != nil {
		t.Fatalf("Set(%+v, %s) error = %v, want nil", pos, block, err)
	}
}
