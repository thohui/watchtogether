package structures

type payloadBase struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func ChatMessagePayload(sender string, message string) payloadBase {
	type chatMessage struct {
		Sender  string `json:"sender"`
		Message string `json:"message"`
	}
	return payloadBase{
		Type: "chat",
		Data: chatMessage{sender, message},
	}
}

func VideoUpdatePayload(videoId string, time int) payloadBase {
	type videoUpdate struct {
		VideoId string `json:"video_id"`
		Time    int    `json:"time"`
	}
	return payloadBase{
		Type: "video_update",
		Data: videoUpdate{videoId, time},
	}
}

func InitPayload(videoId string, time int32) payloadBase {
	type init struct {
		VideoId string `json:"video_id"`
		Time    int32  `json:"time"`
	}
	return payloadBase{
		Type: "init",
		Data: init{videoId, time},
	}
}
