/*
Package aozorafs implements a fs.FS that can be used to display the contents of Aozora Bunko. It makes the documents available as vertically formatted html as well as vertically formatted epub and azw3 (for Kindle readers). It is primarily intended to be used with aozoraBookcase.

Local files are created upon first request. Subsequent request are handled through locally available files.
*/

package aozorafs

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
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
// NewLibrary returns a new Library.
func NewLibrary() *Library {

	return new(Library)
}

/*Open makes available a file corresponding to name. This makes Library implement fs.FS so that it can be served through http's file server.
 */

func (lib *Library) Open(name string) (f fs.File, err error) {

	//block Open while booklist is being updated
	for lib.updating {
	}

	if !isValidFileName(name) {
		name = "index.html"
	}

	log.Println("requested file: ", name)

	if !lib.cache.Exists(name) {

		log.Println(name, "does not exist")

		f, err = createFile(lib, name)

		if err != nil {
			log.Println(err)
		}

	}
	return lib.cache.Open(name)

}

func (lib *Library) Initialize(src string, dir string, clean, verbose, kids, strict bool, checkInt time.Duration) {

	lib.src = src

	lib.root = dir

	os.Mkdir(lib.root, 0766)

	lib.setKids(kids)

	lib.setStrict(strict)

	lib.checkInterval = checkInt

	//	lib.initfs(clean)

	// lib.inittemplates()
}

/*
func (lib *Library) initfs(clean bool) {

	lib.saveCSS()

}
*/
/*
func (lib *Library) inittemplates() {

	lib.mainIndexTemplate()
	lib.authorpageTemplate()
	lib.bookpageTemplate()
	lib.categorypageTemplate()
	lib.recentTemplate()
	lib.randomBookTemplate()

}
*/
// LoadBooklist adds (and possibly updates) the list of books for lib.

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
			defer r.Close()
			if err != nil {
				log.Println(err)
				return
			}
			d, err = io.ReadAll(r)
			if err != nil {
				log.Println(err)
				return
			}
			break
		}
	}
	log.Println("received", len(zf), "bytes")
	return d
}

func createFile(lib *Library, name string) (f fs.File, err error) {

	if !isValidFileName(name) {
		err = errors.New("invalid request 1")
		return
	}

	dir := filepath.Dir(name)
	bname := filepath.Base(name)

	switch {

	case bname == "index.html":
		f, err = lib.genMainIndex()
		info, _ := f.Stat()
		log.Println("created", info.Name())

	case strings.HasPrefix(bname, "author"):
		f, err = genAuthorPage(lib, name)

	case strings.HasPrefix(bname, "book"):
		f, err = genBookPage(lib, name)

	case strings.HasPrefix(bname, "ndc"):
		f, err = genCategoryPage(lib, bname)

	case strings.HasPrefix(dir, "files/files"):
		f, err = generateFile(lib, name)

	case strings.Contains(bname, "recent"):
		f, err = lib.genRecents()

	default:
		err = errors.New("invalid request 2")
		return
	}

	return
}

/*RefreshBooklist checks for updates and if necessary refreshes the database periodically as specified by lib.checkInt. If lib.checkInt <=0, then database is never refreshed. */
func (lib *Library) RefreshBooklist() {

	update := func() {
		lib.UpdateDB()
		lib.UpdateBooklist()
		lib.updatePages()
	}

	if lib.checkInterval <= 0 {
		return
	}
	for {
		time.Sleep(lib.checkInterval)

		if lib.UpstreamUpdated(lib.lastUpdated) {
			update()
		}
	}
}

func isValidFileName(n string) bool {

	if n == "." {
		return true
	}

	if strings.HasPrefix(n, "files/files_") && len(strings.Split(n, "/")) == 3 {
		return true
	}

	if len(strings.Split(n, "/")) > 2 {
		return false
	}

	if strings.HasPrefix(n, "authors/author_") && strings.HasSuffix(n, ".html") {
		return true
	}

	if strings.HasPrefix(n, "books/book_") && strings.HasSuffix(n, ".html") && len(strings.Split(n, "_")) == 3 {
		return true
	}

	if strings.HasPrefix(n, "categories/ndc_") && strings.HasSuffix(n, ".html") {
		return true
	}

	if len(strings.Split(n, "/")) > 1 {
		log.Println("invalid request:", n)
		return false
	}

	if n == "ebooks.css" {
		return true
	}

	if n == "index.html" {
		return true
	}

	if n == "recent.html" {
		return true
	}

	log.Println("invalid request:", n)
	return false
}
