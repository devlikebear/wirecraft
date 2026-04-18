package world

type BlockType uint8

const (
	BlockAir BlockType = iota
	BlockSolid
	BlockDebugMover
	BlockPower
	BlockWire
	BlockButton
	BlockAndGate
	BlockMCUOutput
	BlockPiston
	BlockMotor
	BlockMotorDriver
	BlockTransistorSwitch
)

func (b BlockType) Valid() bool {
	switch b {
	case BlockAir, BlockSolid, BlockDebugMover, BlockPower, BlockWire, BlockButton, BlockAndGate, BlockMCUOutput, BlockPiston, BlockMotor, BlockMotorDriver, BlockTransistorSwitch:
		return true
	default:
		return false
	}
}

func (b BlockType) String() string {
	switch b {
	case BlockAir:
		return "air"
	case BlockSolid:
		return "solid"
	case BlockDebugMover:
		return "debug_mover"
	case BlockPower:
		return "power"
	case BlockWire:
		return "wire"
	case BlockButton:
		return "button"
	case BlockAndGate:
		return "and_gate"
	case BlockMCUOutput:
		return "mcu_output"
	case BlockPiston:
		return "piston"
	case BlockMotor:
		return "motor"
	case BlockMotorDriver:
		return "motor_driver"
	case BlockTransistorSwitch:
		return "transistor_switch"
	default:
		return "invalid"
	}
}
