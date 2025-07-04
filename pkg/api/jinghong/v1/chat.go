package v1

type Message struct {
	ID             string `json:"id"`
	MID            uint   `json:"mId"`
	Content        string `json:"content"`
	Role           string `json:"role"`
	State          string `json:"state"`
	ShowingContent string `json:"showingContent"`
	SendTime       string `json:"sendTime"`
	SessionId      string `json:"sessionId"`
}

type Session struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	ChatModel  string `json:"model"`
	CreateTime string `json:"createTime"`
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
