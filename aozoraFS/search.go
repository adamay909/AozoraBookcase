package aozorafs

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"

	"github.com/adamay909/AozoraConvert/jptools"
	"github.com/adamay909/AozoraConvert/runes"
)

//go:embed resources/searchForm.html
var fileSearchhtml string

//go:embed resources/searchResults.html
var fileSearchresultshtml string

// SearchResultsHandler is a handler function for search results. To be used with http.
func SearchResultsHandler(lib *Library) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var R struct {
			Authors, Titles, Categories []*Record
			Prefix                      string
			FoundA, FoundT, FoundC      bool
		}
		r.ParseForm()

		if r.Form["query"][0] == "about:" {
			w.Write([]byte("AozoraBookcase.\n\nCopyright (C) 2024 Masahiro Yamada.\n\nLicensed under AGPL-3.0. \n\nSource code available at https://github.com/adamay909/AozoraBookcase"))
			return
		}

		R.Authors, R.Titles, R.Categories = search(r.Form["query"][0], lib)

		R.FoundA = len(R.Authors) > 0
		R.FoundT = len(R.Titles) > 0
		R.FoundC = len(R.Categories) > 0

		template.Must(template.New("result").Parse(fileSearchresultshtml)).Execute(w, R)
		return
	}
}

func search(q string, lib *Library) (authors, titles, categories []*Record) {

	log.Println("seach requested for :", q)

	authors = lib.FindMatchingAuthors(q)

	titles = lib.FindMatchingTitles(q)

	categories = lib.FindMatchingCategories(q)

	//get matching titles
	return
}

func toHiragana(s string) string {

	r := []rune(s)

	for i, c := range r {
		r[i] = jptools.ToHiragana(c)
	}

	return string(r)
}

// FindMatchingAuthors finds the authors whose names include q.
func (lib *Library) FindMatchingAuthors(q string) (authors []*Record) {

	s := runes.Runes(q)

	log.Println("searching through authors")

	for _, b := range lib.allAuthors() {

		if runes.Contains(runes.Runes(b.FullName()), s) {
			authors = append(authors, b)
		}

		if runes.Contains(runes.Runes(toHiragana(b.FullNameY())), s) {
			authors = append(authors, b)
		}
	}

	return authors

}

// FindMatchingTitles finds the books whose title+subtitle contain q.
func (lib *Library) FindMatchingTitles(q string) (titles []*Record) {

	log.Println("searching through titles")

	s := runes.Runes(q)

	for _, b := range lib.booklist {

		if listContainsBook(titles, b) {
			continue
		}

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

	return titles
}

// FindMatchingCategories finds the books whose NDC category
// include q
func (lib *Library) FindMatchingCategories(q string) (categories []*Record) {

	log.Println("searching through categories for", q)

	for _, b := range lib.booklist {

		if listContainsBook(categories, b) {
			continue
		}

		for _, e := range b.Categories {

			if e[0] == q {

				categories = append(categories, b)

			}

			if e[1] == q {
				categories = append(categories, b)

			}

		}
	}

	return categories

}

func listContainsBook(l []*Record, q *Record) bool {

	for _, e := range l {
		if e.BookID == q.BookID {
			return true
		}
	}
	return false
}
