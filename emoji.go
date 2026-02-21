package valourgo

import "time"

type Emoji struct {
	ID            EmojiID   `json:"id"`
	CreatorUserID UserID    `json:"creatorUserId"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
}
