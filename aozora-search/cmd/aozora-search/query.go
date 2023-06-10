package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func showAuthors(db *sql.DB) error {
	// TODO: implement this func
	return nil
}

func showTitles(db *sql.DB, authorID string) error {
	// TODO: implement this func
	return nil
}

func showContent(db *sql.DB, authorID string, titleID string) error {
	// TODO: implement this func
	return nil
}

func queryContent(db *sql.DB, query string) error {
	// TODO: implement this func
	return nil
}
