package server

import (
	"sync"
	"time"

	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/sim"
)

type RoomID string

const DefaultRoomID RoomID = "default"

type Room struct {
	id               RoomID
	mu               sync.Mutex
	simulation       *sim.Simulation
	connectedClients int
	subscribers      map[chan netproto.Snapshot]struct{}
}

func NewRoom(id RoomID) *Room {
	if id == "" {
		id = DefaultRoomID
	}
	return &Room{
		id:          id,
		simulation:  sim.NewSimulation(),
		subscribers: make(map[chan netproto.Snapshot]struct{}),
	}
}

func (r *Room) ID() RoomID {
	return r.id
}

func (r *Room) ApplyCommand(command netproto.Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.simulation.ApplyCommand(command)
}

func (r *Room) StepSnapshot(now time.Time) netproto.Snapshot {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.simulation.Step(sim.StepInput{
		ServerTimeMS: now.UnixMilli(),
		Stats: sim.SnapshotStatsInput{
			ClientCount: r.connectedClients,
		},
	})
}

func (r *Room) Subscribe() (<-chan netproto.Snapshot, func()) {
	snapshots := make(chan netproto.Snapshot, 4)

	r.mu.Lock()
	r.subscribers[snapshots] = struct{}{}
	r.connectedClients++
	r.mu.Unlock()

	unsubscribe := func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		if _, ok := r.subscribers[snapshots]; !ok {
			return
		}
		delete(r.subscribers, snapshots)
		r.connectedClients--
	}

	return snapshots, unsubscribe
}

func (r *Room) PublishSnapshot(snapshot netproto.Snapshot) {
	r.mu.Lock()
	subscribers := make([]chan netproto.Snapshot, 0, len(r.subscribers))
	for subscriber := range r.subscribers {
		subscribers = append(subscribers, subscriber)
	}
	r.mu.Unlock()

	for _, subscriber := range subscribers {
		select {
		case subscriber <- snapshot:
		default:
		}
	}
}
