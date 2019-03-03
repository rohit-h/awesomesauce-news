package main

import (
	"log"
	"sync"
)

func main() {

	telegramCheckEnv()

	posts := make(chan Post)
	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(posts)
	}()

	getHackerNewsPosts(posts, &wg)

	db, _ := dbGet("./awesomesauce.db")
	defer db.Close()

	for post := range posts {
		if !dbContainsPost(db, post) {
			if status, err := telegramPostToChannel(post); status {
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
