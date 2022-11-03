package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func dbGet(filePath string) (*sql.DB, error) {

	database, err := sql.Open("sqlite3", filePath)
	if err == nil {
		database.Exec("CREATE TABLE IF NOT EXISTS posts (id INT PRIMARY KEY, score INT, title TEXT, url TEXT, time INT, in_digest INT)")
	}

	return database, err
}

func dbInsertPost(db *sql.DB, item HnItem) {
	statement, _ := db.Prepare("INSERT INTO posts (id, score, title, url, time, in_digest) VALUES (?, ?, ?, ?, ?, 0)")
	result, err := statement.Exec(item.Id, item.Score, item.Title, item.Url, item.Time)
	if err != nil {
		log.Println(result, err)
		dieOnError(err, "error while inserting item into posts table", 4)
	}
}

func dbContainsPost(db *sql.DB, itemId int) bool {

	results, err := db.Query("SELECT * FROM posts WHERE id = ?", itemId)
	if err != nil {
		dieOnError(err, "error while querying posts with id", 1)
	}
	defer results.Close()

	// todo: figure out if we can return results.Count() >= 1 instead
	for results.Next() {
		// log.Println(results)
		return true
	}
	return false
}

// Entries older than 7 days get deleted from the table
func dbPurgeOldPosts(db *sql.DB) {
	olderThan := time.Now().Unix() - (7 * 24 * 60 * 60)
	db.Exec("DELETE FROM posts WHERE time < ?", olderThan)
}

func dbGetPostsNotInAnyDigest(db *sql.DB) []HnItem {

	posts := []HnItem{}

	results, err := db.Query("SELECT id, score, title, url, time FROM posts WHERE in_digest = 0")
	if err != nil {
		dieOnError(err, "error while querying posts not in digest", 1)
	}
	defer results.Close()
	for results.Next() {
		item := HnItem{}
		err = results.Scan(&item.Id, &item.Score, &item.Title, &item.Url, &item.Time)
		if err != nil {
			dieOnError(err, "error parsing row to HnItem", 1)
		}
		posts = append(posts, item)
	}

	return posts
}

func dbMarkPostInDigest(db *sql.DB, item HnItem) {
	statement, _ := db.Prepare("UPDATE posts SET in_digest = ? WHERE id = ?")
	result, err := statement.Exec(time.Now().Unix(), item.Id)
	if err != nil {
		log.Println(result, err)
		dieOnError(err, "error while updating in_digest in posts table", 4)
	}
}

func dbUpdatePost(db *sql.DB, item HnItem) {
	statement, _ := db.Prepare("UPDATE posts SET score = ?, title = ?, url = ? WHERE id = ?")
	result, err := statement.Exec(item.Score, item.Title, item.Url, item.Id)
	if err != nil {
		log.Println(result, err)
		dieOnError(err, "error while updating item in posts table", 4)
	}
}
