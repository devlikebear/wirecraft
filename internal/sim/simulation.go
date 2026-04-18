package sim

import (
	"fmt"

	"github.com/devlikebear/wirecraft/internal/actuator"
	"github.com/devlikebear/wirecraft/internal/circuit"
	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/physics"
	"github.com/devlikebear/wirecraft/internal/world"
)

const defaultStepDeltaSeconds = 1.0 / 20.0

type Simulation struct {
	tick             TickID
	bounds           world.Dimensions
	world            *world.World
	buttonStates     map[world.Position]bool
	actuatorEntities map[physics.EntityID]physics.DynamicEntity
	lastServerTimeMS int64
}

type StepInput struct {
	ServerTimeMS int64
	DeltaSeconds float64
	Stats        SnapshotStatsInput
}

func NewSimulation() *Simulation {
	return NewSimulationWithDimensions(world.DefaultDimensions())
}

func NewSimulationWithDimensions(bounds world.Dimensions) *Simulation {
	return &Simulation{
		bounds:           bounds,
		world:            world.New(bounds),
		buttonStates:     make(map[world.Position]bool),
		actuatorEntities: make(map[physics.EntityID]physics.DynamicEntity),
	}
}

func (s *Simulation) ApplyCommand(command netproto.Command) error {
	if err := command.Validate(s.bounds); err != nil {
		return err
	}

	switch command.Type {
	case netproto.CommandPlaceBlock:
		if err := s.world.Set(command.Position, command.BlockType); err != nil {
			return err
		}
		delete(s.buttonStates, command.Position)
		return nil
	case netproto.CommandRemoveBlock:
		if err := s.world.Remove(command.Position); err != nil {
			return err
		}
		delete(s.buttonStates, command.Position)
		return nil
	case netproto.CommandSetButton:
		block, err := s.world.Get(command.Position)
		if err != nil {
			return err
		}
		if block != world.BlockButton {
			delete(s.buttonStates, command.Position)
			return nil
		}
		if command.ButtonPressed {
			s.buttonStates[command.Position] = true
		} else {
			delete(s.buttonStates, command.Position)
		}
		return nil
	default:
		return netproto.ErrUnknownCommandType
	}
}

func (s *Simulation) Step(input StepInput) netproto.Snapshot {
	s.tick = s.tick.Next()
	s.updateActuators(s.stepDeltaSeconds(input))

	return BuildSnapshot(SnapshotInput{
		Tick:         s.tick,
		ServerTimeMS: input.ServerTimeMS,
		World:        s.world,
		ButtonStates: s.buttonStates,
		Stats:        input.Stats,
	})
}

func (s *Simulation) ActuatorEntities() []physics.DynamicEntity {
	entities := make([]physics.DynamicEntity, 0, len(s.actuatorEntities))
	for _, entity := range s.actuatorEntities {
		entities = append(entities, entity)
	}
	return physics.SortEntities(entities)
}

func (s *Simulation) updateActuators(deltaSeconds float64) {
	graph, err := circuit.ExtractGraphFromWorld(s.world)
	if err != nil {
		s.actuatorEntities = make(map[physics.EntityID]physics.DynamicEntity)
		return
	}

	evaluation := circuit.EvaluateGraphWithInput(graph, circuit.EvaluationInput{
		ButtonStates: s.buttonStatesByNodeID(),
	})

	next := make(map[physics.EntityID]physics.DynamicEntity)
	for _, block := range s.world.OccupiedBlocks() {
		switch block.BlockType {
		case world.BlockPiston:
			entity := s.updatePiston(block.Position, signalForActuatorPosition(s.world, block.Position, evaluation), deltaSeconds)
			next[entity.ID] = entity
		case world.BlockMotor:
			entity := s.updateMotor(block.Position)
			next[entity.ID] = entity
		}
	}

	s.actuatorEntities = next
}

func (s *Simulation) buttonStatesByNodeID() map[circuit.NodeID]circuit.SignalState {
	states := make(map[circuit.NodeID]circuit.SignalState, len(s.buttonStates))
	for pos, pressed := range s.buttonStates {
		if pressed {
			states[circuit.NodeIDForPosition(pos)] = circuit.SignalHigh
		}
	}
	return states
}

func (s *Simulation) updatePiston(pos world.Position, signal actuator.InputSignal, deltaSeconds float64) physics.DynamicEntity {
	id := actuatorEntityID(world.BlockPiston, pos)
	piston := actuator.NewPiston(id, transformForBlockPosition(pos))
	current := s.actuatorEntities[id]
	if current.ID == "" {
		current = physics.DynamicEntity{
			ID:        id,
			Type:      physics.EntityTypePiston,
			Transform: piston.BaseTransform,
			Target:    piston.BaseTransform,
		}
	}

	target := piston.TargetTransform(signal)
	nextTransform := piston.Step(current.Transform, signal, deltaSeconds)
	return physics.DynamicEntity{
		ID:        id,
		Type:      physics.EntityTypePiston,
		Transform: nextTransform,
		Target:    target,
		Velocity:  velocityBetween(current.Transform.Position, nextTransform.Position, deltaSeconds),
	}
}

func (s *Simulation) updateMotor(pos world.Position) physics.DynamicEntity {
	id := actuatorEntityID(world.BlockMotor, pos)
	current := s.actuatorEntities[id]
	if current.ID != "" {
		return current
	}

	transform := transformForBlockPosition(pos)
	return physics.DynamicEntity{
		ID:        id,
		Type:      physics.EntityTypeMotor,
		Transform: transform,
		Target:    transform,
	}
}

func signalForActuatorPosition(w *world.World, pos world.Position, evaluation circuit.Evaluation) actuator.InputSignal {
	hasLow := false
	for _, offset := range actuatorInputOffsets {
		neighborPos := world.Position{X: pos.X + offset.X, Y: pos.Y + offset.Y, Z: pos.Z + offset.Z}
		block, err := w.Get(neighborPos)
		if err != nil || block == world.BlockAir || !circuit.IsCircuitBlock(block) {
			continue
		}

		switch evaluation.State(circuit.NodeIDForPosition(neighborPos)) {
		case circuit.SignalHigh:
			return actuator.InputSignalHigh
		case circuit.SignalLow:
			hasLow = true
		}
	}
	if hasLow {
		return actuator.InputSignalLow
	}
	return actuator.InputSignalUnknown
}

var actuatorInputOffsets = []world.Position{
	{X: 1, Y: 0, Z: 0},
	{X: -1, Y: 0, Z: 0},
	{X: 0, Y: 1, Z: 0},
	{X: 0, Y: -1, Z: 0},
	{X: 0, Y: 0, Z: 1},
	{X: 0, Y: 0, Z: -1},
}

func actuatorEntityID(blockType world.BlockType, pos world.Position) physics.EntityID {
	return physics.EntityID(fmt.Sprintf("%s:%d:%d:%d", blockType.String(), pos.X, pos.Y, pos.Z))
}

func transformForBlockPosition(pos world.Position) physics.Transform {
	return physics.DefaultTransform(physics.Vec3{
		X: float64(pos.X),
		Y: float64(pos.Y) + 0.5,
		Z: float64(pos.Z),
	})
}

func (s *Simulation) stepDeltaSeconds(input StepInput) float64 {
	if input.DeltaSeconds > 0 {
		if input.ServerTimeMS > 0 {
			s.lastServerTimeMS = input.ServerTimeMS
		}
		return input.DeltaSeconds
	}
	if input.ServerTimeMS > 0 {
		if s.lastServerTimeMS > 0 && input.ServerTimeMS > s.lastServerTimeMS {
			delta := float64(input.ServerTimeMS-s.lastServerTimeMS) / 1000
			s.lastServerTimeMS = input.ServerTimeMS
			return delta
		}
		s.lastServerTimeMS = input.ServerTimeMS
	}
	return defaultStepDeltaSeconds
}

func velocityBetween(from physics.Vec3, to physics.Vec3, deltaSeconds float64) physics.Vec3 {
	if deltaSeconds <= 0 {
		return physics.Vec3{}
	}
	return physics.Vec3{
		X: (to.X - from.X) / deltaSeconds,
		Y: (to.Y - from.Y) / deltaSeconds,
		Z: (to.Z - from.Z) / deltaSeconds,
	}
}
