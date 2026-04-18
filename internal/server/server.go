package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devlikebear/wirecraft/internal/sim"
)

type Options struct {
	TickRateHz int
}

type Server struct {
	mux         *http.ServeMux
	clock       sim.Clock
	defaultRoom *Room
	ctx         context.Context
	cancel      context.CancelFunc
}

func New() *Server {
	return NewWithOptions(Options{})
}

func NewWithOptions(options Options) *Server {
	tickRateHz := options.TickRateHz
	if tickRateHz <= 0 {
		tickRateHz = 20
	}

	ctx, cancel := context.WithCancel(context.Background())
	server := &Server{
		clock:       sim.Clock{RateHz: tickRateHz},
		defaultRoom: NewRoom(DefaultRoomID),
		ctx:         ctx,
		cancel:      cancel,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.HandleFunc("GET /ws", server.handleWebSocket)
	server.mux = mux

	go server.runSimulationLoop()

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Close() {
	s.cancel()
}

func (s *Server) runSimulationLoop() {
	ticker := time.NewTicker(s.clock.TargetDuration())
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case now := <-ticker.C:
			snapshot := s.defaultRoom.StepSnapshot(now)
			s.defaultRoom.PublishSnapshot(snapshot)
		}
	}
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, "ok")
}
