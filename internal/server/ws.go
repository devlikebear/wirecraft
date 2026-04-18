package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/devlikebear/wirecraft/internal/netproto"
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

	snapshots, unsubscribe := s.defaultRoom.Subscribe()
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
			_ = s.defaultRoom.ApplyCommand(command)
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

func isWebSocketUpgrade(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket") &&
		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}
