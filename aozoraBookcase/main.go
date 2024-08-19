package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	aozorafs "github.com/adamay909/AozoraBookcase/aozoraFS"
)

var src, root string
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

	flag.StringVar(&src, "src", "https://localhost:8888", "root url of aozorabunko's file")

	flag.Parse()

	signal.Ignore(syscall.SIGHUP)
}

func main() {

	_, err := url.Parse(src)

	if err != nil {
		log.Println("src needs to be a valid URL")
		return
	}

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
	os.RemoveAll(filepath.Join(root, "library"))

	mainLib.SetCache(NewDiskFS(filepath.Join(root, "library")))

	aozorafs.SetDownloader(DownloadFile)

	aozorafs.SetHeader(GetHeader)

	SetTemplates(mainLib)

	mainLib.Initialize(src, root, clean, verbose, kids, strict, ci)

	mainLib.LoadBooklist()

	log.Println("Setting up library done.")

	http.Handle("/", http.FileServer(http.FS(mainLib)))

	http.HandleFunc("/search", SearchResultsHandler(mainLib))

	http.HandleFunc("/random", RandomBookHandler(mainLib))

	fmt.Println("Listening on " + iface + ":" + port)
	log.Fatal(http.ListenAndServe(iface+":"+port, nil))

}

// SearchResultsHandler is a handler function for search results. To be used with http.
func SearchResultsHandler(lib *aozorafs.Library) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		qs := r.Form["query"]

		if len(qs) > 0 {
			log.Println("searching for", qs[0])
		} else {
			return
		}

		w.Write(lib.GenSearchResults(qs[0]))

		return
	}
}

// RandomBook returns a random book from the library
func RandomBookHandler(lib *aozorafs.Library) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Write(lib.GenRandomBook())
		return
	}
}
