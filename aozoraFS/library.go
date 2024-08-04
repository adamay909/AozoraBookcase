package aozorafs

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

/*
Initialize initalized the library lib to the given specifications.
  - dir is the root directory of the library on the local file system.
  - clean specifies whether or not to start with an empty library directory.
  - verbose toggles verbose logging to screen. Logging to aozora.log will always take place.
  - kids specifies toggles children's library (removes books that are not marked as children's book in the  Aozora Bunko database.
  - strict toggles whether or not to include books that are not in the public domain.
  - checkInt specifies the interval for checking for updates to the library upstream.
*/
func (lib *Library) Initialize(dir string, clean, verbose, kids, strict bool, checkInt time.Duration) {

	lib.src = "https://www.aozora.gr.jp"

	lib.root = dir
	lib.cache = filepath.Join(lib.root, "library/")
	lib.resources = filepath.Join(lib.root, "resources/")

	os.Mkdir(lib.root, 0766)

	lib.setupLogging(verbose)

	lib.setKids(kids)

	lib.setStrict(strict)

	lib.checkInterval = checkInt

	lib.initfs(clean)

	lib.inittemplates()
}

func (lib *Library) initfs(clean bool) {

	if clean {
		os.RemoveAll(lib.cache)
		os.Mkdir(lib.cache, 0766)
		os.Mkdir(lib.cache+"/authors", 0766)
		os.Mkdir(lib.cache+"/books", 0766)
		os.Mkdir(lib.cache+"/files", 0766)
		os.Mkdir(lib.cache+"/categories", 0766)
	}
	lib.saveCSS()

}

func (lib *Library) inittemplates() {

	lib.mainIndexTemplate()
	lib.authorpageTemplate()
	lib.bookpageTemplate()
	lib.categorypageTemplate()
	lib.recentTemplate()
	lib.randomBookTemplate()

}

// LoadBooklist adds (and possibly updates) the list of books for lib.

func (lib *Library) setupLogging(verbose bool) {
	var w io.Writer

	f, err := os.OpenFile(filepath.Join(lib.root, "aozora.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Println(err)
	}

	if verbose {
		w = io.MultiWriter(f, os.Stdout)
	} else {
		w = f
	}
	log.SetOutput(w)
	return
}

// setKids sets lib to be a kids library if val is true..
func (lib *Library) setKids(val bool) {

	lib.kids = val

	return

}

// setStrict sets lib to show only documents that are in the public domain.
func (lib *Library) setStrict(val bool) {
	lib.strict = val

	return
}

func unzip(zf []byte) (d []byte) {

	br := bytes.NewReader(zf)

	f, _ := zip.NewReader(br, int64(len(zf)))

	for _, z := range f.File {
		if filepath.Ext(z.Name) == ".csv" {
			r, err := z.Open()
			if err != nil {
				log.Println(err)
				return
			}
			d, err = io.ReadAll(r)
			if err != nil {
				log.Println(err)
				return
			}
			r.Close()
			break
		}
	}
	log.Println("received", len(zf), "bytes")
	return d
}
