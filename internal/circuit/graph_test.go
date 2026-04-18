package circuit

import (
	"errors"
	"testing"
)

func TestNewGraphSortsNodesAndEdgesDeterministically(t *testing.T) {
	graph, err := NewGraph(
		[]Node{
			{ID: NodeID("wire:1"), Type: NodeTypeWire},
			{ID: NodeID("power:1"), Type: NodeTypePowerSource},
			{ID: NodeID("output:1"), Type: NodeTypeMCUOutput},
		},
		[]Edge{
			{
				From: PinRef{NodeID: NodeID("wire:1"), PinID: PinID("out")},
				To:   PinRef{NodeID: NodeID("output:1"), PinID: PinID("in")},
			},
			{
				From: PinRef{NodeID: NodeID("power:1"), PinID: PinID("out")},
				To:   PinRef{NodeID: NodeID("wire:1"), PinID: PinID("in")},
			},
		},
	)
	if err != nil {
		t.Fatalf("NewGraph() error = %v, want nil", err)
	}

	wantNodes := []Node{
		{ID: NodeID("output:1"), Type: NodeTypeMCUOutput},
		{ID: NodeID("power:1"), Type: NodeTypePowerSource},
		{ID: NodeID("wire:1"), Type: NodeTypeWire},
	}
	assertNodes(t, graph.Nodes(), wantNodes)

	wantEdges := []Edge{
		{
			From: PinRef{NodeID: NodeID("power:1"), PinID: PinID("out")},
			To:   PinRef{NodeID: NodeID("wire:1"), PinID: PinID("in")},
		},
		{
			From: PinRef{NodeID: NodeID("wire:1"), PinID: PinID("out")},
			To:   PinRef{NodeID: NodeID("output:1"), PinID: PinID("in")},
		},
	}
	assertEdges(t, graph.Edges(), wantEdges)
}

func TestNewGraphReturnsCopies(t *testing.T) {
	graph, err := NewGraph([]Node{{ID: NodeID("wire:1"), Type: NodeTypeWire}}, nil)
	if err != nil {
		t.Fatalf("NewGraph() error = %v, want nil", err)
	}

	nodes := graph.Nodes()
	nodes[0].ID = NodeID("mutated")

	if got := graph.Nodes()[0].ID; got != NodeID("wire:1") {
		t.Fatalf("graph node was mutated through copy: %s", got)
	}
}

func TestNewGraphRejectsDuplicateNodes(t *testing.T) {
	_, err := NewGraph(
		[]Node{
			{ID: NodeID("wire:1"), Type: NodeTypeWire},
			{ID: NodeID("wire:1"), Type: NodeTypeWire},
		},
		nil,
	)

	if !errors.Is(err, ErrDuplicateNode) {
		t.Fatalf("NewGraph(duplicate) error = %v, want %v", err, ErrDuplicateNode)
	}
}

func TestNewGraphRejectsMissingEdgeEndpoints(t *testing.T) {
	_, err := NewGraph(
		[]Node{{ID: NodeID("wire:1"), Type: NodeTypeWire}},
		[]Edge{
			{
				From: PinRef{NodeID: NodeID("wire:1"), PinID: PinID("out")},
				To:   PinRef{NodeID: NodeID("missing"), PinID: PinID("in")},
			},
		},
	)

	if !errors.Is(err, ErrMissingEdgeEndpoint) {
		t.Fatalf("NewGraph(missing endpoint) error = %v, want %v", err, ErrMissingEdgeEndpoint)
	}
}

func TestNodeTypeForBlockRole(t *testing.T) {
	cases := []struct {
		role BlockRole
		want NodeType
	}{
		{RolePowerSource, NodeTypePowerSource},
		{RoleConductor, NodeTypeWire},
		{RoleSwitch, NodeTypeButton},
		{RoleLogicGate, NodeTypeAndGate},
		{RoleOutput, NodeTypeMCUOutput},
	}

	for _, tc := range cases {
		got, ok := NodeTypeForBlockRole(tc.role)
		if !ok {
			t.Fatalf("NodeTypeForBlockRole(%s) ok = false, want true", tc.role)
		}
		if got != tc.want {
			t.Fatalf("NodeTypeForBlockRole(%s) = %s, want %s", tc.role, got, tc.want)
		}
	}
}

func TestSignalStateString(t *testing.T) {
	if SignalUnknown.String() != "unknown" {
		t.Fatalf("SignalUnknown.String() = %q, want unknown", SignalUnknown.String())
	}
	if SignalLow.String() != "low" {
		t.Fatalf("SignalLow.String() = %q, want low", SignalLow.String())
	}
	if SignalHigh.String() != "high" {
		t.Fatalf("SignalHigh.String() = %q, want high", SignalHigh.String())
	}
}

func assertNodes(t *testing.T, got []Node, want []Node) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(nodes) = %d, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("nodes[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func assertEdges(t *testing.T, got []Edge, want []Edge) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(edges) = %d, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("edges[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}
