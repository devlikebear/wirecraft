package circuit

import "testing"

func TestEvaluateGraphPropagatesPowerThroughWiresToOutput(t *testing.T) {
	graph := mustNewGraph(t,
		[]Node{
			{ID: NodeID("power"), Type: NodeTypePowerSource},
			{ID: NodeID("wire:a"), Type: NodeTypeWire},
			{ID: NodeID("wire:b"), Type: NodeTypeWire},
			{ID: NodeID("out"), Type: NodeTypeMCUOutput},
		},
		[]Edge{
			bodyEdge("power", "wire:a"),
			bodyEdge("wire:a", "wire:b"),
			bodyEdge("wire:b", "out"),
		},
	)

	evaluation := EvaluateGraph(graph)

	assertSignal(t, evaluation, NodeID("power"), SignalHigh)
	assertSignal(t, evaluation, NodeID("wire:a"), SignalHigh)
	assertSignal(t, evaluation, NodeID("wire:b"), SignalHigh)
	assertSignal(t, evaluation, NodeID("out"), SignalHigh)
}

func TestEvaluateGraphDefaultsButtonOff(t *testing.T) {
	graph := mustNewGraph(t,
		[]Node{
			{ID: NodeID("button"), Type: NodeTypeButton},
			{ID: NodeID("wire"), Type: NodeTypeWire},
			{ID: NodeID("out"), Type: NodeTypeMCUOutput},
		},
		[]Edge{
			bodyEdge("button", "wire"),
			bodyEdge("wire", "out"),
		},
	)

	evaluation := EvaluateGraph(graph)

	assertSignal(t, evaluation, NodeID("button"), SignalLow)
	assertSignal(t, evaluation, NodeID("wire"), SignalLow)
	assertSignal(t, evaluation, NodeID("out"), SignalLow)
}

func TestEvaluateGraphEvaluatesAndGateTruthTable(t *testing.T) {
	cases := []struct {
		name  string
		left  NodeType
		right NodeType
		want  SignalState
	}{
		{name: "high high", left: NodeTypePowerSource, right: NodeTypePowerSource, want: SignalHigh},
		{name: "high low", left: NodeTypePowerSource, right: NodeTypeButton, want: SignalLow},
		{name: "low low", left: NodeTypeButton, right: NodeTypeButton, want: SignalLow},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			graph := mustNewGraph(t,
				[]Node{
					{ID: NodeID("left"), Type: tc.left},
					{ID: NodeID("right"), Type: tc.right},
					{ID: NodeID("and"), Type: NodeTypeAndGate},
					{ID: NodeID("out"), Type: NodeTypeMCUOutput},
				},
				[]Edge{
					bodyEdge("left", "and"),
					bodyEdge("right", "and"),
					bodyEdge("and", "out"),
				},
			)

			evaluation := EvaluateGraph(graph)

			assertSignal(t, evaluation, NodeID("and"), tc.want)
			assertSignal(t, evaluation, NodeID("out"), tc.want)
		})
	}
}

func TestEvaluateGraphLeavesFloatingCycleUnknown(t *testing.T) {
	graph := mustNewGraph(t,
		[]Node{
			{ID: NodeID("wire:a"), Type: NodeTypeWire},
			{ID: NodeID("wire:b"), Type: NodeTypeWire},
			{ID: NodeID("out"), Type: NodeTypeMCUOutput},
		},
		[]Edge{
			bodyEdge("wire:a", "wire:b"),
			bodyEdge("wire:b", "out"),
			bodyEdge("out", "wire:a"),
		},
	)

	evaluation := EvaluateGraph(graph)

	assertSignal(t, evaluation, NodeID("wire:a"), SignalUnknown)
	assertSignal(t, evaluation, NodeID("wire:b"), SignalUnknown)
	assertSignal(t, evaluation, NodeID("out"), SignalUnknown)
}

func TestEvaluationReturnsStateCopies(t *testing.T) {
	graph := mustNewGraph(t,
		[]Node{{ID: NodeID("power"), Type: NodeTypePowerSource}},
		nil,
	)
	evaluation := EvaluateGraph(graph)

	states := evaluation.States()
	states[NodeID("power")] = SignalLow

	assertSignal(t, evaluation, NodeID("power"), SignalHigh)
	assertSignal(t, evaluation, NodeID("missing"), SignalUnknown)
}

func bodyEdge(from NodeID, to NodeID) Edge {
	return Edge{
		From: PinRef{NodeID: from, PinID: DefaultPinID},
		To:   PinRef{NodeID: to, PinID: DefaultPinID},
	}
}

func mustNewGraph(t *testing.T, nodes []Node, edges []Edge) Graph {
	t.Helper()

	graph, err := NewGraph(nodes, edges)
	if err != nil {
		t.Fatalf("NewGraph() error = %v, want nil", err)
	}
	return graph
}

func assertSignal(t *testing.T, evaluation Evaluation, nodeID NodeID, want SignalState) {
	t.Helper()

	if got := evaluation.State(nodeID); got != want {
		t.Fatalf("signal %s = %s, want %s", nodeID, got, want)
	}
}
