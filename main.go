package main

import (
	"log/slog"
	"net/http"
	"os"
	"synochatgpt/internal"
)

func mustGetEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		slog.Error("missing env", "name", name)
		os.Exit(1)
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
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

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
