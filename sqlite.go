package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func dbGet(filePath string) (*sql.DB, error) {

	database, err := sql.Open("sqlite3", filePath)
	if err == nil {
		database.Exec("CREATE TABLE IF NOT EXISTS posts (hyperlink TEXT PRIMARY KEY, title TEXT, posted_on INT)")
	}

	return database, err
}

func dbInsertPost(db *sql.DB, post Post) {
	statement, _ := db.Prepare("INSERT INTO posts (hyperlink, title, posted_on) VALUES (?, ?, ?)")
	statement.Exec(post.link, post.title, time.Now().Unix())
}

func dbContainsPost(db *sql.DB, post Post) bool {

	results, err := db.Query("SELECT * FROM posts WHERE hyperlink = ? AND title = ?", post.link, post.title)
	defer results.Close()

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	// todo: figure out if we can return results.Count() >= 1 instead
	for results.Next() {
		return true
	}
	return false
}

// Entries older than 7 days get deleted from the table
func dbPurgeOldPosts(db *sql.DB) {
	olderThan := time.Now().Unix() - (7 * 24 * 60 * 60)
	db.Exec("DELETE FROM posts WHERE posted_on < ?", olderThan)
}
