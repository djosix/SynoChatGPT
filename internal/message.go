package internal

import "time"

type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"-"`
	NumTokens int       `json:"-"`
}
