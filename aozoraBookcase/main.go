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

	home, _ := os.LookupEnv("HOME")

	flag.StringVar(&root, "d", filepath.Join(home, "aozorabunko"), "directory containing server files. Defaults to $HOME/aozorabunko")

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

	mainLib := aozorafs.NewLibrary()

	ci, err := time.ParseDuration(checkint)

	if err != nil {
		log.Println(err)
	}

	mainLib.Initialize(root, clean, verbose, kids, strict, ci)

	mainLib.LoadBooklist()

	log.Println("Setting up library done.")

	http.Handle("/", http.FileServer(http.FS(mainLib)))

	http.HandleFunc("/search", aozorafs.SearchResultsHandler(mainLib))

	fmt.Println("Listening on " + iface + ":" + port)
	log.Fatal(http.ListenAndServe(iface+":"+port, nil))

}
