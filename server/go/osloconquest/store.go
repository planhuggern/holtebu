package osloconquest

import (
	"sync"
)

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
	rooms, err := createWaitingRoom(store.rooms, roomID, player)
	if err != nil {
		return Room{}, err
	}
	store.rooms = rooms
	return rooms[roomID], nil
}

func (store *RoomStore) GetRoom(roomID string) (Room, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	room, exists := store.rooms[roomID]
	return room, exists
}
