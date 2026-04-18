package circuit

type Evaluation struct {
	states map[NodeID]SignalState
}

func EvaluateGraph(graph Graph) Evaluation {
	nodes := graph.Nodes()
	adjacency := buildAdjacency(graph.Edges())
	states := initialStates(nodes)

	maxIterations := len(nodes) + 1
	for range maxIterations {
		next := make(map[NodeID]SignalState, len(states))
		changed := false

		for _, node := range nodes {
			state := evaluateNode(node, neighborStates(adjacency[node.ID], states))
			next[node.ID] = state
			if state != states[node.ID] {
				changed = true
			}
		}

		states = next
		if !changed {
			break
		}
	}

	return Evaluation{states: states}
}

func (e Evaluation) State(nodeID NodeID) SignalState {
	state, ok := e.states[nodeID]
	if !ok {
		return SignalUnknown
	}
	return state
}

func (e Evaluation) States() map[NodeID]SignalState {
	states := make(map[NodeID]SignalState, len(e.states))
	for nodeID, state := range e.states {
		states[nodeID] = state
	}
	return states
}

func buildAdjacency(edges []Edge) map[NodeID][]NodeID {
	adjacency := make(map[NodeID][]NodeID)
	for _, edge := range edges {
		adjacency[edge.From.NodeID] = append(adjacency[edge.From.NodeID], edge.To.NodeID)
		adjacency[edge.To.NodeID] = append(adjacency[edge.To.NodeID], edge.From.NodeID)
	}
	return adjacency
}

func initialStates(nodes []Node) map[NodeID]SignalState {
	states := make(map[NodeID]SignalState, len(nodes))
	for _, node := range nodes {
		states[node.ID] = initialStateForNode(node)
	}
	return states
}

func initialStateForNode(node Node) SignalState {
	switch node.Type {
	case NodeTypePowerSource:
		return SignalHigh
	case NodeTypeButton:
		return SignalLow
	default:
		return SignalUnknown
	}
}

func neighborStates(neighborIDs []NodeID, states map[NodeID]SignalState) []SignalState {
	neighbors := make([]SignalState, 0, len(neighborIDs))
	for _, neighborID := range neighborIDs {
		state, ok := states[neighborID]
		if !ok {
			state = SignalUnknown
		}
		neighbors = append(neighbors, state)
	}
	return neighbors
}

func evaluateNode(node Node, neighbors []SignalState) SignalState {
	switch node.Type {
	case NodeTypePowerSource:
		return SignalHigh
	case NodeTypeButton:
		return SignalLow
	case NodeTypeWire, NodeTypeMCUOutput:
		return relayState(neighbors)
	case NodeTypeAndGate:
		return andGateState(neighbors)
	default:
		return SignalUnknown
	}
}

func relayState(neighbors []SignalState) SignalState {
	hasLow := false
	for _, state := range neighbors {
		switch state {
		case SignalHigh:
			return SignalHigh
		case SignalLow:
			hasLow = true
		}
	}
	if hasLow {
		return SignalLow
	}
	return SignalUnknown
}

func andGateState(neighbors []SignalState) SignalState {
	highCount := 0
	unknownCount := 0
	hasLow := false

	for _, state := range neighbors {
		switch state {
		case SignalHigh:
			highCount++
		case SignalLow:
			hasLow = true
		default:
			unknownCount++
		}
	}

	if highCount >= 2 {
		return SignalHigh
	}
	if hasLow {
		return SignalLow
	}
	if highCount+unknownCount >= 2 {
		return SignalUnknown
	}
	return SignalLow
}
