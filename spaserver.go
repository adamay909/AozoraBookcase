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

var hashHandler map[string]handleFunc

var prefixes []string

func setupJS() {

	setupButtonListener()
	setupInputListener()
	setupHashHandlers()
	setupHashListeners()

}

func spaserver(event js.Value, params ...any) {

	js.Global().Call("detachReaderListeners")
	hash := str.PercentDecode(getHash())

	for _, p := range prefixes {
		if strings.HasPrefix(hash, p) {
			hashHandler[p](strings.TrimPrefix(hash, "#"))
			return
		}
	}
	return
}

func setupHashHandlers() {

	setHashHandler("#index.html", mainPages)

	setHashHandler("#authors", mainPages)

	setHashHandler("#books", mainPages)

	setHashHandler("#categories", mainPages)

	setHashHandler("#recent", recentsPage)

	setHashHandler("#read", readBook)

	setHashHandler("#search=", showSearchResult)

	setHashHandler("#latestReads", latestReads)

	setHashHandler("#favorites", showFavorites)
}

func setupHashListeners() {

	addEventListener(domWindow, "hashchange", spaserver)
}

func setHashHandler(prefix string, f handleFunc) {
	if len(hashHandler) == 0 {
		hashHandler = make(map[string]handleFunc)
	}
	hashHandler[prefix] = f
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

func latestReads(name string) {
	log.Println("request list of recently read books")
	go func() {
		datajs, err := jsAwait(domWindow.Call("loadFileFromIDB", name))
		var list []string
		if err == nil {
			log.Println("found list", datajs[0].String())
			rawlist := strings.Split(datajs[0].String(), `,`)
			list = cleanUpList(rawlist)
		}
		globalLib.SetLatestReads(list)
		mkpage(name, string(getPageData(name)))
	}()
}

func showFavorites(name string) {
	log.Println("request list of favorites:", name)
	go func() {
		datajs, err := jsAwait(domWindow.Call("loadFileFromIDB", name))
		var list []string
		if err == nil {
			log.Println("found list", datajs[0].String())
			rawlist := strings.Split(datajs[0].String(), `,`)
			list = cleanUpList(rawlist)
		}
		globalLib.SetFavorites(list)
		mkpage(name, string(getPageData(name)))
	}()
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

func serveFile(path, ext string) {
	pparts := strings.Split(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), "_")
	authorID, bookID := pparts[1], pparts[2]
	rec := globalLib.GetRecordWithID(authorID, bookID)
	path = "files/files_" + rec.RealBookID() + "/" + rec.FileName() + "_u." + ext
	go serveFileSvc(path)
	return
}

func serveFileSvc(path string) {

	coverScreen()

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

	//uncoverElement(domMainBody)
	uncoverScreen()
	return
}

func _showMore(event js.Value, params ...any) {

	elem, _ := getElementByID("more")

	elem.Set("style", "display: none")

	elem, _ = getElementByID("coveredLinks")

	elem.Set("style", "display: block")

}

func showMore() {

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
	/*
		defer func() {

			if r := recover(); r != nil {
				log.Println("something WRONG. Can't read!")
				showErrorMsg()
				return
			}
		}()
	*/
	var htmldata, navdata string

	coverScreen()

	bookid := globalLib.GetID(path)
	if bookid == "" {
		log.Println("couldn't find id") //this should never happen
		return
	}

	mkpage(path, string(getPageData(path)))

	found, _ := jsAwait(domWindow.Call("fileExists", bookid))
	if !found[0].Bool() {
		log.Println("requesting raw data")
		data := getBookText(path)
		htmldata = data[0]
		navdata = data[1]
		domTemplate.Set("innerHTML", htmldata)
		domNavTemplate.Set("innerHTML", navdata)
	} else {
		log.Println("Found processed material from previous use.")
	}

	uncoverScreen()

	//call the paginator
	js.Global().Call("initAozoraReader", js.ValueOf(bookid))

	//check for star
	toggle, err := getElementByID("favoriteToggle")
	if err == nil {
		if globalLib.IsFavorite(bookid) {
			toggle.Set("checked", true)
		} else {
			toggle.Set("checked", false)
		}
	}

	//add to recently read list
	var list []string
	datajs, err := jsAwait(domWindow.Call("loadFileFromIDB", "latestReads.html"))
	if err == nil {
		list = strings.Split(bookid+`,`+datajs[0].String(), `,`)
		list = cleanUpList(list)
		if len(list) > 200 {
			list = list[:200]
		}
		jsAwait(domWindow.Call("saveFileToIDB", "latestReads.html", strings.Join(list, `,`)))
	}

	return
}

func getBookText(name string) []string {

	return globalLib.GetMonolithicHTML(name)

}

/*
func search(event js.Value, params ...any) {

		q := domDocument.Call("getElementById", "query").Get("value").String()

		setHash("#search=" + q)

		return
	}
*/
func showSearchResult(q string) {

	q = str.PercentDecode(strings.TrimPrefix(q, "search="))

	log.Println("looking for", q)

	mkpage(q, string(globalLib.GenSearchResults(q)))

	return

}

func _randomBook(event js.Value, param ...any) {

	log.Println("finding random book")

	bk := globalLib.RandomBook()

	hash := "#books/book_" + bk.AuthorID + "_" + bk.BookID + ".html"

	setHash(hash)

	return
}

func randomBook() {

	log.Println("finding random book")

	bk := globalLib.RandomBook()

	hash := "#books/book_" + bk.AuthorID + "_" + bk.BookID + ".html"

	setHash(hash)

	return
}

func showAbout() {

	elem, _ := getElementByID("about")
	data := string(readFromResources("about.html"))
	elem.Set("innerHTML", data)
	elem.Set("style", "display: block")
	elem.Call("scrollTo", "0", "0")

	return

}
func closeAbout() {

	elem, _ := getElementByID("about")
	//data := string(readFromResources("about.html"))
	//	elem.Set("innerHTML", data)
	elem.Set("style", "display: none")

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

	uncoverScreen()
	hideErrorMsg()
	replaceBody(data)

	domMainBody.Call("removeAttribute", "class")

	//addPageEventListeners(path)

}

func cleanUpList(list []string) []string {
	var resp []string
	found := make(map[string]struct{})
	for _, item := range list {
		if item == "" {
			continue
		}
		if _, ok := found[item]; ok {
			continue
		}
		resp = append(resp, item)
		found[item] = struct{}{}
	}
	return resp
}

func setupButtonListener() {
	addEventListener(domBody, "click", buttonClickHandler)
}

func setupInputListener() {
	addEventListener(domBody, "change", inputHandler)
}

func inputHandler(event js.Value, params ...any) {
	//inputID := domWindow.Call("eventTargetID", event, "input").String()
	inputID := eventTargetID(event, "input")

	switch inputID {
	case "query":
		q := eventTargetElem(event, "input").Get("value").String()
		setHash("#search=" + q)
	case "favoriteToggle":
		state := eventTargetElem(event, "input").Get("checked").Bool()
		toggleFavorite(state)
	}
}

// central place for handling buttons. One handler function for
// all buttons on site
func buttonClickHandler(event js.Value, params ...any) {
	buttonId := eventTargetID(event, "button")
	path := getHash()
	switch buttonId {
	case "epubdl":
		serveFile(path, "epub")
	case "azw3dl":
		serveFile(path, "azw3")
	case "texdl":
		serveFile(path, "tex")
	case "txtdl":
		serveFile(path, "txt")
	case "htmldl":
		serveFile(path, "html")
	case "jsondl":
		serveFile(path, "json")
	case "aztxtdl":
		serveFile(path, "zip")
	case "rndbk":
		randomBook()
	case "menubutton":
		settingsMenu()
	case "homebutton", "homebuttonMobile":
		setHash("")
		setHash("index.html")
	case "clearStorage":
		clearStorage()
	case "showAbout":
		showAbout()
	case "closeAbout1", "closeAbout2":
		closeAbout()

	case "more":
		showMore()
	}
}

func toggleFavorite(checked bool) {
	path := strings.TrimPrefix(getHash(), "#")
	bookid := globalLib.GetID(path)
	fileName := "favorites.html"
	go func() {
		var list []string
		datajs, err := jsAwait(domWindow.Call("loadFileFromIDB", fileName))
		if err == nil {
			if checked {
				list = strings.Split(bookid+`,`+datajs[0].String(), `,`)
			} else {
				list = removeFromList(strings.Split(datajs[0].String(), `,`), bookid)
			}
			list = cleanUpList(list)
			if len(list) > 200 {
				list = list[:199]
			}
			jsAwait(domWindow.Call("saveFileToIDB", fileName, strings.Join(list, `,`)))
		}
	}()
}

func removeFromList(original []string, item string) (resp []string) {
	for _, elem := range original {
		if elem == item {
			continue
		}
		resp = append(resp, elem)
	}
	return resp
}

func clearStorage() {
	go func() {
		jsAwait(domWindow.Call("clearIDB"))
	}()
}
