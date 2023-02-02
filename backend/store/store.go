package store

import (
	"fmt"
	"sync"

	"github.com/thohui/watchtogether/room"
)

type Store struct {
	mutex        sync.RWMutex
	rooms        map[string]*room.Room
	shutdownChan chan string
}

func New() *Store {
	store := &Store{
		mutex:        sync.RWMutex{},
		rooms:        make(map[string]*room.Room),
		shutdownChan: make(chan string),
	}
	go store.startStoreTask()
	return store
}

func (s *Store) Add(room *room.Room) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	room.ShutdownChan = s.shutdownChan
	s.rooms[room.Id] = room
}

func (s *Store) Remove(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.rooms, id)
}

func (s *Store) Get(id string) *room.Room {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.rooms[id]
}

func (s *Store) startStoreTask() {
	for {
		id := <-s.shutdownChan
		fmt.Printf("removed room %v from store\n", id)
		s.Remove(id)
	}
}
