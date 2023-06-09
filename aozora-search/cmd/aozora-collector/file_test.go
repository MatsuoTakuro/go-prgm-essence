package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
)

func Test_findEntires(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Opening at", r.URL.String())
		if r.URL.String() == "/" {
			// returns a detail of author; テスト 太郎 (author_id: `999999``)
			w.Write([]byte(`
			<table summary="作家データ">
			<tr><td class="header">作家名：</td><td><font size="+2">テスト 太郎</font></td></tr>
			<tr><td class="header">作家名読み：</td><td>テスト 太郎</td></tr>
			<tr><td class="header">ローマ字表記：</td><td>Test, Taro</td></tr>
			</table>
			<ol>
			<li><a href="../cards/999999/card001.html">テスト書籍001</a></li> 
			<li><a href="../cards/999999/card002.html">テスト書籍002</a></li> 
			<li><a href="../cards/999999/card003.html">テスト書籍003</a></li> 
			</ol>
			`))
		} else {
			// returns a detail of work; テスト書籍001 (title_id: `001`), テスト書籍002 (title_id: `002`) and テスト書籍003 (title_id: `003`).
			pattern := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)
			tokens := pattern.FindStringSubmatch(r.URL.String())
			w.Write([]byte(fmt.Sprintf(`
			<table summary="作家データ">
			<tr><td class="header">作家名：</td><td><font size="+2">テスト 太郎</font></td></tr>
			<tr><td class="header">作家名読み：</td><td>テスト 太郎</td></tr>
			<tr><td class="header">ローマ字表記：</td><td>Test, Taro</td></tr>
			</table>
			<table border="1" summary="ダウンロードデータ" class="download">
			<tr>
					<td><a href="./files/%[1]s_%[2]s.zip">%[1]s_%[2]s.zip</a></td>
			</tr>
			</table>
			`,
				tokens[1], // author id
				tokens[2], // title id
			)))
		}
	}))
	defer ts.Close()

	tmp := pageURLFormat
	pageURLFormat = ts.URL + "/cards/%s/card%s.html"
	defer func() {
		pageURLFormat = tmp
	}()

	got, err := findEntires(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}

	want := []Entry{
		{
			AuthorID: "999999",
			Author:   "テスト 太郎",
			TitleID:  "001",
			Title:    "テスト書籍001",
			SiteURL:  ts.URL,
			ZipURL:   ts.URL + "/cards/999999/files/999999_001.zip",
		},
		{
			AuthorID: "999999",
			Author:   "テスト 太郎",
			TitleID:  "002",
			Title:    "テスト書籍002",
			SiteURL:  ts.URL,
			ZipURL:   ts.URL + "/cards/999999/files/999999_002.zip",
		},
		{
			AuthorID: "999999",
			Author:   "テスト 太郎",
			TitleID:  "003",
			Title:    "テスト書籍003",
			SiteURL:  ts.URL,
			ZipURL:   ts.URL + "/cards/999999/files/999999_003.zip",
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %+v, but got %+v", want, got)
	}
}

func Test_extractText(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir(".")))
	defer ts.Close()

	got, err := extractText(ts.URL + "/testdata/example.zip")
	if err != nil {
		t.Fatal(err)
		return
	}
	want := "test data\n"
	if want != got {
		t.Errorf("want %+v, but got %+v", want, got)
	}
}
