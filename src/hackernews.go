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
	Id    int
	Score int
	Time  int
	Title string
	Url   string
}

func getHackerNewsTopStories() []int {
	fullUrl := apiEndpoint + "/topstories.json"
	resp, err := http.Get(fullUrl)
	if err != nil {
		assertNoError(err, "http get: "+fullUrl, 1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		assertNoError(err, "http read response: "+fullUrl, 1)
	}

	var result []int
	json.Unmarshal(body, &result)
	return result
}

func validateItemFields(item HnItem) {
	customErr := fmt.Errorf("HnItem invalid field")
	if item.Id == 0 {
		assertNoError(customErr, "Id field cannot be 0", 2)
	}
	if item.Time == 0 {
		assertNoError(customErr, "Time field cannot be 0", 2)
	}
	if item.Title == "" {
		assertNoError(customErr, "Title field cannot be empty", 2)
	}
}

func getHackerNewsItem(itemId int) HnItem {
	fullUrl := apiEndpoint + "/item/" + strconv.Itoa(int(itemId)) + ".json"
	resp, err := http.Get(fullUrl)
	if err != nil {
		assertNoError(err, "http get: "+fullUrl, 1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		assertNoError(err, "http read response: "+fullUrl, 1)
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
		log.Printf("Post meets threshold [%d<%d] [%d|%s]\n", postThreshold, story.Score, story.Id, story.Title)
		return true
	}
	timeElapsed := time.Now().Unix() - int64(story.Time)
	meanTimeBetweenUpvotes := float32(timeElapsed) / float32(story.Score)
	thresholdRate := float32(60*60) / float32(postThreshold)
	result := meanTimeBetweenUpvotes <= thresholdRate
	if result {
		log.Printf("Post going viral [%.0f<%.0f] [%d|%s]\n", meanTimeBetweenUpvotes, thresholdRate, story.Id, story.Title)
	}
	return result
}
