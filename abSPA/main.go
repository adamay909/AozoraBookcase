package main

import (
	"log"
	"time"

	aozorafs "github.com/adamay909/AozoraBookcase/aozoraFS"
)

var globalLib *aozorafs.Library

func main() {

	s := new(localStorage)

	globalLib = aozorafs.NewLibrary()

	globalLib.SetCache(s)

	aozorafs.SetDownloader(fetchData)

	SetTemplates(globalLib)

	ci, _ := time.ParseDuration("24h")

	log.Println("starting up")

	log.Println("site URL is", getUrl())

	globalLib.Initialize("https://"+getHost(), "", false, true, false, true, ci)

	globalLib.FetchLibrary()

	log.Println("main: done setting up library")

	setHandler("#index.html", mainPages)

	setHandler("#authors", mainPages)

	setHandler("#books", mainPages)

	setHandler("#categories", mainPages)

	setHandler("#recent.html", mainPages)

	//	setHandler("#files", serveFile)

	setHandler("#read", readBook)

	//	setHandler("#random", randomBook)

	setHandler("#search=", showSearchResult)

	setHandler("#menu", showMenu)

	setupjs()

	log.Println("main: done setting up JS")

	setHash("")

	setHash("index.html")

	<-make(chan bool) //prevent exiting

}
