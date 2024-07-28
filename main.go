package main

import (
	"log/slog"
	"net/http"
	"os"
	"synochatgpt/internal"
)

func abort(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func mustGetEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		abort("missing env", "name", name)
	}
	return value
}

func getEnvOr(name string, defaultValue string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}
	return value
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: func() slog.Level {
			var level slog.Level
			if err := level.UnmarshalText([]byte(getEnvOr("LOG_LEVEL", "info"))); err != nil {
				abort("invalid log level", "err", err)
			}
			return level
		}(),
	})))

	http.HandleFunc("/", internal.NewHandler(
		internal.SynoChatBot{
			IncomingURL: mustGetEnv("SYNOCHATGPT_BOT_INCOMING_URL"),
			Token:       mustGetEnv("SYNOCHATGPT_BOT_TOKEN"),
		},
		internal.ChatGPTAPI{
			ChatContext: getEnvOr("SYNOCHATGPT_API_CONTEXT", "Answers must be very concise. Reply in the same language. For Chinese questions, must reply with use traditional chinese characters."),
			ModelName:   getEnvOr("SYNOCHATGPT_API_MODEL", "gpt-3.5-turbo"),
			BearerToken: mustGetEnv("SYNOCHATGPT_API_TOKEN"),
		},
	))

	listenAddr := mustGetEnv("SYNOCHATGPT_LISTEN_ADDR")
	slog.Info("listening on address", "address", listenAddr)
	slog.Error("server shutdown", "err", http.ListenAndServe(listenAddr, nil))
}
