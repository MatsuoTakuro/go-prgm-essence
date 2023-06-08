package main

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Entry struct {
	AuthorID string
	Author   string
	TitleID  string
	Title    string
	SiteURL  string
	ZipURL   string
}

func findEntires(siteURL string) ([]Entry, error) {
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return nil, err
	}
	pat := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)
	entries := []Entry{}
	doc.Find("ol li a").Each(func(n int, elem *goquery.Selection) {
		token := pat.FindStringSubmatch(elem.AttrOr("href", ""))
		// fmt.Printf("token: %#v\n", token)
		if len(token) != 3 {
			return
		}
		title := elem.Text()
		pageURL := fmt.Sprintf(
			"https://www.aozora.gr.jp/cards/%s/card%s.html",
			token[1],
			token[2],
		)
		author, zipURL := findAuthorAndZIP(pageURL)
		if zipURL != "" {
			entries = append(entries, Entry{
				AuthorID: token[1],
				Author:   author,
				TitleID:  token[2],
				Title:    title,
				SiteURL:  siteURL,
				ZipURL:   zipURL,
			})
		}
		// fmt.Printf("author: %#v, zipURL: %#v\n", author, zipURL)
	})
	return entries, nil
}

func findAuthorAndZIP(siteURL string) (author string, zipURL string) {
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return "", ""
	}

	author = doc.Find("table[summary=作家データ] tr:nth-child(2) td:nth-child(2)").Text()
	// println(author)

	doc.Find("table.download a").Each(func(n int, elem *goquery.Selection) {
		href := elem.AttrOr("href", "")
		if strings.HasSuffix(href, ".zip") {
			zipURL = href
		}
	})
	if zipURL == "" {
		return author, ""
	}

	if strings.HasPrefix(zipURL, "http://") || strings.HasPrefix(zipURL, "https://") {
		return author, zipURL
	}

	u, err := url.Parse(siteURL)
	if err != nil {
		return author, ""
	}
	u.Path = path.Join(path.Dir(u.Path), zipURL)

	return author, u.String()
}

func main() {
	URLForListOfWorksByAkutagawa := "http://www.aozora.gr.jp/index_pages/person879.html"

	es, err := findEntires(URLForListOfWorksByAkutagawa)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range es {
		fmt.Println(e.Title, e.ZipURL)
	}
}
