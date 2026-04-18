package actuator

import (
	"testing"

	"github.com/devlikebear/wirecraft/internal/physics"
)

func TestPistonTargetFollowsInputSignal(t *testing.T) {
	piston := NewPiston(physics.EntityID("piston:1"), physics.DefaultTransform(physics.Vec3{X: 4, Y: 2, Z: 1}))
	piston.Axis = physics.Vec3{X: 1, Y: 0, Z: 0}
	piston.StrokeBlocks = 1

	tests := []struct {
		name       string
		signal     InputSignal
		wantTarget physics.Vec3
	}{
		{
			name:       "low retracts to base",
			signal:     InputSignalLow,
			wantTarget: physics.Vec3{X: 4, Y: 2, Z: 1},
		},
		{
			name:       "unknown retracts to base",
			signal:     InputSignalUnknown,
			wantTarget: physics.Vec3{X: 4, Y: 2, Z: 1},
		},
		{
			name:       "high extends one block",
			signal:     InputSignalHigh,
			wantTarget: physics.Vec3{X: 5, Y: 2, Z: 1},
		},
		{
			name:       "pwm placeholder extends like high",
			signal:     InputSignalPWM,
			wantTarget: physics.Vec3{X: 5, Y: 2, Z: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := piston.TargetTransform(tt.signal)
			if target.Position != tt.wantTarget {
				t.Fatalf("target.Position = %+v, want %+v", target.Position, tt.wantTarget)
			}
			if target.Rotation != physics.IdentityQuat() {
				t.Fatalf("target.Rotation = %+v, want identity", target.Rotation)
			}
			if target.Scale != physics.UnitVec3() {
				t.Fatalf("target.Scale = %+v, want unit scale", target.Scale)
			}
		})
	}
}

func TestPistonStepClampsMovementBySpeed(t *testing.T) {
	piston := NewPiston(physics.EntityID("piston:1"), physics.DefaultTransform(physics.Vec3{}))
	piston.Axis = physics.Vec3{X: 1, Y: 0, Z: 0}
	piston.StrokeBlocks = 1
	piston.SpeedBlocksPerSecond = 2

	current := physics.DefaultTransform(physics.Vec3{})

	next := piston.Step(current, InputSignalHigh, 0.25)
	if next.Position != (physics.Vec3{X: 0.5, Y: 0, Z: 0}) {
		t.Fatalf("next.Position = %+v, want half block extension", next.Position)
	}

	next = piston.Step(next, InputSignalHigh, 1)
	if next.Position != (physics.Vec3{X: 1, Y: 0, Z: 0}) {
		t.Fatalf("next.Position = %+v, want clamped full extension", next.Position)
	}

	next = piston.Step(next, InputSignalLow, 0.25)
	if next.Position != (physics.Vec3{X: 0.5, Y: 0, Z: 0}) {
		t.Fatalf("next.Position = %+v, want half block retraction", next.Position)
	}
}

func TestMotorRequiresDriverInput(t *testing.T) {
	motor := NewMotor(physics.EntityID("motor:1"))

	if motor.Enabled(InputSignalHigh, MotorDriverDisconnected) {
		t.Fatal("motor enabled without driver input, want disabled")
	}
	if !motor.Enabled(InputSignalHigh, MotorDriverConnected) {
		t.Fatal("motor disabled with high driver input, want enabled")
	}
	if motor.Enabled(InputSignalLow, MotorDriverConnected) {
		t.Fatal("motor enabled with low driver input, want disabled")
	}
}
