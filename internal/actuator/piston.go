package actuator

import (
	"math"

	"github.com/devlikebear/wirecraft/internal/physics"
)

const (
	DefaultPistonStrokeBlocks         = 1
	DefaultPistonSpeedBlocksPerSecond = 4
)

type Piston struct {
	ID                   physics.EntityID
	Type                 ActuatorType
	BaseTransform        physics.Transform
	Axis                 physics.Vec3
	StrokeBlocks         float64
	SpeedBlocksPerSecond float64
}

func NewPiston(id physics.EntityID, base physics.Transform) Piston {
	return Piston{
		ID:                   id,
		Type:                 ActuatorTypePiston,
		BaseTransform:        base,
		Axis:                 physics.Vec3{X: 1, Y: 0, Z: 0},
		StrokeBlocks:         DefaultPistonStrokeBlocks,
		SpeedBlocksPerSecond: DefaultPistonSpeedBlocksPerSecond,
	}
}

func (p Piston) TargetTransform(signal InputSignal) physics.Transform {
	target := p.BaseTransform
	if signal.IsEnergized() {
		target.Position = addVec3(target.Position, scaleVec3(p.Axis, p.StrokeBlocks))
	}
	return target
}

func (p Piston) Step(current physics.Transform, signal InputSignal, deltaSeconds float64) physics.Transform {
	target := p.TargetTransform(signal)
	if deltaSeconds <= 0 || p.SpeedBlocksPerSecond <= 0 {
		return current
	}

	maxDistance := p.SpeedBlocksPerSecond * deltaSeconds
	return moveToward(current, target, maxDistance)
}

func moveToward(current physics.Transform, target physics.Transform, maxDistance float64) physics.Transform {
	delta := subtractVec3(target.Position, current.Position)
	distance := lengthVec3(delta)
	if distance == 0 || distance <= maxDistance {
		return target
	}

	step := maxDistance / distance
	next := target
	next.Position = addVec3(current.Position, scaleVec3(delta, step))
	return next
}

func addVec3(a physics.Vec3, b physics.Vec3) physics.Vec3 {
	return physics.Vec3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

func subtractVec3(a physics.Vec3, b physics.Vec3) physics.Vec3 {
	return physics.Vec3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func scaleVec3(v physics.Vec3, scale float64) physics.Vec3 {
	return physics.Vec3{
		X: v.X * scale,
		Y: v.Y * scale,
		Z: v.Z * scale,
	}
}

func lengthVec3(v physics.Vec3) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}
