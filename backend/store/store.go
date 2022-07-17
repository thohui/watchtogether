package store

import (
	"sync"

	"github.com/thohui/watchtogether/room"
)

type Store struct {
	mutex sync.Mutex
	rooms map[string]*room.Room
}

func New() *Store {
	return &Store{
		mutex: sync.Mutex{},
		rooms: make(map[string]*room.Room),
	}
}

func (store *Store) Add(room *room.Room) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.rooms[room.Id] = room
}

func (store *Store) Remove(id string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	delete(store.rooms, id)
}

func (store *Store) Get(id string) *room.Room {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return store.rooms[id]
}
