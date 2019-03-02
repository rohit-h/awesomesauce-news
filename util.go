package main

import (
	"fmt"
	"os"
)

func dieOnError(err error, reason string, retcode int) {
	if err != nil {
		fmt.Printf("FAIL while: %s\n", reason)
		fmt.Printf("     cause: %s\n", err)
		os.Exit(retcode)
	}
}

type Post struct {
	link, title string
}
