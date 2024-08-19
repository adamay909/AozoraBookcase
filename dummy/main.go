package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	aozorafs "github.com/adamay909/AozoraBookcase/aozoraFS"
)

var root string
var clean, verbose, kids, strict bool
var iface, port string
var checkint string

func init() {

	flag.StringVar(&root, "d", "aozorabunko", "directory containing server files. Defaults to $HOME/aozorabunko. Must be a relative path and is interpreted as relative to $HOME.")

	flag.BoolVar(&clean, "c", false, "force re-downloading of all data from Aozora Bunko")
	flag.BoolVar(&verbose, "v", false, "log to screen and file")
	flag.StringVar(&iface, "i", "", "network interface")
	flag.StringVar(&port, "p", "3333", "network interface")
	flag.BoolVar(&kids, "children", false, "start a kid's library")
	flag.BoolVar(&strict, "strict", true, "set library to show only public domain texts")
	flag.StringVar(&checkint, "refresh", "24h", "interval between library refreshes.")

	flag.Parse()

	signal.Ignore(syscall.SIGHUP)
}

func main() {

	home, _ := os.LookupEnv("HOME")

	if !rootIsSafe() {
		log.Println("directory must be specified as relative path and be contained within $HOME.")
		return
	}

	root = filepath.Join(home, root)

	setupLogging()

	mainLib := aozorafs.NewLibrary()

	ci, err := time.ParseDuration(checkint)

	if err != nil {
		log.Println(err)
	}

	SetTemplates()

	//	mainLib.SetCache(NewDiskFS(filepath.Join(home, "aozorabunko")))

	ls := new(localStorage)

	mainLib.SetCache(ls)

	aozorafs.SetDownloader(DownloadFile)

	mainLib.Initialize(root, clean, verbose, kids, strict, ci)

	mainLib.LoadBooklist()

	log.Println("Setting up library done.")

	http.Handle("/", http.FileServer(http.FS(mainLib)))

	http.HandleFunc("/search", aozorafs.SearchResultsHandler(mainLib))

	http.HandleFunc("/random", aozorafs.RandomBook(mainLib))

	fmt.Println("Listening on " + iface + ":" + port)
	log.Fatal(http.ListenAndServe(iface+":"+port, nil))

}
