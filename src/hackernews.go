package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var apiEndpoint = "https://hacker-news.firebaseio.com/v0"

type HnItem struct {
	Id       int
	Score    int
	Time     int
	Title    string
	Url      string
	InDigest int
}

func getHackerNewsTopStories() []int {
	fullUrl := apiEndpoint + "/topstories.json"
	resp, err := http.Get(fullUrl)
	if err != nil {
		dieOnError(err, "http get: "+fullUrl, 1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		dieOnError(err, "http read response: "+fullUrl, 1)
	}

	var result []int
	json.Unmarshal(body, &result)
	return result
}

func validateItemFields(item HnItem) {
	customErr := fmt.Errorf("HnItem invalid field")
	if item.Id == 0 {
		dieOnError(customErr, "Id field cannot be 0", 2)
	}
	if item.Time == 0 {
		dieOnError(customErr, "Time field cannot be 0", 2)
	}
	if item.Title == "" {
		dieOnError(customErr, "Title field cannot be empty", 2)
	}
}

func getHackerNewsItem(itemId int) HnItem {
	fullUrl := apiEndpoint + "/item/" + strconv.Itoa(int(itemId)) + ".json"
	resp, err := http.Get(fullUrl)
	if err != nil {
		dieOnError(err, "http get: "+fullUrl, 1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		dieOnError(err, "http read response: "+fullUrl, 1)
	}

	var result HnItem
	json.Unmarshal(body, &result)
	// log.Printf("url: %s | body: %s\n", fullUrl, body)
	validateItemFields(result)
	return result
}

func isStoryWorthPosting(story HnItem) bool {
	var postThreshold = 100
	if story.Score >= postThreshold {
		log.Printf("Posting due to score [%d<%d] [%d|%s]\n", postThreshold, story.Score, story.Id, story.Title)
		return true
	}
	timeElapsed := time.Now().Unix() - int64(story.Time)
	meanTimeBetweenUpvotes := float32(timeElapsed) / float32(story.Score)
	thresholdRate := float32(60*60) / float32(postThreshold)
	result := meanTimeBetweenUpvotes <= thresholdRate
	if result {
		log.Printf("Posting due to popularity [%.0f<%.0f] [%d|%s]\n", meanTimeBetweenUpvotes, thresholdRate, story.Id, story.Title)
	}
	return result
}
