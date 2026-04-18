package sensor

import "github.com/devlikebear/wirecraft/internal/world"

type SensorType string

const (
	SensorTypeButton        SensorType = "button"
	SensorTypeProximityStub SensorType = "proximity_stub"
)

type InputState struct {
	Type     SensorType
	Position world.Position
	Active   bool
}

type Store struct {
	inputs map[world.Position]InputState
}

func NewStore() *Store {
	return &Store{
		inputs: make(map[world.Position]InputState),
	}
}

func ButtonInput(pos world.Position, pressed bool) InputState {
	return InputState{
		Type:     SensorTypeButton,
		Position: pos,
		Active:   pressed,
	}
}

func ProximityStubInput(pos world.Position) InputState {
	return InputState{
		Type:     SensorTypeProximityStub,
		Position: pos,
		Active:   false,
	}
}

func (s *Store) Set(input InputState) {
	if !input.Active {
		s.Clear(input.Position)
		return
	}
	s.inputs[input.Position] = input
}

func (s *Store) Clear(pos world.Position) {
	delete(s.inputs, pos)
}

func (s *Store) ButtonStates() map[world.Position]bool {
	states := make(map[world.Position]bool)
	for pos, input := range s.inputs {
		if input.Type == SensorTypeButton && input.Active {
			states[pos] = true
		}
	}
	return states
}
