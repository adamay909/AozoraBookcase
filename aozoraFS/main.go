/*
Package aozorafs implements a fs.FS that can be used to display the contents of Aozora Bunko. It makes the documents available as vertically formatted html as well as vertically formatted epub and azw3 (for Kindle readers). It is primarily intended to be used with aozoraBookcase.

Local files are created upon first request. Subsequent request are handled through locally available files.
*/
package aozorafs

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

	os.Chdir(lib.cache)

	log.Println("requested file: ", name)

	_, err = os.Stat(name)
	if os.IsNotExist(err) {
		createFile(lib, name)
	}
	os.Chdir(lib.cache)
	return os.Open(name)

}

func createFile(lib *Library, name string) {

	if !isValidFileName(name) {
		return
	}

	dir := filepath.Dir(name)
	bname := filepath.Base(name)

	switch {

	case strings.HasPrefix(bname, "author"):
		genAuthorPage(lib, name)

	case strings.HasPrefix(bname, "book"):
		genBookPage(lib, name)

	case strings.HasPrefix(bname, "ndc"):
		genCategoryPage(lib, bname)

	case strings.HasPrefix(dir, "files/files"):
		generateFiles(lib, name)

	case strings.Contains(bname, "recent"):
		lib.genRecents()

	default:
		return
	}
}

/*RefreshBooklist checks for updates and if necessary refreshes the database periodically as specified by lib.checkInt. If lib.checkInt <=0, then database is never refreshed. */
func (lib *Library) RefreshBooklist() {

	if lib.checkInterval <= 0 {
		return
	}
	for {
		if lib.UpstreamUpdated(lib.lastUpdated) {
			lib.UpdateDB()
			lib.UpdateBooklist()
			lib.updatePages()
		}
		time.Sleep(lib.checkInterval)
	}
}

func isValidFileName(n string) bool {

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

	return false
}
