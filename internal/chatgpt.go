package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type ChatGPTAPI struct {
	ChatContext string
	ModelName   string
	BearerToken string
}

func (c ChatGPTAPI) Call(messages []Message) (string, error) {
	type ReqBody struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}

	reqBody := ReqBody{
		Model:    c.ModelName,
		Messages: messages,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.BearerToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http client do: %w", err)
	}

	type RespBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	var respBody RespBody
	if err := json.Unmarshal(body, &respBody); err != nil {
		return "", fmt.Errorf("unmarshal response body: %w", err)
	}

	var content string
	if n := len(respBody.Choices); n == 0 {
		content = string(body)
	} else {
		n = rand.Intn(n)
		content = respBody.Choices[n].Message.Content
	}

	return content, nil
}

func (c *ChatGPTAPI) BuildMessages(messages ...Message) []Message {
	maxNumTokens := 1000
	minTimestamp := time.Now().Add(-30 * time.Minute)

	numTokens := 0

	n := len(messages)
	for n > 0 {
		message := messages[n-1]
		if !message.Timestamp.After(minTimestamp) {
			break
		}
		numTokens += message.NumTokens
		if numTokens > maxNumTokens {
			break
		}
		n--
	}

	if n == len(messages) {
		n = len(messages) - 1 // ensure at least one message
	}

	outputMessages := []Message{}
	if c.ChatContext != "" {
		outputMessages = append(outputMessages, Message{
			Role:    "system",
			Content: c.ChatContext,
		})
	}
	return append(outputMessages, messages[n:]...)
}
