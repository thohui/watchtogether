package room

import (
	"fmt"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/teris-io/shortid"
	"github.com/thohui/watchtogether/structures"
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
	Id          string
	connections map[string]*connection
	mutex       sync.Mutex
	video       youtube.Video
	time        *atomic.Int32
	host        *atomic.String
	paused      *atomic.Bool
	pauseChan   chan bool
}

func New(video youtube.Video) *Room {
	r := &Room{
		Id:          shortid.MustGenerate(),
		connections: make(map[string]*connection),
		mutex:       sync.Mutex{},
		video:       video,
		time:        atomic.NewInt32(0),
		host:        atomic.NewString(""),
		paused:      atomic.NewBool(true),
		pauseChan:   make(chan bool),
	}
	go r.startRoomTask()
	return r
}
func (room *Room) Handle(conn *websocket.Conn) {
	room.mutex.Lock()
	id := shortid.MustGenerate()
	// TODO: generate a unique name
	connection := &connection{Name: "User", Id: id, Conn: conn}
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
	defer r.remove(conn)
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
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
				payload := structures.ChatMessagePayload(conn.Name, chatMessage.Message)
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
	var ticker *time.Ticker = &time.Ticker{}
	for {
		select {
		case val := <-r.pauseChan:
			if val {
				ticker.Stop()
			} else {
				ticker = time.NewTicker(time.Second * 1)
			}
		case <-ticker.C:
			if r.time.Load() > int32(r.video.Duration) {
				fmt.Println("video has ended")
				ticker.Stop()
			}
			time := r.time.Inc()
			if time%5 == 0 {
				fmt.Println("sending payload")
				payload := structures.VideoUpdatePayload(time, r.paused.Load())
				data, _ := json.Marshal(payload)
				r.broadcast(data)
			}
		}
	}
}
