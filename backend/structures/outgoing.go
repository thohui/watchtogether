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
