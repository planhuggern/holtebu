package osloconquest

import (
	"errors"
	"testing"
)

func TestRoomStoreCreatesAndGetsWaitingRoom(t *testing.T) {
	store := NewRoomStore()
	player := Player{ID: "player-1", Name: "Espen"}

	createdRoom, err := store.CreateWaitingRoom("room-1", player)
	if err != nil {
		t.Fatalf("CreateWaitingRoom returned an error: %v", err)
	}
	if createdRoom.ID != "room-1" {
		t.Errorf("created room ID = %q, want %q", createdRoom.ID, "room-1")
	}

	storedRoom, found := store.GetRoom("room-1")
	if !found {
		t.Fatal("GetRoom should find the created room")
	}
	if storedRoom.ID != "room-1" {
		t.Errorf("stored room ID = %q, want %q", storedRoom.ID, "room-1")
	}
	if len(storedRoom.Players) != 1 || storedRoom.Players[0].ID != player.ID {
		t.Errorf("stored players = %v, want player %q", storedRoom.Players, player.ID)
	}

	_, err = store.CreateWaitingRoom("room-2", Player{})

	if !errors.Is(err, ErrInvalidPlayer) {
		t.Errorf("error = %v, want ErrInvalidPlayer", err)
	}
	if _, found := store.GetRoom("room-2"); found {
		t.Error("invalid player should not create a room")
	}
}
