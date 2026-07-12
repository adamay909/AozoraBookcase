package main

//go:generate ./scripts/generate.sh

import (
	"embed"
	"log"
	"path"
	"strings"

	aozorafs "github.com/adamay909/AozoraBookcase/aozoraFS"
	ac "github.com/adamay909/AozoraConvert/v2"
)

var globalLib *aozorafs.Library

var globalSettings struct {
	kids    bool
	strict  bool
	verbose bool
	clean   bool
	jis0213 bool
}

//go:embed resources/*
var resourceFiles embed.FS

func main() {
	/*
		defer func() {

			if r := recover(); r != nil {
				log.Println("something WRONG")
				showErrorMsg()
				return
			}
		}()
	*/
	globalSettings.kids = false
	globalSettings.strict = false
	globalSettings.clean = true
	globalSettings.verbose = true
	globalSettings.jis0213 = false

	var hash string
	if !reloaded() {
		hash = domWindow.Get("location").Get("hash").String()
		log.Println("hash is:", hash, "len:", len(hash))
	}

	ac.SetFullUnicode()

	ac.SetStrict(false)

	ac.SetRubyEmph(true)

	initLibrary()

	writeConsoleLog("main: done setting up library")

	setupJS()

	log.Println("main: done setting up JS")

	hideSplashscreen()

	setHash("")
	if hash != "" {
		setHash(strings.TrimPrefix(hash, "#"))
	} else {
		setHash("index.html")
	}

	//	loadMainPage()

	<-make(chan bool) //prevent exiting

}

func initLibrary() {

	log.Println("initializing library")

	globalLib = aozorafs.NewLibrary()

	s := new(localStorage)

	s.RemoveAll()

	globalLib.SetCache(s)

	aozorafs.SetDownloader(fetchData)

	globalLib.ImportTemplates(resourceFiles)

	log.Println("starting up")

	log.Println("site URL is", getURL())

	globalLib.Initialize("https://"+getHost()+getPathname(), "", globalSettings.clean, globalSettings.verbose, globalSettings.kids, globalSettings.strict)

	var libdata []byte
	mustDownload := newVersionAvailable()

	if !mustDownload {
		rawdata, err := jsAwait(domWindow.Call("loadFileFromIDB", "librarydata.zip"))
		if err != nil {
			mustDownload = true
		} else {
			libdata = bytesOf(rawdata[0])
			log.Println("loaded from indexDB")
		}
	}
	if mustDownload {
		libdata = globalLib.FetchLibraryData()
		log.Println("fetched from server")
		//		jsAwait(domWindow.Call("saveFileToIDB", "librarydata.zip", uint8arrayOf(libdata))) //no point. depend on browser cache
	}

	//globalLib.FetchLibrary()
	globalLib.ConstructLibrary(libdata)

	setupFavorites()
}

func loadMainPage() {

	setHash("")
	setHash("index.html")

}

func setupFavorites() {
	fileName := "favorites.html"
	go func() {
		datajs, err := jsAwait(domWindow.Call("loadFileFromIDB", fileName))
		if err == nil {
			list := strings.Split(datajs[0].String(), `,`)
			globalLib.SetFavorites(list)
		}
	}()
}

// check for updated library db
func newVersionAvailable() bool {
	return true //just depend on browser cache!!
	src := globalLib.URL()
	log.Println("check for updated library")
	pathStr := path.Join(src, "sha512.sum")
	sum, err := fetchData(pathStr)
	if err != nil {
		log.Println("can't download sum")
		return true
	}

	defer func() {
		jsAwait(domWindow.Call("saveFileToIDB", "sha512.sum", string(sum)))
	}()

	oldsumjs, err := jsAwait(domWindow.Call("loadFileFromIDB", "sha512.sum"))
	if err != nil {
		log.Println("can't find old sum")
		return true
	}
	oldsum := oldsumjs[0].String()
	return oldsum != string(sum)
}
