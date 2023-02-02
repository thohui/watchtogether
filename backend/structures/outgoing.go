package structures

type payloadBase struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func ChatMessagePayload(sender string, message string, owner bool) payloadBase {
	type chatMessage struct {
		Sender  string `json:"sender"`
		Message string `json:"message"`
		Owner   bool   `json:"owner"`
	}
	return payloadBase{
		Type: "chat",
		Data: chatMessage{sender, message, owner},
	}
}

func InitPayload(videoId string, host bool, paused bool) payloadBase {
	type init struct {
		VideoId string `json:"video_id"`
		Host    bool   `json:"host"`
		Paused  bool   `json:"paused"`
	}
	return payloadBase{
		Type: "init",
		Data: init{videoId, host, paused},
	}
}

func VideoUpdatePayload(time int32, paused bool) payloadBase {
	type videoUpdate struct {
		Time   int32 `json:"time"`
		Paused bool  `json:"paused"`
	}
	return payloadBase{
		Type: "video_update",
		Data: videoUpdate{time, paused},
	}
}
