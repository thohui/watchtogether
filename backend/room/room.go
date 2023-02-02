package room

import (
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/teris-io/shortid"
	"github.com/thohui/watchtogether/structures"
	"github.com/thohui/watchtogether/util"
	"github.com/thohui/watchtogether/youtube"
	"go.uber.org/atomic"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type connection struct {
	Name string
	Id   string
	*websocket.Conn
}

type Room struct {
	Id               string
	connections      map[string]*connection
	mutex            sync.RWMutex
	video            youtube.Video
	host             *atomic.String
	paused           *atomic.Bool
	timePaused       *atomic.Duration
	pauseChan        chan bool
	startTime        time.Time
	once             sync.Once
	done             *atomic.Bool
	syncShutdownChan chan struct{}
}

func New(video youtube.Video) *Room {
	return &Room{
		Id:               shortid.MustGenerate(),
		connections:      make(map[string]*connection),
		mutex:            sync.RWMutex{},
		video:            video,
		host:             atomic.NewString(""),
		paused:           atomic.NewBool(true),
		timePaused:       atomic.NewDuration(time.Second),
		pauseChan:        make(chan bool),
		startTime:        time.Time{},
		once:             sync.Once{},
		done:             atomic.NewBool(false),
		syncShutdownChan: make(chan struct{}),
	}
}

func (r *Room) Handle(conn *websocket.Conn) {
	id := shortid.MustGenerate()
	connection := &connection{Name: util.GenerateRandomName(), Id: id, Conn: conn}
	var isHost bool
	if len(r.connections) == 0 {
		r.host.Store(id)
		isHost = true
	}
	r.mutex.Lock()
	r.connections[id] = connection
	r.mutex.Unlock()
	data, _ := json.Marshal(structures.InitPayload(r.video.ID, isHost, r.paused.Load()))
	connection.WriteMessage(websocket.TextMessage, data)
	r.handle(connection)
}

func (r *Room) remove(conn *connection) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.connections, conn.Id)
}

func (r *Room) broadcast(data []byte) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, conn := range r.connections {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (r *Room) handle(conn *connection) {
	defer func() {
		r.remove(conn)
		if len(r.connections) > 0 {
			if r.host.String() == conn.Id {
				for _, conn := range r.connections {
					r.host.Store(conn.Id)
					break
				}
			}
		} else {
			r.shutdown()
		}
	}()
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
		if messageType != websocket.TextMessage {
			continue
		}
		incoming := &structures.IncomingMessage{}
		err = json.Unmarshal(msg, incoming)
		if err != nil {
			continue
		}
		switch incoming.Type {
		case "chat":
			chatMessage := &structures.IncomingChatMessage{}
			err := json.Unmarshal(msg, chatMessage)
			if err != nil {
				continue
			}
			payload := structures.ChatMessagePayload(conn.Name, chatMessage.Message, r.host.Load() == conn.Id)
			data, _ := json.Marshal(payload)
			r.broadcast(data)
		case "pause":
			if r.host.Load() != conn.Id {
				continue
			}
			paused := r.paused.Load()
			if paused {
				continue
			}
			r.pauseChan <- true
			r.paused.Store(true)
			payload := structures.VideoUpdatePayload(r.elapsedSeconds(), true)
			data, _ := json.Marshal(payload)
			r.broadcast(data)
		case "resume":
			if r.host.Load() != conn.Id {
				continue
			}
			paused := r.paused.Load()
			if !paused {
				continue
			}
			r.once.Do(func() {
				r.startTime = time.Now() // the room starts paused, this is why we cannot set this variable on struct initialization.
				go r.syncVideo()
			})
			r.pauseChan <- false
			r.paused.Store(false)
			payload := structures.VideoUpdatePayload(r.elapsedSeconds(), false)
			data, _ := json.Marshal(payload)
			r.broadcast(data)
		}
	}
}

func (r *Room) elapsedSeconds() int32 {
	if (r.startTime == time.Time{}) {
		return 0
	}
	return int32((time.Since(r.startTime) - r.timePaused.Load()) / time.Second)
}

func (r *Room) syncVideo() {
	var ticker = &time.Ticker{}
	var init bool

	var pausedStart = time.Time{}
	for {
		select {
		case isPaused := <-r.pauseChan:
			if isPaused {
				pausedStart = time.Now()
				ticker.Stop()
				continue
			}
			ticker = time.NewTicker(time.Second * 1)
			// The first time the video gets resumed, the pausedStart variable is a zero-value. We prevent adding seconds to this zero-value by having this guard clause.
			if !init {
				init = true
				continue
			}
			r.timePaused.Add(time.Since(pausedStart))

		case <-ticker.C:
			currentTime := r.elapsedSeconds()
			if currentTime < int32(r.video.Duration) {
				if currentTime%5 == 0 { // every 5 seconds.
					payload := structures.VideoUpdatePayload(currentTime, r.paused.Load())
					data, _ := json.Marshal(payload)
					r.broadcast(data)
				}
				continue
			}
			ticker.Stop()
			time.Sleep(time.Second * 5)
			r.shutdown()

		case <-r.syncShutdownChan:
			break
		}
	}
}

func (r *Room) IsDone() bool {
	return r.done.Load()
}

func (r *Room) shutdown() {
	r.done.Store(true)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, conn := range r.connections {
		// Disconnect user
		conn.SetReadDeadline(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	}
	r.syncShutdownChan <- struct{}{}
}
