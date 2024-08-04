package aozorafs

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adamay909/AozoraConvert/azrconvert"
)

func jpSortOrder() []rune {
	return []rune("あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをん")
}

func (lib *Library) genMainIndex() {

	type Sec struct {
		Char string
		List []*Record
	}
	type PageData struct {
		Prefix      string
		Files       int
		Authors     int
		SectionData []struct {
			Char string
			List []*Record
		}
	}

	var Page PageData

	f, _ := os.Create(filepath.Join(lib.cache, "index.html"))
	defer f.Close()

	for _, c := range jpSortOrder() {
		Page.SectionData = append(Page.SectionData, Sec{string(c), lib.getAuthorsByInitial(string(c))})
	}

	err := lib.indexT.Execute(f, Page)
	if err != nil {
		log.Println(err)
	}

	f.Sync()
}

func (lib *Library) genRecents() {

	type PageData struct {
		Books []*Record
	}

	log.Println("Creating list of recent texts.")
	var P PageData
	P.Books = append(P.Books, lib.getRecents(100)...)

	f, _ := os.Create(filepath.Join(lib.cache, "recent.html"))

	defer f.Close()
	err := lib.recentT.Execute(f, P)

	if err != nil {
		log.Println(err)
	}

	f.Sync()
	log.Println("Created list of recent texts.")
}

func genAuthorPage(lib *Library, name string) {
	type Page struct {
		Books []*Record
		NextAuthor,
		PrevAuthor *Record
		Prefix string
	}

	var P Page

	os.Chdir(lib.resources)
	authorID := getID(name)
	log.Println("looking for ", authorID)

	P.Books = append(P.Books, lib.getBooksByAuthor(authorID)...)
	log.Println("found ", len(P.Books), "books by", authorID)
	P.NextAuthor = lib.NextAuthor(P.Books[0])
	P.PrevAuthor = lib.PrevAuthor(P.Books[0])

	f, err := os.Create(filepath.Join(lib.cache, "authors", "author_"+authorID+".html"))
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	err = lib.authorT.Execute(f, P)
	if err != nil {
		log.Println(err)
	}

	f.Sync()

}

func genBookPage(lib *Library, name string) {

	type Page struct {
		B *Record
		PrevBook,
		NextBook,
		PrevAuthor,
		NextAuthor *Record
		Prefix string
	}
	var P Page

	id := strings.Split(name, `_`)

	authorID := id[1]
	bookID := strings.TrimSuffix(id[2], ".html")

	var k int

	P.B, k = lib.getBookRecord(authorID, bookID)

	if P.B.BookID != bookID {
		log.Println("book not found:", name)
		log.Println("\t resorting to first book in catalog")
	}

	if k == 0 {
		P.PrevBook = lib.booklist[len(lib.booklist)-1]
	} else {
		P.PrevBook = lib.booklist[k-1]
	}

	if k == len(lib.booklist)-1 {
		P.NextBook = lib.booklist[0]
	} else {
		P.NextBook = lib.booklist[k+1]
	}

	P.NextAuthor = lib.NextAuthor(P.B)
	P.PrevAuthor = lib.PrevAuthor(P.B)

	f, _ := os.Create(filepath.Join(lib.cache, "books", "book_"+authorID+"_"+bookID+".html"))
	defer f.Close()

	err := lib.bookT.Execute(f, P)

	if err != nil {
		log.Println(err)
	}

	f.Sync()

}

func genCategoryPage(lib *Library, name string) {

	type Page struct {
		Category string
		Books    []*Record
	}

	var P Page

	q := strings.TrimSuffix(strings.TrimPrefix(name, "ndc_"), ".html")

	log.Println("making category page for", q)

	os.Chdir(lib.resources)

	P.Books = append(P.Books, lib.FindMatchingCategories(q)...)

	log.Println("found", len(P.Books), "items")
	P.Category = ndcmap()[q]

	f, err := os.Create(filepath.Join(lib.cache, "categories", "ndc_"+q+".html"))
	if err != nil {
		log.Println("1", err)
		os.Exit(1)
	}
	defer f.Close()

	err = lib.categoryT.Execute(f, P)
	if err != nil {
		log.Println("2", err)
		os.Exit(1)
	}

	f.Sync()

}

func generateFiles(lib *Library, name string) {

	bookID := getID(name)
	bk := lib.getBookRecordSimple(bookID)

	if bk.BookID != bookID {
		log.Println("book not found:", name)
		return
	}

	path := lib.cache
	book := getBook(lib, bk)
	id := bk.BookID
	err := os.Mkdir(filepath.Join(path, "files", "files_"+id), 0766)
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile(filepath.Join(path, "files", "files_"+id, bk.FileName()+"_u.html"), book.RenderWebpage(), 0644)
	if err != nil {
		log.Println(err)
	}
	for _, file := range book.Files {
		err = os.WriteFile(filepath.Join(path, "files", "files_"+id, file.Name), file.Data, 0644)
		if err != nil {
			log.Println(err)
		}
	}

	writeEpub(filepath.Join(path, "files", "files_"+id, bk.FileName()+"_u.epub"), book)

	writeAZW3(filepath.Join(path, "files", "files_"+id, bk.FileName()+"_u.azw3"), book)

	return

}

func getID(name string) string {

	dir := filepath.Dir(name)
	switch {
	case strings.HasPrefix(dir, "files/files_"):
		name := strings.TrimSuffix(filepath.Base(name), "_u"+filepath.Ext(name))
		id := strings.Split(name, "_")
		for len(id[0]) < 6 {
			id[0] = "0" + id[0]
		}
		return id[0]
	default:
		name = strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
		id := strings.Split(name, "_")
		if len(id) != 2 {
			return ""
		}
		return id[1]
	}
	return ""
}

func getBook(lib *Library, bk *Record) *azrconvert.Book {

	path, _ := url.Parse(lib.src + bk.URI)
	d := downloadFile(path)
	book := azrconvert.NewBookFrom(d)
	book.SetURI(path.String())
	if bk.Subtitle != "" {
		book.SetTitle(bk.Title + "─" + bk.Subtitle + "─")
	} else {
		book.SetTitle(bk.Title)
	}
	book.SetCreator(bk.FullName())
	book.SetPublisher("青空文庫")
	book.AddFiles()
	return book
}

func downloadFile(path *url.URL) []byte {
	if path.Host == "" {
		return getLocalFile(path.Path)
	}
	r, err := http.Get(path.String())
	if err != nil {
		log.Println("server download Remote:", err)
	}
	data, _ := io.ReadAll(r.Body)
	r.Body.Close()
	return data
}

func getLocalFile(path string) (data []byte) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("server retrieving local file: ", err)
	}
	return
}

func writeEpub(fname string, book *azrconvert.Book) {

	err := os.WriteFile(fname, book.RenderEpub(), 0644)
	if err != nil {
		log.Println(err)
	}
	return
}

func writeAZW3(fname string, book *azrconvert.Book) {

	err := os.WriteFile(fname, book.RenderAZW3(), 0644)
	if err != nil {
		log.Println(err)
	}
	return
}

func genAZW3(fname string, limit chan int) {

	calibre := "/usr/bin/ebook-convert"

	src := fname

	dest := strings.TrimSuffix(src, ".epub") + ".azw3"

	cmd := exec.Command(calibre, src, dest)

	err := cmd.Run()

	if err != nil {
		log.Println("calibre azw3 fail", fname, "Error was", err)
	}
	_ = <-limit

	return
}

/*
func calibreConvert(c chan string) {
	limit := make(chan int, 12)

	for {
		fname := <-c
		limit <- 1
		go genAZW3(fname, limit)
	}
}
*/
