package osloconquest

import (
	"errors"
	"fmt"
	"sync"
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

func TestRoomStoreUpdatesRoomAndPersistsTheResult(t *testing.T) {
	store := NewRoomStore()
	firstPlayer := Player{ID: "player-1", Name: "Espen"}
	secondPlayer := Player{ID: "player-2", Name: "Ada"}

	if _, err := store.CreateWaitingRoom("room-1", firstPlayer); err != nil {
		t.Fatalf("CreateWaitingRoom returned an error: %v", err)
	}

	updatedRoom, err := store.UpdateRoom("room-1", func(room Room) (Room, error) {
		return AddPlayer(room, secondPlayer)
	})
	if err != nil {
		t.Fatalf("UpdateRoom returned an error: %v", err)
	}
	if len(updatedRoom.Players) != 2 {
		t.Errorf("updated player count = %d, want 2", len(updatedRoom.Players))
	}

	storedRoom, found := store.GetRoom("room-1")
	if !found {
		t.Fatal("GetRoom should find the updated room")
	}
	if len(storedRoom.Players) != 2 {
		t.Errorf("stored player count = %d, want 2", len(storedRoom.Players))
	}
}

func TestRoomStoreListsRooms(t *testing.T) {
	store := NewRoomStore()
	player := Player{ID: "player-1", Name: "Espen"}

	if _, err := store.CreateWaitingRoom("room-1", player); err != nil {
		t.Fatalf("CreateWaitingRoom returned an error: %v", err)
	}

	roomInfos := store.ListRooms()

	if len(roomInfos) != 1 {
		t.Fatalf("room count = %d, want 1", len(roomInfos))
	}
	if roomInfos[0].Room != "room-1" {
		t.Errorf("room ID = %q, want %q", roomInfos[0].Room, "room-1")
	}
	if roomInfos[0].PlayerCount != 1 {
		t.Errorf("player count = %d, want 1", roomInfos[0].PlayerCount)
	}
}

func TestRoomStoreGetRoomReturnsCopy(t *testing.T) {
	store := NewRoomStore()
	player := Player{ID: "player-1", Name: "Espen"}

	if _, err := store.CreateWaitingRoom("room-1", player); err != nil {
		t.Fatalf("CreateWaitingRoom returned an error: %v", err)
	}

	room, found := store.GetRoom("room-1")
	if !found {
		t.Fatal("GetRoom should find the created room")
	}
	room.Players[0].Name = "Changed outside the store"

	storedRoom, found := store.GetRoom("room-1")
	if !found {
		t.Fatal("GetRoom should still find the room")
	}
	if storedRoom.Players[0].Name != "Espen" {
		t.Errorf("stored player name = %q, want %q", storedRoom.Players[0].Name, "Espen")
	}
}

func TestRoomStoreCreatesRoomsConcurrently(t *testing.T) {
	store := NewRoomStore()
	const roomCount = 10

	errCh := make(chan error, roomCount)
	var waitGroup sync.WaitGroup
	for i := range roomCount {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			roomID := fmt.Sprintf("room-%d", i)
			player := Player{ID: PlayerID(fmt.Sprintf("player-%d", i))}
			if _, err := store.CreateWaitingRoom(roomID, player); err != nil {
				errCh <- err
			}
		}()
	}
	waitGroup.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("CreateWaitingRoom returned an error: %v", err)
	}

	if roomInfos := store.ListRooms(); len(roomInfos) != roomCount {
		t.Errorf("room count = %d, want %d", len(roomInfos), roomCount)
	}
}
