package main

import (
	"log"
)

func main() {
	db, err := setupDB("database.sqlite")
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
		err = addEntry(db, &e, content)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
