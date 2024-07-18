package aozorafs

import (
	"html/template"
	"log"
	"net/http"

	"github.com/adamay909/AozoraConvert/jptools"
	"github.com/adamay909/AozoraConvert/runes"
)

// SearchResultsHandler is a handler function for search results. To be used with http.
func SearchResultsHandler(lib *Library) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var R struct {
			Authors, Titles, Categories []*Record
			Prefix                      string
		}
		r.ParseForm()

		if r.Form["query"][0] == "about:" {
			w.Write([]byte("AozoraBookcase.\n\nCopyright (C) 2024 Masahiro Yamada.\n\nLicensed under AGPL-3.0. \n\nSource code available at https://github.com/adamay909/AozoraBookcase"))
			return
		}

		R.Authors, R.Titles, R.Categories = search(r.Form["query"][0], lib)
		template.Must(template.New("result").Parse(searchResultsTemplate())).Execute(w, R)
		return
	}
}

func search(t string, lib *Library) (authors, titles, categories []*Record) {

	s := runes.Runes(t)

	log.Println("seach requested for :", s)

	contains := func(s []*Record, r *Record) bool {
		for _, c := range s {
			if c.BookID == r.BookID {
				return true
			}
		}
		return false
	}
	log.Println("searching through authors")

	for _, b := range lib.allAuthors() {

		if runes.Contains(runes.Runes(b.FullName()), s) {
			authors = append(authors, b)
		}

		if runes.Contains(runes.Runes(toHiragana(b.FullNameY())), s) {
			authors = append(authors, b)
		}
	}

	log.Println("searching through titles")
	//get matching titles
	for _, b := range lib.booklist {

		switch {

		case runes.Contains(runes.Runes(b.Title+b.Subtitle), s):
			titles = append(titles, b)

		case runes.Contains(runes.Runes(b.TitleY+b.SubtitleY), s):
			titles = append(titles, b)

		case runes.Contains(runes.Runes(b.TitleSort+b.SubtitleSort), s):
			titles = append(titles, b)

		default:
			continue
		}

	}
	log.Println("searching through categories")

	//get matching categories
	for _, b := range lib.booklist {

		if contains(categories, b) {
			continue
		}
		if runes.Contains(runes.Runes(b.Category), s) {
			categories = append(categories, b)
		}
		//b = lib.NextAuthor(b)
	}

	return
}

func toHiragana(s string) string {

	r := []rune(s)

	for i, c := range r {
		r[i] = jptools.ToHiragana(c)
	}

	return string(r)
}
