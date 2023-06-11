package main

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	db, err := setupDB("database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var close func()
	logger, close, err = setupLogger("log.json")
	if err != err {
		log.Fatal(err)
	}
	defer close()

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
