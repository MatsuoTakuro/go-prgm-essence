package main

import (
	"database/sql"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	_ "github.com/mattn/go-sqlite3"
)

func setupDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS authors(author_id TEXT, author TEXT, PRIMARY KEY (author_id))`,
		`CREATE TABLE IF NOT EXISTS contents(author_id TEXT, title_id TEXT, title TEXT, content TEXT, PRIMARY KEY (author_id, title_id))`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS contents_fts USING fts4(words)`,
	}
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func addEntry(db *sql.DB, entry *Entry, content string) error {
	_, err := db.Exec(
		`REPLACE INTO authors(author_id, author) values(?,?)`,
		entry.AuthorID,
		entry.Author,
	)
	if err != nil {
		return err
	}
	logger.Sugar().Infof("author updated successfully; author_id: %s, author: %s\n", entry.AuthorID, entry.Author)

	res, err := db.Exec(
		`REPLACE INTO contents(author_id, title_id, title, content) values(?,?,?,?)`,
		entry.AuthorID,
		entry.TitleID,
		entry.Title,
		content,
	)
	if err != nil {
		return err
	}
	logger.Sugar().Infof("content updated successfully; title_id: %s, title: %s\n", entry.TitleID, entry.Title)

	docID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return err
	}

	seg := t.Wakati(content)
	_, err = db.Exec(
		`INSERT INTO contents_fts(docid, words) values(?,?)`,
		docID,
		strings.Join(seg, " "),
	)
	if err != nil {
		return err
	}
	logger.Sugar().Infof("contents_fts updated successfully; docid: %d\n", docID)

	return nil
}
