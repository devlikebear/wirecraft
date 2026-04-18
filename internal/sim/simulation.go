package sim

import (
	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/world"
)

type Simulation struct {
	tick         TickID
	bounds       world.Dimensions
	world        *world.World
	buttonStates map[world.Position]bool
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
		bounds:       bounds,
		world:        world.New(bounds),
		buttonStates: make(map[world.Position]bool),
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

	return BuildSnapshot(SnapshotInput{
		Tick:         s.tick,
		ServerTimeMS: input.ServerTimeMS,
		World:        s.world,
		ButtonStates: s.buttonStates,
		Stats:        input.Stats,
	})
}
