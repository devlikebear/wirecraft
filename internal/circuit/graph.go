package circuit

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrDuplicateNode       = errors.New("duplicate circuit node")
	ErrMissingEdgeEndpoint = errors.New("missing circuit edge endpoint")
)

type NodeID string

type PinID string

type Node struct {
	ID   NodeID
	Type NodeType
}

type PinRef struct {
	NodeID NodeID
	PinID  PinID
}

type Edge struct {
	From PinRef
	To   PinRef
}

type Graph struct {
	nodes []Node
	edges []Edge
}

func NewGraph(nodes []Node, edges []Edge) (Graph, error) {
	nodeSet := make(map[NodeID]struct{}, len(nodes))
	for _, node := range nodes {
		if _, exists := nodeSet[node.ID]; exists {
			return Graph{}, fmt.Errorf("%w: %s", ErrDuplicateNode, node.ID)
		}
		nodeSet[node.ID] = struct{}{}
	}

	for _, edge := range edges {
		if _, exists := nodeSet[edge.From.NodeID]; !exists {
			return Graph{}, fmt.Errorf("%w: %s", ErrMissingEdgeEndpoint, edge.From.NodeID)
		}
		if _, exists := nodeSet[edge.To.NodeID]; !exists {
			return Graph{}, fmt.Errorf("%w: %s", ErrMissingEdgeEndpoint, edge.To.NodeID)
		}
	}

	return Graph{
		nodes: SortNodes(nodes),
		edges: SortEdges(edges),
	}, nil
}

func (g Graph) Nodes() []Node {
	return append([]Node(nil), g.nodes...)
}

func (g Graph) Edges() []Edge {
	return append([]Edge(nil), g.edges...)
}

func SortNodes(nodes []Node) []Node {
	sorted := append([]Node(nil), nodes...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID < sorted[j].ID
	})
	return sorted
}

func SortEdges(edges []Edge) []Edge {
	sorted := append([]Edge(nil), edges...)
	sort.Slice(sorted, func(i, j int) bool {
		return edgeSortKey(sorted[i]) < edgeSortKey(sorted[j])
	})
	return sorted
}

func edgeSortKey(edge Edge) string {
	return string(edge.From.NodeID) + "\x00" +
		string(edge.From.PinID) + "\x00" +
		string(edge.To.NodeID) + "\x00" +
		string(edge.To.PinID)
}
