package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SynoChatPayload struct {
	Text    string `json:"text,omitempty"`
	UserIDs []int  `json:"user_ids,omitempty"`
}

type SynoChatBot struct {
	IncomingURL string `json:"incoming_url"`
	Token       string `json:"token"`
}

func (chat *SynoChatBot) SendText(text string, userIDs ...int) error {
	payload := SynoChatPayload{
		Text:    text,
		UserIDs: userIDs,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	values := url.Values{}
	values.Add("payload", string(data))

	resp, err := http.PostForm(chat.IncomingURL, values)
	if err != nil {
		return fmt.Errorf("post form: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %v", resp.StatusCode)
	}

	return nil
}
