package actuator

import "github.com/devlikebear/wirecraft/internal/physics"

type ActuatorType string

const (
	ActuatorTypePiston ActuatorType = "piston"
	ActuatorTypeMotor  ActuatorType = "motor"
)

type InputSignal string

const (
	InputSignalUnknown InputSignal = "unknown"
	InputSignalLow     InputSignal = "low"
	InputSignalHigh    InputSignal = "high"
	InputSignalPWM     InputSignal = "pwm"
)

type MotorDriverState uint8

const (
	MotorDriverDisconnected MotorDriverState = iota
	MotorDriverConnected
)

type Motor struct {
	ID             physics.EntityID
	Type           ActuatorType
	RequiresDriver bool
}

func NewMotor(id physics.EntityID) Motor {
	return Motor{
		ID:             id,
		Type:           ActuatorTypeMotor,
		RequiresDriver: true,
	}
}

func (m Motor) Enabled(signal InputSignal, driver MotorDriverState) bool {
	if m.RequiresDriver && driver != MotorDriverConnected {
		return false
	}
	return signal.IsEnergized()
}

func (s InputSignal) IsEnergized() bool {
	return s == InputSignalHigh || s == InputSignalPWM
}
