package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
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

type UnpinChatMessagePayload struct {
	ChatId string `json:"chat_id"`
}

type PinChatMessagePayload struct {
	ChatId              string `json:"chat_id"`
	MessageId           int    `json:"message_id"`
	DisableNotification bool
}

type MessageIdPayload struct {
	MessageId int `json:"message_id"`
}

type PostApiResponse struct {
	Ok               bool `json:"ok"`
	MessageIdPayload `json:"result"`
}

func telegramPostToChannel(apiToken, chatID, text string, disableLinkPreview bool) (PostApiResponse, error) {

	if apiToken == "dummy" || chatID == "dummy" {
		log.Println("[TG] Post to channel", text)
		return PostApiResponse{
			Ok:               true,
			MessageIdPayload: MessageIdPayload{-1},
		}, nil
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
	assertNoError(err, "error while marshaling json for SendMessagePayload", 1)

	fullURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", apiToken)
	response, err := http.Post(fullURL, "application/json", bytes.NewBuffer(payloadBytes))

	if err != nil || response.StatusCode != http.StatusOK {
		b, _ := httputil.DumpResponse(response, true)
		log.Println("Error occurred in telegramPostToChannel: ", string(b))
	}

	responseBytes, err := io.ReadAll(response.Body)
	assertNoError(err, "failed to read response body", 1)

	sentMessageResponse := PostApiResponse{}
	err = json.Unmarshal(responseBytes, &sentMessageResponse)
	assertNoError(err, "failed to parse json from sendMessage response", 1)

	log.Println(sentMessageResponse)

	return sentMessageResponse, err

}

func telegramPinNewMessage(apiToken, chatID string, messageId int) {
	unpinEndpoint := fmt.Sprintf("https://api.telegram.org/bot%s/unpinAllChatMessages", apiToken)
	pinEndpoint := fmt.Sprintf("https://api.telegram.org/bot%s/pinChatMessage", apiToken)
	payloadArgs := UnpinChatMessagePayload{chatID}
	payloadBytes, err := json.Marshal(payloadArgs)
	if err != nil {
		assertNoError(err, "error while marshaling json for SendMessagePayload", 1)
	}
	response, err := http.Post(unpinEndpoint, "application/json", bytes.NewBuffer(payloadBytes))
	body, _ := io.ReadAll(response.Body)
	log.Printf("%s %s", string(body), err)

	pinPayload, _ := json.Marshal(PinChatMessagePayload{
		ChatId:              chatID,
		MessageId:           messageId,
		DisableNotification: true,
	})
	response, err = http.Post(pinEndpoint, "application/json", bytes.NewBuffer(pinPayload))
	body, _ = io.ReadAll(response.Body)
	log.Printf("%s %s", string(body), err)
}

func tgSendPost(apiToken, chatID string, item HnItem) error {
	log.Println("tgSendPost item.Id:", item.Id)
	text := fmt.Sprintf("<a href='%s'>%s</a> | <a href='%s?id=%d'>discuss</a>",
		item.Url,
		html.EscapeString(item.Title),
		urlRoot,
		item.Id)
	_, err := telegramPostToChannel(apiToken, chatID, text, false)
	return err
}

func tgSendDigest(apiToken, chatID string, items []HnItem) error {
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
	response, err := telegramPostToChannel(apiToken, chatID, text, true)
	telegramPinNewMessage(apiToken, chatID, response.MessageId)
	return err
}
