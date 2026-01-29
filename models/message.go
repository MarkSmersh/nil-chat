package models

type Message struct {
	ID        int    `json:"id"`
	ChatID    int    `json:"chatId"`
	FromID    int    `json:"fromId"`
	Text      string `json:"text"`
	Timestamp int    `json:"timestamp"`
}

type DeleteMessage struct {
	ID     int `json:"id"`
	ChatID int `json:"chatId"`
}
