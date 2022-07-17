package room

import (
	"sync"

	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/teris-io/shortid"
	"github.com/thohui/watchtogether/structures"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type connection struct {
	Name string
	Id   string
	Host bool
	*websocket.Conn
}

type Room struct {
	Id          string
	connections map[string]*connection
	mutex       sync.Mutex
	VideoId     string
	time        int32
}

func New(videoId string) *Room {
	return &Room{
		Id:          shortid.MustGenerate(),
		connections: make(map[string]*connection),
		mutex:       sync.Mutex{},
		VideoId:     videoId,
		time:        0,
	}
}

func (room *Room) Handle(conn *websocket.Conn) {
	room.mutex.Lock()
	id := shortid.MustGenerate()
	// TODO: generate a unique name
	connection := &connection{Name: "User", Id: id, Host: len(room.connections) == 0, Conn: conn}
	room.connections[id] = connection
	room.mutex.Unlock()
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
