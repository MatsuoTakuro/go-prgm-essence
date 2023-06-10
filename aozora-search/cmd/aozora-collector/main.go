package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	URLForListOfWorksByAkutagawa := "http://www.aozora.gr.jp/index_pages/person879.html"

	es, err := findEntires(URLForListOfWorksByAkutagawa)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range es {
		content, err := extractText(e.ZipURL)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(e.SiteURL)
		fmt.Println(content)
	}
}
