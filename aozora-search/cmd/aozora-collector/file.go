package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
)

type Entry struct {
	AuthorID string
	Author   string
	TitleID  string
	Title    string
	SiteURL  string
	ZipURL   string
}

var pageURLFormat = "https://www.aozora.gr.jp/cards/%s/card%s.html"

func findEntires(siteURL string) ([]Entry, error) {
	resp, err := http.Get(siteURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	pat := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)
	entries := []Entry{}
	doc.Find("ol li a").Each(func(n int, elem *goquery.Selection) {
		token := pat.FindStringSubmatch(elem.AttrOr("href", ""))
		if len(token) != 3 {
			return
		}
		title := elem.Text()
		pageURL := fmt.Sprintf(
			pageURLFormat,
			token[1], // author id
			token[2], // title id
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
		logger.Sugar().Infof("found the src; author: %s, zipURL: %s\n", author, zipURL)
	})
	return entries, nil
}

func findAuthorAndZIP(siteURL string) (string, string) {
	resp, err := http.Get(siteURL)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", ""
	}

	author := doc.Find("table[summary=作家データ] tr:nth-child(2) td:nth-child(2)").First().Text()

	var zipURL string
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

func extractText(zipURL string) (string, error) {
	resp, err := http.Get(zipURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// read a zip file
	zipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// read a list of files from the zip file
	listOfFiles, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return "", err
	}

	for _, f := range listOfFiles.File {
		if path.Ext(f.Name) == ".txt" {
			textFile, err := f.Open()
			if err != nil {
				return "", err
			}
			textBytes, err := ioutil.ReadAll(textFile)
			textFile.Close()
			if err != nil {
				return "", err
			}
			encodedToUTF8, err := japanese.ShiftJIS.NewDecoder().Bytes(textBytes)
			if err != nil {
				return "", err
			}

			return string(encodedToUTF8), nil
		}
	}

	return "", errors.New("no more contents(text files) found in the zip file")
}
