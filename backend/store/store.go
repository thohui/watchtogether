package store

import (
	"fmt"
	"sync"

	"github.com/thohui/watchtogether/room"
)

type Store struct {
	mutex        sync.Mutex
	rooms        map[string]*room.Room
	shutdownChan chan string
}

func New() *Store {
	store := &Store{
		mutex:        sync.Mutex{},
		rooms:        make(map[string]*room.Room),
		shutdownChan: make(chan string),
	}
	go store.startStoreTask()
	return store
}

func (store *Store) Add(room *room.Room) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	room.ShutdownChan = store.shutdownChan
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

func (store *Store) startStoreTask() {
	for {
		id := <-store.shutdownChan
		fmt.Printf("removed room %v from store\n", id)
		store.Remove(id)
	}
}
