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
	Id           string
	connections  map[string]*connection
	mutex        sync.Mutex
	video        youtube.Video
	time         *atomic.Int32
	host         *atomic.String
	paused       *atomic.Bool
	pauseChan    chan bool
	ShutdownChan chan<- string
}

func New(video youtube.Video) *Room {
	r := &Room{
		Id:           shortid.MustGenerate(),
		connections:  make(map[string]*connection),
		mutex:        sync.Mutex{},
		video:        video,
		time:         atomic.NewInt32(0),
		host:         atomic.NewString(""),
		paused:       atomic.NewBool(true),
		pauseChan:    make(chan bool),
		ShutdownChan: nil,
	}
	go r.startRoomTask()
	return r
}

func (room *Room) Handle(conn *websocket.Conn) {
	room.mutex.Lock()
	id := shortid.MustGenerate()
	name := util.GenerateRandomName()
	connection := &connection{Name: name, Id: id, Conn: conn}
	if len(room.connections) == 0 {
		room.host.Store(id)
	}
	room.connections[id] = connection
	room.mutex.Unlock()
	data, _ := json.Marshal(structures.InitPayload(room.video.ID, room.time.Inc(), room.host.String() == id, room.paused.Load()))
	connection.WriteMessage(websocket.TextMessage, data)
	room.handle(connection)
}

func (r *Room) remove(conn *connection) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.connections, conn.Id)
}

func (r *Room) broadcast(data []byte) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
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
			r.ShutdownChan <- r.Id
		}
	}()
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
		if messageType == websocket.TextMessage {
			incoming := &structures.IncomingMessage{}
			err := json.Unmarshal(msg, incoming)
			if err != nil {
				continue
			}
			switch incoming.Type {
			case "chat":
				chatMessage := &structures.IncomingChatMessage{}
				err := json.Unmarshal(msg, chatMessage)
				if err != nil {
					return
				}
				payload := structures.ChatMessagePayload(conn.Name, chatMessage.Message, r.host.Load() == conn.Id)
				data, _ := json.Marshal(payload)
				r.broadcast(data)
			case "pause":
				if r.host.Load() != conn.Id {
					return
				}
				paused := r.paused.Load()
				if paused {
					return
				}
				r.pauseChan <- true
				r.paused.Store(true)
				payload := structures.VideoUpdatePayload(r.time.Load(), true)
				data, _ := json.Marshal(payload)
				r.broadcast(data)
			case "resume":
				if r.host.Load() != conn.Id {
					return
				}
				paused := r.paused.Load()
				if !paused {
					return
				}
				r.pauseChan <- false
				r.paused.Store(false)
				payload := structures.VideoUpdatePayload(r.time.Load(), false)
				data, _ := json.Marshal(payload)
				r.broadcast(data)
			}
		}

	}
}

func (r *Room) startRoomTask() {
	var ticker = &time.Ticker{}
	for {
		select {
		case val := <-r.pauseChan:
			if val {
				ticker.Stop()
			} else {
				ticker = time.NewTicker(time.Second * 1)
			}
		case <-ticker.C:
			if r.time.Load() < int32(r.video.Duration) {
				time := r.time.Inc()
				// we are only sending a sync update every 5 seconds
				if time%5 == 0 {
					payload := structures.VideoUpdatePayload(time, r.paused.Load())
					data, _ := json.Marshal(payload)
					r.broadcast(data)
				}
			} else {
				ticker.Stop()
				time.Sleep(time.Second * 5)
				r.shutdown()
			}
		}
	}
}
func (r *Room) shutdown() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, conn := range r.connections {
		// hacky solution to close the connection outside of the handle function goroutine
		conn.SetReadDeadline(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	}
	r.ShutdownChan <- r.Id
}
