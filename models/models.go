package models

type MessageCounter struct {
	FromUserID  string `json:"from_user_id"`
	ToUserID    string `json:"to_user_id"`
	CountTotal  int    `json:"count_total"`
	CountUnread int    `json:"count_unread"`
}
