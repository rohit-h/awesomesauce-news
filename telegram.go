package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func htmlSafe(text string) string {
	text = strings.Replace(text, "&", "&amp;", -1)
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)
	return text
}

var apiToken = os.Getenv("TELEGRAM_BOT_TOKEN")
var chatID = os.Getenv("TELEGRAM_CHAT_ID")

func telegramCheckEnv() {

	if apiToken == "" {
		log.Fatal(fmt.Errorf("envvar not defined: TELEGRAM_BOT_TOKEN"))
		os.Exit(1)
	}

	if chatID == "" {
		log.Fatal(fmt.Errorf("envvar not defined: TELEGRAM_CHAT_ID"))
		os.Exit(2)
	}
}

func telegramPostToChannel(post Post) (bool, error) {

	text := fmt.Sprintf("<a href='%s'>%s</a>\n<a href='%s'>[backlink]</a>", post.link, htmlSafe(post.title), post.backlink)
	text = strings.Replace(text, "'", "\\\"", -1)
	payload := fmt.Sprintf("{\"chat_id\": \"%s\", \"parse_mode\": \"html\", \"disable_notification\": true, \"text\": \"%s\"}", chatID, text)

	fullURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", apiToken)
	response, err := http.Post(fullURL, "application/json", bytes.NewBufferString(payload))

	return response.StatusCode == http.StatusOK, err

}