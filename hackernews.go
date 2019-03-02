package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

var urlRoot = "https://news.ycombinator.com"

// Returns html root element of HN page
func getHackerNewsPage(pageNum int) soup.Root {
	fullUrl := urlRoot + "/news?p=" + strconv.Itoa(pageNum)
	fmt.Println(fullUrl)
	resp, err := soup.Get(fullUrl)
	dieOnError(err, "url fetch", 1)

	return soup.HTMLParse(resp)
}

func getStoriesInPage(doc soup.Root) []string {
	storyIds := make([]string, 0)

	scoreRows := doc.FindAll("tr", "class", "athing")
	for _, element := range scoreRows {
		storyID := element.Attrs()["id"]
		storyIds = append(storyIds, storyID)
	}
	return storyIds
}

func isStoryWorthPosting(doc soup.Root, storyID string) bool {

	scoreSpan := doc.Find("span", "id", "score_"+storyID)
	if scoreSpan.Pointer == nil {
		fmt.Println("scoreSpan element does not exist")
		return false
	}
	hnPointsText := scoreSpan.Text() // Expected: "[0-9]+ points"
	hnPointsTextTokens := strings.Split(hnPointsText, " ")

	hnPoints, err := strconv.ParseInt(hnPointsTextTokens[0], 10, 16)
	dieOnError(err, "Converting HN points to int", 2)

	return hnPoints >= 100
}

func getStoryData(doc soup.Root, storyID string) Post {

	storyRow := doc.Find("tr", "id", storyID)
	hrefElem := storyRow.Find("a", "class", "storylink")

	hyperlink := hrefElem.Attrs()["href"]
	if strings.Index(hyperlink, "item") == 0 {
		hyperlink = urlRoot + "/" + hyperlink
	}

	return Post{hyperlink, hrefElem.Text()}
}

func appendNewsworthyPosts(pageNum int, posts chan Post) {

	page := getHackerNewsPage(pageNum)

	for _, storyID := range getStoriesInPage(page) {
		if isStoryWorthPosting(page, storyID) {
			hnPost := getStoryData(page, storyID)
			// fmt.Println(hnPost)
			posts <- hnPost
		}
	}

}
