package main

import (
	"embed"
	"log"

	aozorafs "github.com/adamay909/AozoraBookcase/aozoraFS"
)

var globalLib *aozorafs.Library

var globalSettings struct {
	kids    bool
	strict  bool
	verbose bool
	clean   bool
}

//go:embed resources/*
var resourceFiles embed.FS

func main() {

	globalSettings.kids = false
	globalSettings.strict = false
	globalSettings.clean = true
	globalSettings.verbose = true

	initLibrary()

	writeConsoleLog("main: done setting up library")

	setupJS()

	log.Println("main: done setting up JS")

	loadMainPage()

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

	globalLib.Initialize("https://"+getHost(), "", globalSettings.clean, globalSettings.verbose, globalSettings.kids, globalSettings.strict)

	globalLib.FetchLibrary()
}

func loadMainPage() {

	setHash("")
	setHash("index.html")

}
