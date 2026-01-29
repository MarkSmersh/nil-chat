package models

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	LastSeenAt  int    `json:"lastSeenAt"`
	CreatedAt   int    `json:"createdAt"`
}
