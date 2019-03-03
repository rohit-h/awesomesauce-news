package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"strings"
)

func telegramPostToChannel(apiToken, chatID string, post Post) (bool, error) {

	text := fmt.Sprintf("<a href='%s'>%s</a>\n<a href='%s'>[backlink]</a>", post.link, html.EscapeString(post.title), post.backlink)
	text = strings.Replace(text, "'", "\\\"", -1)
	payload := fmt.Sprintf("{\"chat_id\": \"%s\", \"parse_mode\": \"html\", \"disable_notification\": true, \"text\": \"%s\"}", chatID, text)

	fullURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", apiToken)
	response, err := http.Post(fullURL, "application/json", bytes.NewBufferString(payload))

	return response.StatusCode == http.StatusOK, err

}
