package main

import (
	"embed"
	_ "embed" //for embedding data
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	aozorafs "github.com/adamay909/AozoraBookcase/aozoraFS"
)

var src, root string
var clean, verbose, kids, strict bool
var iface, port string

//go:embed resources/*
var resourceFiles embed.FS

func init() {

	flag.StringVar(&root, "d", "aozorabunko", "directory containing server files. Defaults to $HOME/aozorabunko. Must be a relative path and is interpreted as relative to $HOME.")

	flag.BoolVar(&clean, "c", false, "force re-downloading of all data from Aozora Bunko")
	flag.BoolVar(&verbose, "v", false, "log to screen and file")
	flag.StringVar(&iface, "i", "", "network interface")
	flag.StringVar(&port, "p", "3333", "network interface")
	flag.BoolVar(&kids, "children", false, "start a kid's library")
	flag.BoolVar(&strict, "strict", true, "set library to show only public domain texts")

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

	fsys := NewDiskFS(filepath.Join(root, "library"))

	fsys.RemoveAll()

	mainLib.SetCache(fsys)

	aozorafs.SetDownloader(DownloadFile)

	mainLib.ImportTemplates(resourceFiles)

	mainLib.Initialize(src, root, clean, verbose, kids, strict)

	mainLib.FetchLibrary()

	mainLib.SortByAvailDate()

	log.Println("Setting up library done.")

	http.Handle("/", http.FileServer(http.FS(mainLib)))

	http.HandleFunc("/search", SearchResultsHandler(mainLib))

	fmt.Println("Listening on " + iface + ":" + port)

	log.Fatal(http.ListenAndServe(iface+":"+port, nil))

}

// SearchResultsHandler is a handler function for search results. To be used with http.
func SearchResultsHandler(lib *aozorafs.Library) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		qs := r.Form["query"]

		if len(qs) > 0 {
			log.Println("search handler:", qs[0])
		} else {
			return
		}

		w.Write(lib.GenSearchResults(qs[0]))

		return
	}
}

// RandomBook returns a random book from the library
/*
func RandomBookHandler(lib *aozorafs.Library) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Write(lib.GenRandomBook())
		return
	}
}
*/
