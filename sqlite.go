package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func getDatabase(filePath string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", filePath)
	if err == nil {
		database.Exec("CREATE TABLE IF NOT EXISTS posts (hyperlink TEXT PRIMARY KEY, title TEXT, posted_on INT)")
	}

	return database, err
}

func markPosted(db *sql.DB, post Post) {
	statement, _ := db.Prepare("INSERT INTO posts (hyperlink, title, posted_on) VALUES (?, ?, ?)")
	statement.Exec(post.link, post.title, time.Now().Unix())
}

func hasBeenPosted(db *sql.DB, post Post) bool {
	results, _ := db.Query("SELECT * FROM posts WHERE hyperlink = ? AND title = ?", post.link, post.title)
	defer results.Close()
	for results.Next() {
		return true
	}
	return false
}

func purgeOldEntries(db *sql.DB) {
	olderThan := time.Now().Unix() - (7 * 24 * 60 * 60)
	db.Exec("DELETE FROM posts WHERE posted_on < ?", olderThan)
}
