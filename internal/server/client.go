package server

import (
	"fmt"

	"github.com/devlikebear/wirecraft/internal/netproto"
)

type ClientID string

type ClientPresence struct {
	ID          ClientID
	DisplayName string
	sequence    int
}

func NewClientPresence(index int) ClientPresence {
	return ClientPresence{
		ID:          ClientID(fmt.Sprintf("client-%d", index)),
		DisplayName: fmt.Sprintf("Client %d", index),
		sequence:    index,
	}
}

func (c ClientPresence) Snapshot() netproto.ClientPresenceSnapshot {
	return netproto.ClientPresenceSnapshot{
		ID:          string(c.ID),
		DisplayName: c.DisplayName,
	}
}
