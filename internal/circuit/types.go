package circuit

import "github.com/devlikebear/wirecraft/internal/world"

type BlockRole string

const (
	RolePowerSource BlockRole = "power_source"
	RoleConductor   BlockRole = "conductor"
	RoleSwitch      BlockRole = "switch"
	RoleLogicGate   BlockRole = "logic_gate"
	RoleOutput      BlockRole = "output"
)

type NodeType string

const (
	NodeTypePowerSource NodeType = "power_source"
	NodeTypeWire        NodeType = "wire"
	NodeTypeButton      NodeType = "button"
	NodeTypeAndGate     NodeType = "and_gate"
	NodeTypeMCUOutput   NodeType = "mcu_output"
)

type SignalState uint8

const (
	SignalUnknown SignalState = iota
	SignalLow
	SignalHigh
)

type BlockMetadata struct {
	BlockType   world.BlockType
	DisplayName string
	Role        BlockRole
	InputPins   int
	OutputPins  int
	Conductive  bool
	Interactive bool
}

var blockTypes = []world.BlockType{
	world.BlockPower,
	world.BlockWire,
	world.BlockButton,
	world.BlockAndGate,
	world.BlockMCUOutput,
}

var metadataByBlock = map[world.BlockType]BlockMetadata{
	world.BlockPower: {
		BlockType:   world.BlockPower,
		DisplayName: "Power",
		Role:        RolePowerSource,
		OutputPins:  1,
	},
	world.BlockWire: {
		BlockType:   world.BlockWire,
		DisplayName: "Wire",
		Role:        RoleConductor,
		Conductive:  true,
	},
	world.BlockButton: {
		BlockType:   world.BlockButton,
		DisplayName: "Button",
		Role:        RoleSwitch,
		InputPins:   1,
		OutputPins:  1,
		Interactive: true,
	},
	world.BlockAndGate: {
		BlockType:   world.BlockAndGate,
		DisplayName: "AND Gate",
		Role:        RoleLogicGate,
		InputPins:   2,
		OutputPins:  1,
	},
	world.BlockMCUOutput: {
		BlockType:   world.BlockMCUOutput,
		DisplayName: "MCU Output",
		Role:        RoleOutput,
		InputPins:   1,
	},
}

func BlockTypes() []world.BlockType {
	return append([]world.BlockType(nil), blockTypes...)
}

func IsCircuitBlock(block world.BlockType) bool {
	_, ok := metadataByBlock[block]
	return ok
}

func MetadataForBlock(block world.BlockType) (BlockMetadata, bool) {
	metadata, ok := metadataByBlock[block]
	return metadata, ok
}

func NodeTypeForBlockRole(role BlockRole) (NodeType, bool) {
	switch role {
	case RolePowerSource:
		return NodeTypePowerSource, true
	case RoleConductor:
		return NodeTypeWire, true
	case RoleSwitch:
		return NodeTypeButton, true
	case RoleLogicGate:
		return NodeTypeAndGate, true
	case RoleOutput:
		return NodeTypeMCUOutput, true
	default:
		return "", false
	}
}

func (s SignalState) String() string {
	switch s {
	case SignalLow:
		return "low"
	case SignalHigh:
		return "high"
	default:
		return "unknown"
	}
}
