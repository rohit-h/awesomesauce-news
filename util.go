package main

import (
	"log"
	"os"
)

func dieOnError(err error, reason string, retcode int) {
	if err != nil {
		log.Panic("FAIL while:" + reason)
		os.Exit(retcode)
	}
}

// Post "model" that will be passed around
type Post struct {
	link, title, backlink string
}
