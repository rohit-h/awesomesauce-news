package main

import (
	"fmt"
	"sync"
)

func main() {

	posts := make(chan Post)

	go func() {
		for post := range posts {
			fmt.Println(post)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(3)

	for pageNum := 1; pageNum <= 3; pageNum++ {
		go func(pageNum int) {
			appendNewsworthyPosts(pageNum, posts)
			wg.Done()
		}(pageNum)
	}

	wg.Wait()
	close(posts)

}
