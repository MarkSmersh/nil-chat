package models

var (
	Global  ChatType = "global"
	Private ChatType = "private"
)

type ChatType string

type ChatPrivate struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type ChatGlobal struct {
	Title string `json:"title"`
}

type Chat struct {
	ID   int      `json:"id"`
	Type ChatType `json:"type"`
	ChatPrivate
	ChatGlobal
}
