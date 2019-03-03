package main

import (
	"log"
	"sync"
)

func main() {

	posts := make(chan Post)
	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(posts)
	}()

	telegramCheckEnv()

	getHackerNewsPosts(posts, &wg)

	db, _ := getDatabase("./awesomesauce.db")
	defer db.Close()

	for post := range posts {
		if !hasBeenPosted(db, post) {
			if status, err := telegramPostToChannel(post); status {
				markPosted(db, post)
			} else {
				log.Fatalln(err)
			}
		} else {
			log.Println("has been posted, skipping", post)
		}
	}

	purgeOldEntries(db)
}
