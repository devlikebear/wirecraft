package circuit

import (
	"fmt"

	"github.com/devlikebear/wirecraft/internal/world"
)

const DefaultPinID PinID = "body"

var adjacencyOffsets = []world.Position{
	{X: 1, Y: 0, Z: 0},
	{X: -1, Y: 0, Z: 0},
	{X: 0, Y: 1, Z: 0},
	{X: 0, Y: -1, Z: 0},
	{X: 0, Y: 0, Z: 1},
	{X: 0, Y: 0, Z: -1},
}

func ExtractGraphFromWorld(w *world.World) (Graph, error) {
	if w == nil {
		return NewGraph(nil, nil)
	}

	nodes := make([]Node, 0)
	nodesByPosition := make(map[world.Position]Node)
	for _, block := range w.OccupiedBlocks() {
		metadata, ok := MetadataForBlock(block.BlockType)
		if !ok {
			continue
		}

		nodeType, ok := NodeTypeForBlockRole(metadata.Role)
		if !ok {
			return Graph{}, fmt.Errorf("unsupported circuit block role %q for block %s", metadata.Role, metadata.BlockType)
		}

		node := Node{
			ID:   NodeIDForPosition(block.Position),
			Type: nodeType,
		}
		nodes = append(nodes, node)
		nodesByPosition[block.Position] = node
	}

	edges := make([]Edge, 0)
	for pos, node := range nodesByPosition {
		for _, offset := range adjacencyOffsets {
			neighbor, ok := nodesByPosition[positionWithOffset(pos, offset)]
			if !ok || node.ID >= neighbor.ID {
				continue
			}

			edges = append(edges, Edge{
				From: PinRef{NodeID: node.ID, PinID: DefaultPinID},
				To:   PinRef{NodeID: neighbor.ID, PinID: DefaultPinID},
			})
		}
	}

	return NewGraph(nodes, edges)
}

func NodeIDForPosition(pos world.Position) NodeID {
	return NodeID(fmt.Sprintf("%d:%d:%d", pos.X, pos.Y, pos.Z))
}

func positionWithOffset(pos world.Position, offset world.Position) world.Position {
	return world.Position{
		X: pos.X + offset.X,
		Y: pos.Y + offset.Y,
		Z: pos.Z + offset.Z,
	}
}
