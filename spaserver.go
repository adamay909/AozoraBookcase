package main

import (
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"syscall/js"

	str "github.com/adamay909/AozoraBookcase/aozoraFS/stringops"
)

type handleFunc func(string)

var handler map[string]handleFunc

var prefixes []string

func spaserver(event js.Value, params ...any) {

	defer func() {

		if r := recover(); r != nil {

			log.Println("something wrong!!!")
			uncoverElement(domBody)
			return
		}
	}()

	hash := str.PercentDecode(getHash())

	for _, p := range prefixes {

		if strings.HasPrefix(hash, p) {

			handler[p](strings.TrimPrefix(hash, "#"))

			return
		}
	}

	return
}

func setupJS() {

	setupHashHandlers()

	setupGlobalEventListeners()

}

func setupHashHandlers() {

	setHashHandler("#index.html", mainPages)

	setHashHandler("#authors", mainPages)

	setHashHandler("#books", mainPages)

	setHashHandler("#categories", mainPages)

	setHashHandler("#recent", recentsPage)

	setHashHandler("#read", readBook)

	setHashHandler("#search=", showSearchResult)

	setHashHandler("#about", showAbout)

}

func setupGlobalEventListeners() {

	addEventListener(domWindow, "hashchange", spaserver)

}

func setHashHandler(prefix string, f handleFunc) {

	if len(handler) == 0 {
		handler = make(map[string]handleFunc)
	}

	handler[prefix] = f

	prefixes = append(prefixes, prefix)

	sortPrefixes(prefixes)

	return

}

func mainPages(path string) {

	mkpage(path, string(getPageData(strings.Split(path, `::`)[0])))

	domHTML.Set("style", "writing-mode: horizontal-tb")

	if elem, err := getElementByID(path); err == nil {

		scrollTo(elem)

	} else {

		domWindow.Call("scrollTo", map[string]any{"top": 0, "left": 0})
	}

	log.Println("spaserver: done constructing page", path)

	return

}

func recentsPage(path string) {

	globalLib.SortByAvailDate()

	n, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(path, ".html"), "recent"))

	if err != nil {
		log.Println(err)
		return
	}

	if n < 1 {
		n = 1
	}

	if (n-1)*100 > globalLib.LenDistinctBooks() {

		n = globalLib.LenDistinctBooks() / 100

	}

	path = "recent" + strconv.Itoa(n) + ".html"

	mkpage(path, string(getPageData(path)))

	domHTML.Set("style", "writing-mode: horizontal-tb")

	domWindow.Call("scrollTo", map[string]any{"top": 0, "left": 0})

	log.Println("spaserver: done constructing page", path)

	return

}

func serveFile(event js.Value, params ...any) {

	path := params[0].(string)

	ext := params[1].(string)

	pparts := strings.Split(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), "_")

	authorID, bookID := pparts[1], pparts[2]

	rec := globalLib.GetRecordWithID(authorID, bookID)

	path = "files/files_" + rec.RealBookID() + "/" + rec.FileName() + "_u." + ext

	go serveFileSvc(path)

	return
}

func serveFileSvc(path string) {

	coverAndWait(domBody, 20)

	log.Println("creating", filepath.Base(path))

	data := getPageData(path)

	bk, _ := globalLib.GetBookRecord(path)

	name := ""

	switch filepath.Ext(path) {

	case ".tex", ".txt", ".html", ".json":
		name = bk.TxtFileName + "_" + strings.TrimPrefix(filepath.Ext(path), ".") + ".zip"

	default:
		name = bk.TxtFileName + filepath.Ext(path)

	}

	saveFile(createJSFile(data, name))

	log.Println("file downloaded to", name)

	uncoverElement(domBody)

	return
}

func showMore(event js.Value, params ...any) {

	elem, _ := getElementByID("more")

	elem.Set("style", "display: none")

	elem, _ = getElementByID("coveredLinks")

	elem.Set("style", "display: block")

}

func readBook(path string) {

	go readBookSvc(path)

	return
}

func readBookSvc(path string) {

	path = strings.TrimSuffix(path, ".html") + ".mono"

	mkpage(path, string(getPageData(path)))

	domHTML.Set("style", "writing-mode: vertical-rl")

	return
}

func search(event js.Value, params ...any) {

	q := domDocument.Call("getElementById", "query").Get("value").String()

	setHash("#search=" + q)

	return
}

func showSearchResult(q string) {

	q = str.PercentDecode(strings.TrimPrefix(q, "search="))

	log.Println("lookin for", q)

	mkpage(q, string(globalLib.GenSearchResults(q)))

	return

}

func randomBook(event js.Value, param ...any) {

	log.Println("finding random book")

	bk := globalLib.RandomBook()

	hash := "#books/book_" + bk.AuthorID + "_" + bk.BookID + ".html"

	setHash(hash)

	return
}

func showAbout(s string) {

	mkpage("", string(readFromResources("about.html")))

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

func readFromResources(name string) []byte {

	f, _ := resourceFiles.Open("resources/" + name)

	defer f.Close()

	return readFrom(f)

}

func mkpage(path string, data string) {

	replaceBody(data)

	domBody.Call("removeAttribute", "class")

	addPageEventListeners(path)

}

func addPageEventListeners(path string) {

	epubdl, err := getElementByID("epubdl")

	if err == nil {

		addEventListener(epubdl, "click", serveFile, path, "epub")

	}

	azw3dl, err := getElementByID("azw3dl")

	if err == nil {

		addEventListener(azw3dl, "click", serveFile, path, "azw3")

	}

	texdl, err := getElementByID("texdl")

	if err == nil {

		addEventListener(texdl, "click", serveFile, path, "tex")

	}

	txtdl, err := getElementByID("txtdl")

	if err == nil {

		addEventListener(txtdl, "click", serveFile, path, "txt")

	}

	htmldl, err := getElementByID("htmldl")

	if err == nil {

		addEventListener(htmldl, "click", serveFile, path, "html")

	}

	jsondl, err := getElementByID("jsondl")

	if err == nil {

		addEventListener(jsondl, "click", serveFile, path, "json")

	}

	aztxtdl, err := getElementByID("aztxtdl")

	if err == nil {

		addEventListener(aztxtdl, "click", serveFile, path, "zip")

	}

	showlinks, err := getElementByID("more")

	if err == nil {

		addEventListener(showlinks, "click", showMore)

	}

	rndbk, err := getElementByID("rndbk")

	if err == nil {

		addEventListener(rndbk, "click", randomBook)

	}

	queryBox, err := getElementByID("query")

	if err == nil {

		addEventListener(queryBox, "change", search)

	}

	setbtn, err := getElementByID("menubutton")

	if err == nil {

		addEventListener(setbtn, "click", settingsMenu)

	}
	return
}
