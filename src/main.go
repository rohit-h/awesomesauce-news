package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"sort"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment does not contain variable '%s'", key)
		os.Exit(1)
	}
	return value
}

func scrapeAndPost(db *sql.DB, tgToken string, tgChatId string) {
	for _, itemId := range getHackerNewsTopStories() {
		if dbContainsPost(db, itemId) {
			log.Println("has been posted, skipping", itemId)
		} else {
			itemData := getHackerNewsItem(itemId)
			if isStoryWorthPosting(itemData) {
				err := tgSendPost(tgToken, tgChatId, itemData)
				assertNoError(err, "failed to sendMessage", 1)
				dbInsertPost(db, itemData)
			}
		}
	}

}

func sendDigestPost(db *sql.DB, tgToken string, tgChatId string) {
	// fetch posts that havent been included in a digest

	log.Println("Fetching unsent posts and updating scores")
	posts := dbGetPostsNotInAnyDigest(db)
	for _, post := range posts {
		itemData := getHackerNewsItem(post.Id)
		dbUpdatePost(db, itemData)
	}

	// scan and sort
	log.Println("Picking top 10 posts from DB")
	posts = dbGetPostsNotInAnyDigest(db)
	sort.Slice(posts[:], func(a, b int) bool {
		return posts[a].Score > posts[b].Score
	})

	var topPosts = []HnItem{}
	if len(posts) < 10 {
		topPosts = posts[:]
	} else {
		topPosts = posts[:10]
	}

	err := tgSendDigest(tgToken, tgChatId, topPosts)
	if err != nil {
		assertNoError(err, "failed to post digest, http error", 2)
	}

	for _, post := range topPosts {
		dbMarkPostInDigest(db, post)
	}
}

func main() {

	digestFlag := flag.Bool("digest", false, "send daily digest post")
	flag.Parse()

	tgToken := getEnv("TELEGRAM_BOT_TOKEN")
	tgChatId := getEnv("TELEGRAM_CHAT_ID")

	db, err := dbGet("awesomesauce.db")
	if err != nil {
		assertNoError(err, "failed to open db file", 1)
	}
	defer db.Close()

	if *digestFlag {
		log.Println("Compiling digest message")
		sendDigestPost(db, tgToken, tgChatId)
		return
	}
	scrapeAndPost(db, tgToken, tgChatId)
	dbPurgeOldPosts(db)
}
