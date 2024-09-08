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

	var P PageData
	P.Books = append(P.Books, lib.getRecents(100)...)

	br := new(bytes.Buffer)
	err := lib.recentT.Execute(br, P)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateFile("recent.html", br.Bytes())
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

	sortList(lib.booksByAuthor[authorID], byTitle)

	P.Books = append(P.Books, lib.booksByAuthor[authorID]...)
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

	lib.consolidateRecords(bookID)

	booklist := lib.booksByAuthor[authorID]

	sortList(booklist, byTitle)

	k := 0
	for k = 0; k < len(booklist); k++ {
		if booklist[k].BookID == bookID {
			P.B = booklist[k]
			break
		}
	}

	if k == 0 {
		P.PrevBook = lib.LastBookBy(lib.PrevAuthor(P.B))
	} else {
		P.PrevBook = booklist[k-1]
	}

	if k == len(booklist)-1 {
		P.NextBook = lib.FirstBookBy(lib.NextAuthor(P.B))
	} else {
		P.NextBook = booklist[k+1]
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

	P.Books = append(P.Books, lib.FindBooksWithMatchingCategories(q)...)

	//	P.Category = lib.Categories[q]
	P.Category = lib.Categories[q[:1]]
	if len(q) > 1 {
		P.Category = P.Category + " : " + lib.Categories[q[:2]]
	}
	if q[:1] == "9" && len(q) > 2 {
		P.Category = P.Category + " : " + lib.Categories[q[:3]]
	}

	br := new(bytes.Buffer)
	err := lib.categoryT.Execute(br, P)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateEphemeral(filepath.Join("categories", "ndc_"+q+".html"), br.Bytes())

}

func genReadingPage(lib *Library, name string) (fs.File, error) {

	var rname string
	var book *azrconvert.Book
	var err error

	if strings.HasSuffix(name, ".mono") {
		rname = strings.TrimSuffix(name, ".mono") + ".html"
	} else {
		rname = name
	}

	book, err = getBookData(lib, rname)

	var realbody string

	if strings.HasSuffix(name, "mono") {
		realbody = book.RenderBodyInnerMonolithic()
	} else {
		realbody = book.RenderBodyInner()
		for _, file := range book.Files {
			name1 := filepath.Join(filepath.Dir(rname), file.Name)
			lib.cache.CreateEphemeral(name1, file.Data)
		}
	}
	br := new(bytes.Buffer)
	err = lib.readingT.Execute(br, book)

	text := string(br.Bytes())

	text = strings.Replace(text, "!!!###TEXT###!!!", realbody, 1)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateEphemeral(name, []byte(text))
}

func (lib *Library) GetBookRecord(name string) (*Record, error) {

	var err error
	bookID := getID(name)
	bk := lib.getBookRecordSimple(bookID)

	if bk.BookID != bookID {
		err := errors.New("book not found: " + name)
		return bk, err
	}

	return bk, err
}

func getBookData(lib *Library, name string) (book *azrconvert.Book, err error) {

	book = new(azrconvert.Book)

	bk, err := lib.GetBookRecord(name)

	if err != nil {
		return
	}

	zn := strings.TrimSuffix(name, filepath.Ext(name)) + `.zip`

	if lib.cache.Exists(zn) {
		log.Println("generating file from local material.")
		book = lib.getBookFromZip(zn)
	} else {
		book = lib.getBook(bk)
		lib.cache.CreateEphemeral(zn, book.RenderWebpagePackage())
	}

	return
}

func generateFile(lib *Library, name string) (fs.File, error) {

	book, _ := getBookData(lib, name)

	var br []byte

	switch filepath.Ext(name) {

	case ".epub":
		br = book.RenderEpub()

	case ".azw3":
		br = book.RenderAZW3()

	default:
		br = book.RenderMonolithicHTML()

	}

	return lib.cache.CreateEphemeral(name, br)
}

func getID(name string) string {

	dir := filepath.Dir(name)
	if strings.HasPrefix(dir, "read") {
		dir = strings.ReplaceAll(dir, "read", "files")
	}

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
