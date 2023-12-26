package domain

type IncomingNewTweetNotification struct {
	FromUserID string `json:"from_user_id"`
	ShortTitle string `json:"short_title"`
}
