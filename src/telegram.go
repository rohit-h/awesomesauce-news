package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

var urlRoot = "https://news.ycombinator.com/item"

type SendMessagePayload struct {
	ChatId              string `json:"chat_id"`
	ParseMode           string `json:"parse_mode"`
	DisableNotification bool   `json:"disable_notification"`
	DisableLinkPreview  bool   `json:"disable_web_page_preview"`
	Text                string `json:"text"`
}

func telegramPostToChannel(apiToken, chatID, text string, disableLinkPreview bool) (bool, error) {

	if apiToken == "dummy" || chatID == "dummy" {
		log.Println("[TG] Post to channel", text)
		return true, nil
	}

	payloadArgs := SendMessagePayload{
		ChatId:              chatID,
		ParseMode:           "html",
		DisableNotification: true,
		DisableLinkPreview:  disableLinkPreview,
		Text:                text,
	}
	// payload := fmt.Sprintf("{\"chat_id\": \"%s\", \"parse_mode\": \"html\", \"disable_notification\": true, \"text\": \"%s\"}", chatID, text)
	payloadBytes, err := json.Marshal(payloadArgs)
	if err != nil {
		dieOnError(err, "error while marshaling json for SendMessagePayload", 1)
	}

	fullURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", apiToken)
	response, err := http.Post(fullURL, "application/json", bytes.NewBuffer(payloadBytes))

	if err != nil || response.StatusCode != http.StatusOK {
		b, _ := httputil.DumpResponse(response, true)
		log.Println("Error occurred in telegramPostToChannel: ", string(b))
	}

	return response.StatusCode == http.StatusOK, err

}

func tgSendPost(apiToken, chatID string, item HnItem) (bool, error) {
	text := fmt.Sprintf("<a href='%s'>%s</a> | <a href='%s?id=%d'>discuss</a>",
		item.Url,
		html.EscapeString(item.Title),
		urlRoot,
		item.Id)
	return telegramPostToChannel(apiToken, chatID, text, false)
}

func tgSendDigest(apiToken, chatID string, items []HnItem) (bool, error) {
	messages := []string{}
	messages = append(messages, "Top posts #digest", "─────────────────────")
	for _, item := range items {
		text := fmt.Sprintf("• (<a href='%s?id=%d'>%d</a>) <a href='%s'>%s</a>",
			urlRoot,
			item.Id,
			item.Score,
			item.Url,
			html.EscapeString(item.Title))

		messages = append(messages, text)
	}
	if len(items) == 0 {
		messages = append(messages, "<em>No top posts found</em>")
	}
	text := strings.Join(messages, "\n")
	log.Println(text)
	return telegramPostToChannel(apiToken, chatID, text, true)
}
