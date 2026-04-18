package sim

import (
	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/world"
)

type Simulation struct {
	tick   TickID
	bounds world.Dimensions
	world  *world.World
}

type StepInput struct {
	ServerTimeMS int64
	Stats        SnapshotStatsInput
}

func NewSimulation() *Simulation {
	return NewSimulationWithDimensions(world.DefaultDimensions())
}

func NewSimulationWithDimensions(bounds world.Dimensions) *Simulation {
	return &Simulation{
		bounds: bounds,
		world:  world.New(bounds),
	}
}

func (s *Simulation) ApplyCommand(command netproto.Command) error {
	if err := command.Validate(s.bounds); err != nil {
		return err
	}

	switch command.Type {
	case netproto.CommandPlaceBlock:
		return s.world.Set(command.Position, command.BlockType)
	case netproto.CommandRemoveBlock:
		return s.world.Remove(command.Position)
	default:
		return netproto.ErrUnknownCommandType
	}
}

func (s *Simulation) Step(input StepInput) netproto.Snapshot {
	s.tick = s.tick.Next()

	return BuildSnapshot(SnapshotInput{
		Tick:         s.tick,
		ServerTimeMS: input.ServerTimeMS,
		World:        s.world,
		Stats:        input.Stats,
	})
}
