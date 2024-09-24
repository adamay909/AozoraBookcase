package main

import (
	"log"
	"time"

	aozorafs "github.com/adamay909/AozoraBookcase/aozoraFS"
)

var globalLib *aozorafs.Library

var globalSettings struct {
	kids bool
}

func main() {

	globalSettings.kids = false

	initLibrary()

	log.Println("main: done setting up library")

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

	SetTemplates(globalLib)

	ci, _ := time.ParseDuration("24h")

	log.Println("starting up")

	log.Println("site URL is", getUrl())

	globalLib.Initialize("https://"+getHost(), "", false, true, globalSettings.kids, true, ci)

	globalLib.FetchLibrary()
}

func loadMainPage() {

	setHash("")

	setHash("index.html")

}
