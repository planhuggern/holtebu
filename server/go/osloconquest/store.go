package osloconquest

import (
	"sync"
	"errors"
)

var ErrRoomNotFound = errors.New("room not found")

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

func (store *RoomStore) UpdateRoom(roomID string, updateFunc func(Room) (Room, error)) (Room, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	room, exists := store.rooms[roomID]
	if !exists {
		return Room{}, ErrRoomNotFound
	}
	room, err := updateFunc(room)
	if err != nil {
		return Room{}, err
	}
	store.rooms[roomID] = room
	return room, nil
}
