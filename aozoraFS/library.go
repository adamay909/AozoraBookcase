/*
Package aozorafs implements a fs.FS that can be used to display the contents of Aozora Bunko. It makes the documents available as vertically formatted html as well as vertically formatted epub and azw3 (for Kindle readers). It is primarily intended to be used with aozoraBookcase.

Local files are created upon first request. Subsequent request are handled through locally available files.
*/

package aozorafs

import (
	"errors"
	"io/fs"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	ac "github.com/adamay909/AozoraConvert/v2"
)

// NewLibrary returns a new Library.
func NewLibrary() *Library {

	return new(Library)
}

/*Open makes available a file corresponding to name. This makes Library implement fs.FS so that it can be served through http's file server.
 */

func (lib *Library) Open(name string) (f fs.File, err error) {

	log.Println("requested file: ", name)

	if !isValidFileName(name) {
		name = "index.html"
	}
	f, err = createFile(lib, name)

	if err != nil {
		log.Println(err)
	}
	return f, err

}

/*
Initialize initalized the library lib to the given specifications.
  - dir is the root directory of the library on the local file system.
  - clean specifies whether or not to start with an empty library directory.
  - verbose toggles verbose logging to screen. Logging to aozora.log will always take place.
  - kids specifies toggles children's library (removes books that are not marked as children's book in the  Aozora Bunko database.
  - strict toggles whether or not to include books that are not in the public domain.
  - checkInt specifies the interval for checking for updates to the library upstream.
*/
func (lib *Library) Initialize(src string, dir string, clean, verbose, kids, strict bool) {

	u, err := url.Parse(src)

	if err != nil {
		log.Println(err)
		return
	}

	lib.src = u.Path

	lib.root = dir

	lib.setKids(kids)

	lib.setStrict(strict)

	lib.booksByID = make(map[string][]*Record)

	lib.booksByAuthor = make(map[string][]*Record)

	lib.posOfAuthor = make(map[string]int)

	lib.Categories = ndcmap()

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

func createFile(lib *Library, name string) (f fs.File, err error) {

	dir := filepath.Dir(name)
	bname := filepath.Base(name)

	switch {

	case bname == "index.html":
		f, err = lib.genMainIndex()

	case bname == "aozora.css":
		return lib.cache.CreateFile("aozora.css", []byte(ac.AozoraCSS))

	case strings.HasPrefix(bname, "author"):
		f, err = lib.genAuthorPage(name)

	case strings.HasPrefix(bname, "book"):
		f, err = lib.genBookPage(name)

	case strings.HasPrefix(bname, "ndc"):
		f, err = lib.genCategoryPage(bname)

	case strings.HasPrefix(dir, "files/files"):
		f, err = lib.generateFile(name)

	case strings.HasPrefix(dir, "read/files"):
		f, err = lib.genReadingPage(name)

	case strings.HasSuffix(bname, "latestReads.html"):
		f, err = lib.genLatestReadPage(bname)

	case strings.HasSuffix(bname, "favorites.html"):
		f, err = lib.genFavoritesPage(bname)

	case strings.HasPrefix(bname, "recent"):
		f, err = lib.genRecents(bname)

	case strings.HasPrefix(bname, "random"):
		f, err = lib.GenRandomBook()

	default:
		err = errors.New("invalid request")
		return
	}

	info, _ := f.Stat()
	log.Println("created", info.Name())

	return
}

func isValidFileName(n string) bool {

	return true

	/***********************************************************

	Need a better scalable way of determining valid file names

	***********************************************************/

	if n == "." {
		return true
	}

	if strings.HasPrefix(n, "files/files_") && len(strings.Split(n, "/")) == 3 {
		return true
	}

	if strings.HasPrefix(n, "read/files_") && len(strings.Split(n, "/")) == 3 {
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

	if strings.HasSuffix(n, "latestReads.html") {
		return true
	}

	if strings.HasSuffix(n, "favorites.html") {
		return true
	}

	if strings.HasPrefix(n, "recent") && strings.HasSuffix(n, ".html") {
		return true
	}

	if strings.HasPrefix(n, "random") && strings.HasSuffix(n, ".html") {
		return true
	}

	if strings.HasSuffix(n, "about.html") {
		return true
	}

	if len(strings.Split(n, "/")) > 1 {
		log.Println("invalid request:", n)
		return false
	}

	if n == "ebooks.css" {
		return true
	}

	if n == "readingpane.css" {
		return true
	}
	if n == "index.html" {
		return true
	}

	log.Println("invalid request:", n)
	return false
}

func (lib *Library) BookList() []*Record {
	return lib.booklist
}

func (lib *Library) URL() string {
	return lib.src
}
