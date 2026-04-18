package server

import (
	"sort"
	"sync"
	"time"

	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/sim"
)

type RoomID string

const DefaultRoomID RoomID = "default"

type Room struct {
	id                RoomID
	mu                sync.Mutex
	simulation        *sim.Simulation
	nextClientIndex   int
	clients           map[ClientID]ClientPresence
	subscriberClients map[chan netproto.Snapshot]ClientID
}

func NewRoom(id RoomID) *Room {
	if id == "" {
		id = DefaultRoomID
	}
	return &Room{
		id:                id,
		simulation:        sim.NewSimulation(),
		clients:           make(map[ClientID]ClientPresence),
		subscriberClients: make(map[chan netproto.Snapshot]ClientID),
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
		Presence:     r.presenceSnapshotLocked(),
		Stats: sim.SnapshotStatsInput{
			ClientCount: len(r.clients),
		},
	})
}

func (r *Room) Subscribe() (<-chan netproto.Snapshot, func()) {
	snapshots := make(chan netproto.Snapshot, 4)

	r.mu.Lock()
	r.nextClientIndex++
	presence := NewClientPresence(r.nextClientIndex)
	r.clients[presence.ID] = presence
	r.subscriberClients[snapshots] = presence.ID
	r.mu.Unlock()

	unsubscribe := func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		clientID, ok := r.subscriberClients[snapshots]
		if !ok {
			return
		}
		delete(r.subscriberClients, snapshots)
		delete(r.clients, clientID)
	}

	return snapshots, unsubscribe
}

func (r *Room) PublishSnapshot(snapshot netproto.Snapshot) {
	r.mu.Lock()
	subscribers := make([]chan netproto.Snapshot, 0, len(r.subscriberClients))
	for subscriber := range r.subscriberClients {
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

func (r *Room) presenceSnapshotLocked() netproto.PresenceSnapshot {
	clients := make([]ClientPresence, 0, len(r.clients))
	for _, client := range r.clients {
		clients = append(clients, client)
	}
	sort.Slice(clients, func(i, j int) bool {
		return clients[i].sequence < clients[j].sequence
	})

	snapshots := make([]netproto.ClientPresenceSnapshot, 0, len(clients))
	for _, client := range clients {
		snapshots = append(snapshots, client.Snapshot())
	}
	return netproto.PresenceSnapshot{Clients: snapshots}
}
