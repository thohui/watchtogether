package room

import (
	"sync"
	"sync/atomic"

	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/teris-io/shortid"
	"github.com/thohui/watchtogether/structures"
	"github.com/thohui/watchtogether/youtube"
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
	time        int32
	host        string
}

func New(video youtube.Video) *Room {
	return &Room{
		Id:          shortid.MustGenerate(),
		connections: make(map[string]*connection),
		mutex:       sync.Mutex{},
		video:       video,
		time:        0,
		host:        "",
	}
}

func (room *Room) Handle(conn *websocket.Conn) {
	room.mutex.Lock()
	id := shortid.MustGenerate()
	// TODO: generate a unique name
	connection := &connection{Name: "User", Id: id, Conn: conn}
	if len(room.connections) == 0 {
		room.host = id
	}
	room.connections[id] = connection
	room.mutex.Unlock()
	time := atomic.LoadInt32(&room.time)
	data, _ := json.Marshal(structures.InitPayload(room.video.ID, time))
	connection.WriteMessage(websocket.TextMessage, data)
	room.handle(connection)
}
func (r *Room) remove(conn *connection) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.connections, conn.Id)
}

func (r *Room) broadcast(sender, message string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	payload := structures.ChatMessagePayload(sender, message)
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
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
				r.broadcast(conn.Name, chatMessage.Message)
			}
		}
	}
}
