package store

import (
	"errors"
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
	return store
}

func (s *Store) Add(room *room.Room) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.rooms[room.Id] = room
}

func (s *Store) Remove(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.rooms, id)
}

func (s *Store) Get(id string) (*room.Room, error) {
	s.mutex.RLock()
	room, ok := s.rooms[id]
	if !ok {
		return nil, errors.New("room does not exist")
	}
	if room.IsDone() {
		s.mutex.RUnlock()
		s.mutex.Lock()
		delete(s.rooms, id)
		s.mutex.Unlock()
		return nil, errors.New("room is finished!")
	}
	s.mutex.RUnlock()
	return room, nil
}

// TODO: implement this
func (s *Store) cleanRooms() {
}
