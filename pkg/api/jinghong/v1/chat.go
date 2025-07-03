package v1

type Message struct {
	ID             uint   `json:"id"`
	MID            uint   `json:"mid"`
	Content        string `json:"content"`
	Role           string `json:"role"`
	State          string `json:"state"`
	ShowingContent string `json:"showingContent"`
	SendTime       string `json:"sendTime"`
	SessionId      uint   `json:"sessionId"`
}

type Session struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	ChatModel string `json:"chatModel"`
}

type NewSessionRequest struct {
	Session
}

type NewMessageRequest struct {
	Message
}

type GetMessagesResponse struct {
	Messages []Message `json:"messages"`
}

type GetMessagesRequest struct {
	Session Session `json:"session"`
}

type GetSessionsResponse struct {
	Sessions []Session `json:"sessions"`
}
