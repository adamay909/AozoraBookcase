package main

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"syscall/js"

	"github.com/mtibben/percent"
)

type handleFunc func(newhash, oldhash string)

var handler map[string]handleFunc

var prefixes []string

var clicked bool

func setupjs() {

	domWindow.Call("addEventListener", "hashchange", js.FuncOf(func(this js.Value, args []js.Value) any {

		event := args[0]

		spaserver(event)
		return true
	}), true)

	domWindow.Call("addEventListener", "change", js.FuncOf(func(js.Value, []js.Value) any {
		search()
		return true
	}), true)

	domWindow.Call("addEventListener", "click", js.FuncOf(func(js.Value, []js.Value) any {
		click()
		return true
	}), true)

	return

}

func click() {

	clicked = true

}

func oldhash(event js.Value) string {

	s := strings.Split(event.Get("oldURL").String(), "#")

	if len(s) == 1 {
		return ""
	}
	return s[1]

}

func spaserver(event js.Value) {

	hash := percent.Decode(getHash())

	for _, p := range prefixes {

		if strings.HasPrefix(hash, p) {

			handler[p](strings.TrimPrefix(hash, "#"), oldhash(event))

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

func mainPages(path, oldpath string) {

	replaceBody(string(getPageData(strings.Split(path, `::`)[0])))

	domHTML.Set("style", "writing-mode: horizontal-tb")

	if elem, err := getElementById(path); err == nil {

		scrollTo(elem)

	} else {

		domWindow.Call("scrollTo", map[string]any{"top": 0, "left": 0})
	}

	if isBookPage(path) {

		epubdl, _ := getElementById("epubdl")

		azwdl, _ := getElementById("azwdl")

		epubdl.Call("addEventListener", "click", js.FuncOf(func(js.Value, []js.Value) any {
			serveFile(path, "epub")
			return true
		}), true)

		azwdl.Call("addEventListener", "click", js.FuncOf(func(js.Value, []js.Value) any {
			serveFile(path, "azw3")
			return true
		}), true)
	}

	log.Println("spaserver: done constructing page", path)

	return

}

// Thanks to https://javascript.plainenglish.io/javascript-create-file-c36f8bccb3be for how to do this

func serveFile(path string, ext string) {

	path = strings.TrimSuffix(path, filepath.Ext(path))

	/*
	   ****
	   need to do more !!!

	   ****
	*/
	patt
	= path + "." + ext

	go serveFileSvc(path)

	clicked = false

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

func readBook(path, oldpath string) {

	go readBookSvc(path)

	return
}

func readBookSvc(path string) {

	path = strings.TrimSuffix(path, ".html") + ".mono"

	replaceBody(string(getPageData(path)))

	domHTML.Set("style", "writing-mode: vertical-rl")

	return
}

func search() {

	q := domDocument.Call("getElementById", "query").Get("value").String()

	setHash("#search=" + q)

	return
}

func showSearchResult(q, old string) {

	q = percent.Decode(strings.TrimPrefix(q, "search="))

	log.Println("lookin for", q)

	replaceBody(string(globalLib.GenSearchResults(q)))

	return

}

func randomBook(s, old string) {

	log.Println("finding random book")

	bk := globalLib.RandomBook()

	hash := "#books/book_" + bk.AuthorID + "_" + bk.BookID + ".html"

	setHash(hash)

	return
}

func showMenu(s, old string) {

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
