package osloconquest

import (
	"errors"
	"sync"
)

var ErrRoomAlreadyExists = errors.New("room already exists")

type RoomStore struct {
	mu    sync.RWMutex
	rooms map[string]Room
}

func NewRoomStore() *RoomStore {
	return &RoomStore{
		rooms: make(map[string]Room),
	}
}

func (store *RoomStore) CreateWaitingRoom(roomID string, player Player) (Room, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if _, exists := store.rooms[roomID]; exists {
		return Room{}, ErrRoomAlreadyExists
	}
	room := Room{
		ID:      roomID,
		Players: []Player{player},
	}
	store.rooms[roomID] = room
	return room, nil
}

func (store *RoomStore) GetRoom(roomID string) (Room, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	room, exists := store.rooms[roomID]
	return room, exists
}
