package main

import (
	"log"
	"os"
	"sync"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment does not contain variable '%s'", key)
		os.Exit(1)
	}
	return value
}

func main() {

	tgToken := getEnv("TELEGRAM_BOT_TOKEN")
	tgChatID := getEnv("TELEGRAM_CHAT_ID")

	posts := make(chan Post)
	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(posts)
	}()

	getHackerNewsPosts(posts, &wg)

	db, err := dbGet("awesomesauce.db")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer db.Close()

	for post := range posts {
		if !dbContainsPost(db, post) {
			if success, err := telegramPostToChannel(tgToken, tgChatID, post); success {
				dbInsertPost(db, post)
			} else {
				log.Fatalln(err)
			}
		} else {
			log.Println("has been posted, skipping", post)
		}
	}

	dbPurgeOldPosts(db)

}
