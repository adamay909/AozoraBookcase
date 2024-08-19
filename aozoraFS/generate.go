package aozorafs

import (
	"bytes"
	"errors"
	"io/fs"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/adamay909/AozoraConvert/azrconvert"
)

func jpSortOrder() []rune {
	return []rune("あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをん")
}

func (lib *Library) genMainIndex() (fs.File, error) {

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

	for _, c := range jpSortOrder() {
		Page.SectionData = append(Page.SectionData, Sec{string(c), lib.getAuthorsByInitial(string(c))})
	}

	br := new(bytes.Buffer)
	err := lib.indexT.Execute(br, Page)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateEphemeral("index.html", br.Bytes())
}

func (lib *Library) genRecents() (fs.File, error) {

	type PageData struct {
		Books []*Record
	}

	log.Println("Creating list of recent texts.")
	var P PageData
	P.Books = append(P.Books, lib.getRecents(1000)...)

	br := new(bytes.Buffer)
	err := lib.recentT.Execute(br, P)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateEphemeral("recent.html", br.Bytes())
}

func genAuthorPage(lib *Library, name string) (fs.File, error) {
	type Page struct {
		Books []*Record
		NextAuthor,
		PrevAuthor *Record
		Prefix string
	}

	var P Page

	authorID := getID(name)
	log.Println("looking for ", authorID)

	P.Books = append(P.Books, lib.getBooksByAuthor(authorID)...)
	log.Println("found ", len(P.Books), "books by", authorID)
	P.NextAuthor = lib.NextAuthor(P.Books[0])
	P.PrevAuthor = lib.PrevAuthor(P.Books[0])

	br := new(bytes.Buffer)
	err := lib.authorT.Execute(br, P)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateEphemeral(filepath.Join("authors", "author_"+authorID+".html"), br.Bytes())

}

func genBookPage(lib *Library, name string) (fs.File, error) {

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

	br := new(bytes.Buffer)
	err := lib.bookT.Execute(br, P)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateEphemeral(filepath.Join("books", "book_"+authorID+"_"+bookID+".html"), br.Bytes())

}

func genCategoryPage(lib *Library, name string) (fs.File, error) {

	type Page struct {
		Category string
		Books    []*Record
	}

	var P Page

	q := strings.TrimSuffix(strings.TrimPrefix(name, "ndc_"), ".html")

	log.Println("making category page for", q)

	P.Books = append(P.Books, lib.FindBooksWithMatchingCategories(q)...)

	log.Println("found", len(P.Books), "items")
	//	P.Category = ndcmap()[q]
	P.Category = ndcmap()[q[:1]]
	if len(q) > 1 {
		P.Category = P.Category + " : " + ndcmap()[q[:2]]
	}
	if q[:1] == "9" && len(q) > 2 {
		P.Category = P.Category + " : " + ndcmap()[q[:3]]
	}

	br := new(bytes.Buffer)
	err := lib.categoryT.Execute(br, P)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateEphemeral(filepath.Join("categories", "ndc_"+q+".html"), br.Bytes())

}

func generateFile(lib *Library, name string) (fs.File, error) {

	bookID := getID(name)
	bk := lib.getBookRecordSimple(bookID)

	if bk.BookID != bookID {
		err := errors.New("book not found: " + name)
		return *new(LibFile), err
	}

	zn := strings.TrimSuffix(name, filepath.Ext(name)) + `.zip`

	book := new(azrconvert.Book)

	if lib.cache.Exists(zn) {
		log.Println("generating file from local material.")
		book = lib.getBookFromZip(zn)
	} else {
		book = lib.getBook(bk)
		lib.cache.CreateFile(zn, book.RenderWebpagePackage())
	}

	var br []byte

	switch filepath.Ext(name) {

	case ".epub":
		br = book.RenderEpub()

	case ".azw3":
		br = book.RenderAZW3()

	case ".mono":
		br = book.RenderMonolithicHTML()

	default:
		br = book.RenderWebpage()
		for _, file := range book.Files {
			name1 := filepath.Join(filepath.Dir(name), file.Name)
			lib.cache.CreateEphemeral(name1, file.Data)
		}
	}

	return lib.cache.CreateEphemeral(name, br)
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

func (lib *Library) getBookFromZip(name string) *azrconvert.Book {

	f, err := lib.cache.Open(name)

	info, err := f.Stat()
	if err != nil {
		log.Println(err)

		return new(azrconvert.Book)
	}

	d := make([]byte, info.Size())

	_, err = f.Read(d)

	if err != nil {
		log.Println(err)

		return new(azrconvert.Book)
	}

	return azrconvert.NewBookFromZip(d)
}

func (lib *Library) getBook(bk *Record) *azrconvert.Book {

	path, _ := url.Parse(bk.URI)

	d := download(path)
	book := azrconvert.NewBook()
	book.SetURI(path.String())
	book.GetBookFrom(d)
	if bk.Subtitle != "" {
		book.SetTitle(bk.Title + "─" + bk.Subtitle + "─")
	} else {
		book.SetTitle(bk.Title)
	}
	book.SetCreator(bk.FullName())
	book.SetPublisher("青空文庫")
	book.GenTitlePage()
	return book
}
