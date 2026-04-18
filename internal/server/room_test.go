package server

import (
	"reflect"
	"testing"
	"time"

	"github.com/devlikebear/wirecraft/internal/netproto"
	"github.com/devlikebear/wirecraft/internal/world"
)

func TestRoomSubscribeTracksJoinedClientsInSnapshotStats(t *testing.T) {
	room := NewRoom(DefaultRoomID)

	_, unsubscribeFirst := room.Subscribe()
	defer unsubscribeFirst()
	_, unsubscribeSecond := room.Subscribe()
	defer unsubscribeSecond()

	snapshot := room.StepSnapshot(time.UnixMilli(1000))
	if snapshot.Stats.ClientCount != 2 {
		t.Fatalf("snapshot.Stats.ClientCount = %d, want 2", snapshot.Stats.ClientCount)
	}

	unsubscribeSecond()
	snapshot = room.StepSnapshot(time.UnixMilli(1050))
	if snapshot.Stats.ClientCount != 1 {
		t.Fatalf("snapshot.Stats.ClientCount after unsubscribe = %d, want 1", snapshot.Stats.ClientCount)
	}
}

func TestRoomIncludesClientPresenceInSnapshots(t *testing.T) {
	room := NewRoom(DefaultRoomID)

	_, unsubscribeFirst := room.Subscribe()
	defer unsubscribeFirst()
	_, unsubscribeSecond := room.Subscribe()
	defer unsubscribeSecond()

	snapshot := room.StepSnapshot(time.UnixMilli(1000))
	want := []netproto.ClientPresenceSnapshot{
		{ID: "client-1", DisplayName: "Client 1"},
		{ID: "client-2", DisplayName: "Client 2"},
	}
	if !reflect.DeepEqual(snapshot.Presence.Clients, want) {
		t.Fatalf("snapshot.Presence.Clients = %+v, want %+v", snapshot.Presence.Clients, want)
	}

	unsubscribeFirst()
	snapshot = room.StepSnapshot(time.UnixMilli(1050))
	want = []netproto.ClientPresenceSnapshot{
		{ID: "client-2", DisplayName: "Client 2"},
	}
	if !reflect.DeepEqual(snapshot.Presence.Clients, want) {
		t.Fatalf("snapshot.Presence.Clients after unsubscribe = %+v, want %+v", snapshot.Presence.Clients, want)
	}
}

func TestRoomOwnsIndependentSimulationInstance(t *testing.T) {
	first := NewRoom(RoomID("first"))
	second := NewRoom(RoomID("second"))
	pos := world.Position{X: 1, Y: 2, Z: 3}

	if err := first.ApplyCommand(netproto.Command{
		Type:      netproto.CommandPlaceBlock,
		ClientID:  "client-1",
		CommandID: "cmd-1",
		Position:  pos,
		BlockType: world.BlockSolid,
	}); err != nil {
		t.Fatalf("first.ApplyCommand(place) error = %v, want nil", err)
	}

	firstSnapshot := first.StepSnapshot(time.UnixMilli(1000))
	secondSnapshot := second.StepSnapshot(time.UnixMilli(1000))

	if !snapshotContainsBlock(firstSnapshot, pos, world.BlockSolid) {
		t.Fatalf("first snapshot blocks = %+v, want placed block", firstSnapshot.Blocks)
	}
	if snapshotContainsBlock(secondSnapshot, pos, world.BlockSolid) {
		t.Fatalf("second snapshot blocks = %+v, want independent empty room", secondSnapshot.Blocks)
	}
}

func TestRoomPublishesSnapshotsToSubscribers(t *testing.T) {
	room := NewRoom(DefaultRoomID)
	snapshots, unsubscribe := room.Subscribe()
	defer unsubscribe()

	want := room.StepSnapshot(time.UnixMilli(1000))
	room.PublishSnapshot(want)

	select {
	case got := <-snapshots:
		if got.Tick != want.Tick {
			t.Fatalf("published snapshot tick = %d, want %d", got.Tick, want.Tick)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for published snapshot")
	}
}
