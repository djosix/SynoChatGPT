package internal

import "time"

type Session struct {
	Name      string
	UserID    int
	Timestamp time.Time
	Messages  Queue[Message]
	Done      chan struct{}
}

func NewSession(name string, userID int) *Session {
	return &Session{
		Name:      name,
		UserID:    userID,
		Timestamp: time.Now(),
		Messages:  NewQueue[Message](16),
		Done:      make(chan struct{}),
	}
}
