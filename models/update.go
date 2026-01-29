package models

type Update struct {
	Message       *Message       `json:"message,omitempty"`
	DeleteMessage *DeleteMessage `json:"deleteMessage,omitempty"`
}
