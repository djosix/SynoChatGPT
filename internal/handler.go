package internal

import (
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	tokenizer "github.com/samber/go-gpt-3-encoder"
)

func NewHandler(chat SynoChatBot, api ChatGPTAPI) func(w http.ResponseWriter, r *http.Request) {
	sessions := map[int]*Session{}
	sessionMu := sync.Mutex{}

	encoder, err := tokenizer.NewEncoder()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			slog.Error("cannot parse form", "err", err)
			writeTextResp(w, http.StatusInternalServerError, "error")
			return
		}

		token := r.Form.Get("token")
		if token != chat.Token {
			slog.Error("invalid chat token", "token", token)
			writeTextResp(w, http.StatusUnauthorized, "error")
			return
		}

		userID, err := strconv.Atoi(r.Form.Get("user_id"))
		if err != nil {
			slog.Error("cannot parse user_id", "err", err)
			writeTextResp(w, http.StatusBadRequest, "error")
			return
		}

		name := r.Form.Get("username")
		text := r.Form.Get("text")

		slog.Info("request", "name", name, "userID", userID, "text", text)

		sessionMu.Lock()
		defer sessionMu.Unlock()

		session, ok := sessions[userID]
		if !ok {
			session = NewSession(name, userID)
			sessions[userID] = session
		}

		var numTokens int
		if tokens, err := encoder.Encode(text); err == nil {
			numTokens = len(tokens)
		} else {
			slog.Error("cannot tokenize", "err", err)
		}

		session.Messages.Push(Message{
			Role:      "user",
			Content:   text,
			Timestamp: time.Now(),
			NumTokens: numTokens,
		})

		messages := api.BuildMessages(session.Messages.Data()...)

		go func() {
			text, err := api.Call(messages)

			if err != nil {
				if err := chat.SendText(err.Error(), userID); err != nil {
					slog.Error("cannot send text", "err", err)
				}
				return
			}

			slog.Debug("response", "name", name, "userID", userID, "text", text)
			if err := chat.SendText(text, userID); err != nil {
				slog.Error("cannot send text", "err", err)
				return
			}

			sessionMu.Lock()
			defer sessionMu.Unlock()

			var numTokens int
			if tokens, err := encoder.Encode(text); err == nil {
				numTokens = len(tokens)
			}

			session.Messages.Push(Message{
				Role:      "assistant",
				Content:   text,
				Timestamp: time.Now(),
				NumTokens: numTokens,
			})
		}()
	}
}

func writeTextResp(w http.ResponseWriter, statusCode int, errMsg string) error {
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(errMsg))
	return err
}
