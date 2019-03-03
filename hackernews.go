package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/anaskhan96/soup"
)

var urlRoot = "https://news.ycombinator.com"

// Returns html root element of HN page
func getHackerNewsPage(pageNum int) soup.Root {
	fullURL := urlRoot + "/news?p=" + strconv.Itoa(pageNum)
	resp, err := soup.Get(fullURL)
	dieOnError(err, "url fetch : "+fullURL, 1)

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
		log.Println("scoreSpan element does not exist")
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
	backlink := fmt.Sprintf("%s/item?id=%s", urlRoot, storyID)

	return Post{hyperlink, hrefElem.Text(), backlink}
}

func getHackerNewsPosts(posts chan Post, wg *sync.WaitGroup) {

	var pagesToScrape = 3
	wg.Add(pagesToScrape)

	for pageNo := 1; pageNo <= pagesToScrape; pageNo++ {

		go func(pageNo int) {
			page := getHackerNewsPage(pageNo)
			stories := getStoriesInPage(page)

			for _, storyID := range stories {
				if isStoryWorthPosting(page, storyID) {
					posts <- getStoryData(page, storyID)
				}
			}
			wg.Done()
		}(pageNo)
	}

}
