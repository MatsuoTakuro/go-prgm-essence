package main

import (
	"log"
)

var logfilePath = "./cmd/aozora-collector/log.json"

func main() {
	db, err := setupDB("database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var closeLogger func() error
	logger, closeLogger, err = setupLogger(logfilePath)
	if err != err {
		log.Fatal(err)
	}
	defer func() {
		if closeLogger == nil {
			// in case of testing, closing logger will be done at teardown fuction
			return
		}
		if err := closeLogger(); err != nil {
			log.Println("Error during closing logger:", err)
		}
	}()

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
