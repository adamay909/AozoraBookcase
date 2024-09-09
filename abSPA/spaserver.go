package main

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"syscall/js"

	"github.com/mtibben/percent"
)

type handleFunc func(string)

var handler map[string]handleFunc

var prefixes []string

var clicked bool

func setupjs() {

	addEventListener(domWindow, "hashchange", spaserver)

	addEventListener(domWindow, "change", search)

}

func oldhash(event js.Value) string {

	s := strings.Split(event.Get("oldURL").String(), "#")

	if len(s) == 1 {
		return ""
	}
	return s[1]

}

func spaserver(event js.Value, params ...any) {

	hash := percent.Decode(getHash())

	for _, p := range prefixes {

		if strings.HasPrefix(hash, p) {

			handler[p](strings.TrimPrefix(hash, "#"))

			clicked = false

			return
		}
	}
	clicked = false

	return
}

func setHandler(prefix string, f handleFunc) {

	if len(handler) == 0 {
		handler = make(map[string]handleFunc)
	}

	handler[prefix] = f

	prefixes = append(prefixes, prefix)

	sortPrefixes(prefixes)

	return

}

func mainPages(path string) {

	replaceBody(string(getPageData(strings.Split(path, `::`)[0])))

	domHTML.Set("style", "writing-mode: horizontal-tb")

	if elem, err := getElementById(path); err == nil {

		scrollTo(elem)

	} else {

		domWindow.Call("scrollTo", map[string]any{"top": 0, "left": 0})
	}

	if isBookPage(path) {

		setupBookPage(path)

	}

	log.Println("spaserver: done constructing page", path)

	return

}

func setupBookPage(path string) {

	epubdl, _ := getElementById("epubdl")

	azw3dl, _ := getElementById("azw3dl")

	addEventListener(epubdl, "click", serveFile, path, "epub")

	addEventListener(azw3dl, "click", serveFile, path, "azw3")

}

func serveFile(event js.Value, params ...any) {

	path := params[0].(string)

	ext := params[1].(string)

	log.Println("creating download buttons for:", path, "type:", ext)

	pparts := strings.Split(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), "_")

	authorID, bookID := pparts[1], pparts[2]

	log.Println("authorID:", authorID, "bookID:", bookID)

	rec := globalLib.GetRecordWithID(authorID, bookID)

	path = "files/files_" + rec.RealBookID() + "/" + rec.FileName() + "_u." + ext

	go serveFileSvc(path)

	return
}

func serveFileSvc(path string) {

	log.Println("creating", filepath.Base(path))

	bk, _ := globalLib.GetBookRecord(path)

	name := bk.Title + filepath.Ext(path)

	fi := createJSFile(getPageData(path), name)

	saveFile(fi)

	log.Println("file downloaded to", name)
	return
}

func readBook(path string) {

	go readBookSvc(path)

	return
}

func readBookSvc(path string) {

	path = strings.TrimSuffix(path, ".html") + ".mono"

	replaceBody(string(getPageData(path)))

	domHTML.Set("style", "writing-mode: vertical-rl")

	return
}

func search(event js.Value, params ...any) {

	q := domDocument.Call("getElementById", "query").Get("value").String()

	setHash("#search=" + q)

	return
}

func showSearchResult(q string) {

	q = percent.Decode(strings.TrimPrefix(q, "search="))

	log.Println("lookin for", q)

	replaceBody(string(globalLib.GenSearchResults(q)))

	return

}

func randomBook(s string) {

	log.Println("finding random book")

	bk := globalLib.RandomBook()

	hash := "#books/book_" + bk.AuthorID + "_" + bk.BookID + ".html"

	setHash(hash)

	return
}

func showMenu(s string) {

	f, _ := templateFiles.Open("resources/menu.html")

	data := readFrom(f)

	replaceBody(string(data))

	return

}

// sort by length in descending order
func sortPrefixes(s []string) {

	if len(s) < 2 {
		return
	}

	for changed := true; !changed; {

		changed = false

		for i := 1; i < len(s); i++ {

			if len(s[i-1]) < len(s[i]) {

				t := s[i]

				s[i] = s[i-1]

				s[i-1] = t

				changed = true
			}
		}
	}

	return
}

func getPageData(path string) []byte {
	f, _ := globalLib.Open(path)

	defer f.Close()

	fc := f.(*cacheFile)
	defer fc.Close()

	return readFrom(fc)

}

func readFrom(f fs.File) []byte {

	info, _ := f.Stat()

	r := make([]byte, info.Size())

	f.Read(r)

	return r
}

func isBookPage(path string) bool {

	return strings.HasPrefix(path, "books/book_")

}
