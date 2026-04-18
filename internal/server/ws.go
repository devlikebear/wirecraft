package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/sim"
)

const websocketWriteTimeout = time.Second

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	if !isWebSocketUpgrade(r) {
		http.Error(w, "websocket upgrade required", http.StatusUpgradeRequired)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	snapshots, unsubscribe := s.subscribe()
	defer unsubscribe()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	commands := make(chan netproto.Command)
	readErrors := make(chan error, 1)
	go readCommands(ctx, conn, commands, readErrors)

	for {
		select {
		case <-ctx.Done():
			return
		case <-readErrors:
			return
		case command := <-commands:
			_ = s.applyCommand(command)
		case snapshot := <-snapshots:
			writeCtx, writeCancel := context.WithTimeout(ctx, websocketWriteTimeout)
			err := wsjson.Write(writeCtx, conn, snapshot)
			writeCancel()
			if err != nil {
				return
			}
		}
	}
}

func readCommands(ctx context.Context, conn *websocket.Conn, commands chan<- netproto.Command, errors chan<- error) {
	for {
		var command netproto.Command
		if err := wsjson.Read(ctx, conn, &command); err != nil {
			select {
			case errors <- err:
			default:
			}
			return
		}

		select {
		case commands <- command:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Server) applyCommand(command netproto.Command) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.simulation.ApplyCommand(command)
}

func (s *Server) stepSnapshot(now time.Time) netproto.Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.simulation.Step(sim.StepInput{
		ServerTimeMS: now.UnixMilli(),
		Stats: sim.SnapshotStatsInput{
			ClientCount: s.connectedClients,
		},
	})
}

func (s *Server) subscribe() (<-chan netproto.Snapshot, func()) {
	snapshots := make(chan netproto.Snapshot, 4)

	s.mu.Lock()
	s.subscribers[snapshots] = struct{}{}
	s.connectedClients++
	s.mu.Unlock()

	unsubscribe := func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		if _, ok := s.subscribers[snapshots]; !ok {
			return
		}
		delete(s.subscribers, snapshots)
		s.connectedClients--
	}

	return snapshots, unsubscribe
}

func (s *Server) publishSnapshot(snapshot netproto.Snapshot) {
	s.mu.Lock()
	subscribers := make([]chan netproto.Snapshot, 0, len(s.subscribers))
	for subscriber := range s.subscribers {
		subscribers = append(subscribers, subscriber)
	}
	s.mu.Unlock()

	for _, subscriber := range subscribers {
		select {
		case subscriber <- snapshot:
		default:
		}
	}
}

func isWebSocketUpgrade(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket") &&
		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}
