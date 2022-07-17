package structures

type IncomingMessage struct {
	Type string `json:"type"`
}

type IncomingChatMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
