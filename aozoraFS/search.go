package aozorafs

import (
	"bytes"
	"log"
	"strconv"
	"strings"

	azrconv "github.com/adamay909/AozoraConvert/v2"
)

func (lib *Library) GenSearchResults(q string) []byte {

	var R struct {
		Authors, Titles        []*Record
		Categories             []string
		Prefix                 string
		FoundA, FoundT, FoundC bool
	}

	w := new(bytes.Buffer)

	if q == "about:" {
		w.Write([]byte("AozoraBookcase.\n\nCopyright (C) 2024 Masahiro Yamada.\n\nLicensed under AGPL-3.0. \n\nSource code available at https://github.com/adamay909/AozoraBookcase"))
		return w.Bytes()
	}

	R.Authors, R.Titles, R.Categories = lib.search(q)

	R.FoundA = len(R.Authors) > 0
	R.FoundT = len(R.Titles) > 0
	R.FoundC = len(R.Categories) > 0

	err := lib.searchresultT.Execute(w, R)
	if err != nil {
		log.Println(err)
	}

	return w.Bytes()

}

func (lib *Library) search(q string) (authors, titles []*Record, categories []string) {

	q = strings.TrimSpace(q)

	log.Println("seach requested for :", q)

	authors = lib.FindMatchingAuthors(q)

	titles = lib.FindMatchingTitles(q)

	categories = lib.FindMatchingCategories(q)

	//get matching titles
	return
}

func toHiragana(s string) string {

	var r []rune

	for _, c := range s {
		r = append(r, azrconv.ToHiragana(c))
	}

	return string(r)
}

// FindMatchingAuthors finds the authors whose names include q.
func (lib *Library) FindMatchingAuthors(q string) (authors []*Record) {

	for _, b := range lib.allAuthors() {

		if strings.Contains(b.FullName(), q) {
			authors = append(authors, b)
		}

		if strings.Contains(toHiragana(b.FullNameY()), q) {
			authors = append(authors, b)
		}
	}

	return authors

}

// FindMatchingTitles finds the books whose title+subtitle contain q.
func (lib *Library) FindMatchingTitles(q string) (titles []*Record) {

	var listed = make(map[string]bool)

	for _, b := range lib.booklist {

		if _, ok := listed[b.BookID]; ok {
			continue
		}

		switch {

		case strings.Contains(b.Title+b.Subtitle, q):
			titles = append(titles, b)
			listed[b.BookID] = true

		case strings.Contains(b.TitleY+b.SubtitleY, q):
			titles = append(titles, b)
			listed[b.BookID] = true

		case strings.Contains(b.TitleSort+b.SubtitleSort, q):
			titles = append(titles, b)
			listed[b.BookID] = true

		case strings.Contains(toHiragana(b.Title+b.Subtitle), q):
			titles = append(titles, b)
			listed[b.BookID] = true

		case strings.Contains(toHiragana(b.TitleY+b.SubtitleY), q):
			titles = append(titles, b)
			listed[b.BookID] = true

		case strings.Contains(toHiragana(b.TitleSort+b.SubtitleSort), q):
			titles = append(titles, b)
			listed[b.BookID] = true
		default:
			continue
		}

	}

	return titles
}

// FindMatchingCategories finds the books whose NDC category
// include q
func (lib *Library) FindMatchingCategories(q string) (categories []string) {

	for n := 1; n < 1000; n++ {

		code := ndcNumeric(n)

		if _, ok := lib.Categories[code]; !ok {
			continue
		}

		found := false

		if strings.Contains(lib.Categories[code[:1]], q) {
			found = true
		}
		if len(code) > 1 {

			if strings.Contains(lib.Categories[code[:2]], q) {
				found = true
			}
		}
		if len(code) > 2 {

			if strings.Contains(lib.Categories[code], q) {
				found = true
			}
		}

		if found {
			if lib.categoryHasBooks(code) {
				categories = append(categories, code)
			}
		}
	}

	return
}

func (lib *Library) categoryHasBooks(code string) bool {

	for _, b := range lib.booklist {

		for _, c := range b.Categories {
			if c[0] == code {
				return true
			}
			if c[1] == code {
				return true
			}
			if c[2] == code {
				return true
			}
		}
	}
	return false
}

func (lib *Library) FindBooksWithMatchingCategories(q string) (categories []*Record) {

	var matched map[*Record]bool

	for _, b := range lib.booklist {

		if _, ok := matched[b]; ok {
			continue
		}

		for _, e := range b.Categories {
			found := false
			if e[0] == q {
				found = true
			}
			if e[1] == q {
				found = true
			}
			if e[2] == q {
				found = true
			}
			if found == true {
				categories = append(categories, b)
			}
		}
	}

	return categories

}

func ndcNumeric(n int) string {

	var val string

	switch {
	case n < 10:
		val = strconv.Itoa(n)

	case n > 10 && n < 20:
		val = "0" + strconv.Itoa(n-10)

	default:
		val = strconv.Itoa(n - 10)

	}

	n++

	return val
}
