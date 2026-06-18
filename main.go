package main

import (
	"embed"
	"log"
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

	defer func() {

		if r := recover(); r != nil {
			log.Println("something WRONG")
			showErrorMsg()
			return
		}
	}()
	globalSettings.kids = false
	globalSettings.strict = false
	globalSettings.clean = true
	globalSettings.verbose = true
	globalSettings.jis0213 = false

	var hash string
	reloaded := domWindow.Call("reloaded").Bool()
	if !reloaded {
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

	setHash("")
	if hash != "" {
		setHash(strings.TrimPrefix(hash, "#"))
	} else {
		setHash("index.html")
	}

	//	loadMainPage()

	uncoverScreen()

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

	globalLib.FetchLibrary()

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
