package main

import (
	"log"
	"path/filepath"
	"strings"
	"syscall/js"

	"github.com/mtibben/percent"
)

type handleFunc func(string)

var handler map[string]handleFunc

var prefixes []string

func setupjs() {

	domWindow.Call("addEventListener", "hashchange", js.FuncOf(func(js.Value, []js.Value) any {
		spaserver()
		return true
	}), true)

	domWindow.Call("addEventListener", "change", js.FuncOf(func(js.Value, []js.Value) any {
		search()
		return true
	}), true)

	domWindow.Call("addEventListener", "popstate", js.FuncOf(func(js.Value, []js.Value) any {
		backbutton()
		return true
	}), true)

	return

}

func backbutton() {

	log.Println("backbutton pressed")

}

func spaserver() {

	hash := getHash()

	for _, p := range prefixes {

		if strings.HasPrefix(hash, p) {

			handler[p](strings.TrimPrefix(hash, "#"))

			return
		}
	}
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

	f, _ := globalLib.Open(path)

	defer f.Close()

	fc := f.(*cacheFile)
	defer fc.Close()

	r := make([]byte, fc.Size())

	fc.Read(r)

	log.Println("spaserver: done constructing page", path)

	replaceBody(string(r))

	domHTML.Set("style", "writing-mode: horizontal-tb")

	domWindow.Call("scrollTo", map[string]any{"top": 0, "left": 0})

	return

}

// Thanks to https://javascript.plainenglish.io/javascript-create-file-c36f8bccb3be for how to do this

func serveFile(path string) {

	go serveFileSvc(path)

	return
}

func serveFileSvc(path string) {

	log.Println("creating", filepath.Base(path))

	f, err := globalLib.Open(path)
	defer f.Close()

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("done")

	info, err := f.Stat()

	data := make([]byte, info.Size())

	f.Read(data)

	bk, _ := globalLib.GetBookRecord(path)

	fi := createJSFile(data, bk.Title+filepath.Ext(path))

	saveFile(fi)

	return
}

func readBook(path string) {

	go readBookSvc(path)

	return
}

func readBookSvc(path string) {

	path = strings.TrimSuffix(path, ".html") + ".mono"

	f, _ := globalLib.Open(path)

	defer f.Close()

	fc := f.(*cacheFile)
	defer fc.Close()

	r := make([]byte, fc.Size())

	fc.Read(r)

	html := string(r)

	log.Println("spaserver: done constructing page", path)

	replaceBody(html)

	domHTML.Set("style", "writing-mode: vertical-rl")

	return
}

func getBookFile(path string) []byte {

	f, err := globalLib.Open(path)

	if err != nil {
		log.Println("something wrong")
		log.Println(err)
		return []byte("")
	}

	info, err := f.Stat()

	data := make([]byte, info.Size())

	f.Read(data)

	return data

}

func search() {

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
