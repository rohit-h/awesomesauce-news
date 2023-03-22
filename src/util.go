package main

import (
	"log"
	"os"
)

func assertNoError(err error, reason string, retcode int) {
	if err != nil {
		log.Panic("FAIL while:" + reason)
		log.Panic("error:" + err.Error())
		os.Exit(retcode)
	}
}
