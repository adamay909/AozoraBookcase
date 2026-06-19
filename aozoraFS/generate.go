package aozorafs

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	azrconvert "github.com/adamay909/AozoraConvert/v2"
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
		SectionData []Sec
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

	return lib.cache.CreateFile("index.html", br.Bytes())
}

func (lib *Library) genRecents(name string) (f fs.File, err error) {

	type PageData struct {
		Books []*Record
		N     string
		NP    string
		NN    string
		NPT   string
		NNT   string
		ORD   string
	}

	n, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(name, ".html"), "recent"))

	if err != nil {
		log.Println(err)
		return
	}

	var P PageData
	P.Books = append(P.Books, lib.getRecents(n-1)...)

	for _, b := range P.Books {
		b.SetCategoryString(lib.Categories)
	}

	np := n - 1
	nn := n + 1
	npt := n - 10
	nnt := n + 10

	if np < 1 {
		np = 0
	}
	if (nn-1)*100 > lib.LenDistinctBooks() {
		nn = 0
	}
	if npt < 1 {
		npt = 0
	}
	if (nnt-1)*100 > lib.LenDistinctBooks() {
		nnt = 0
	}

	P.N = strconv.Itoa(n)
	P.NP = strconv.Itoa(np)
	P.NPT = strconv.Itoa(npt)
	P.NN = strconv.Itoa(nn)
	P.NNT = strconv.Itoa(nnt)
	P.ORD = strconv.Itoa((n-1)*100 + 1)

	br := new(bytes.Buffer)
	err = lib.recentT.Execute(br, P)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateFile(name, br.Bytes())
}

func (lib *Library) genAuthorPage(name string) (fs.File, error) {
	type Page struct {
		Books []*Record
		NextAuthor,
		PrevAuthor *Record
		Prefix string
	}

	var P Page

	authorID := getID(name)

	sortList(lib.booksByAuthor[authorID], byTitle)

	for _, b := range lib.booksByAuthor[authorID] {

		b.SetCategoryString(lib.Categories)
		P.Books = append(P.Books, b)
	}
	//P.Books = append(P.Books, lib.booksByAuthor[authorID]...)
	P.NextAuthor = lib.NextAuthor(P.Books[0])
	P.PrevAuthor = lib.PrevAuthor(P.Books[0])

	br := new(bytes.Buffer)
	err := lib.authorT.Execute(br, P)
	if err != nil {
		log.Println(err)
	}
	log.Println(P.Books[0].NameSei, P.Books[0].NameMei, "has ", len(P.Books), "books")
	return lib.cache.CreateFile(filepath.Join("authors", "author_"+authorID+".html"), br.Bytes())

}

func (lib *Library) genBookPage(name string) (fs.File, error) {

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
	P.B.SetCategoryString(lib.Categories)

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

	return lib.cache.CreateFile(filepath.Join("books", "book_"+authorID+"_"+bookID+".html"), br.Bytes())

}

func (lib *Library) genCategoryPage(name string) (fs.File, error) {

	type Page struct {
		Category string
		Books    []*Record
	}

	var P Page

	q := strings.TrimSuffix(strings.TrimPrefix(name, "ndc_"), ".html")

	P.Books = append(P.Books, lib.FindBooksWithMatchingCategories(q)...)

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

	return lib.cache.CreateFile(filepath.Join("categories", "ndc_"+q+".html"), br.Bytes())

}

func (lib *Library) genReadingPage(name string) (fs.File, error) {

	var rname string

	if strings.HasSuffix(name, ".mono") {
		rname = strings.TrimSuffix(name, ".mono") + ".html"
	} else {
		rname = name
	}

	book, _ := lib.getBookData(rname)

	br := new(bytes.Buffer)
	_ = lib.readingT.Execute(br, book)

	text := string(br.Bytes())

	return lib.cache.CreateFile(name, []byte(text))
}

func (lib *Library) genLatestReadPage(name string) (fs.File, error) {
	type Page struct {
		Books []*Record
	}

	var P Page

	for _, id := range lib.latestReads {
		if id == "" {
			continue
		}
		lib.consolidateRecords(id)
		P.Books = append(P.Books, lib.booksByID[id][0])
	}

	if len(P.Books) == 0 {
		P.Books = append(P.Books, lib.booksByID["056572"][0])
	}
	br := new(bytes.Buffer)
	err := lib.latestReadT.Execute(br, P)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateFile(name, br.Bytes())

}

func (lib *Library) genFavoritesPage(name string) (fs.File, error) {
	type Page struct {
		Books         []*Record
		BooksByAuthor []*Record
		BooksByTitle  []*Record
	}

	var P Page

	for id, _ := range lib.favorites {
		lib.consolidateRecords(id)
		P.Books = append(P.Books, lib.booksByID[id][0])
	}

	if len(P.Books) == 0 {
		P.Books = append(P.Books, lib.booksByID["056572"][0])
	}

	sortList(P.Books, byAuthor)

	for _, book := range P.Books {
		P.BooksByAuthor = append(P.BooksByAuthor, book)
	}

	sortList(P.Books, byTitle)

	for _, book := range P.Books {
		P.BooksByTitle = append(P.BooksByTitle, book)
	}

	br := new(bytes.Buffer)
	err := lib.favoritesT.Execute(br, P)
	if err != nil {
		log.Println(err)
	}

	return lib.cache.CreateFile(name, br.Bytes())

}

func (lib *Library) GetMonolithicHTML(name string) []string {

	var rname string

	if strings.HasSuffix(name, ".mono") {
		rname = strings.TrimSuffix(name, ".mono") + ".html"
	} else {
		rname = name
	}

	book, _ := lib.getBookData(rname)

	if book.Body == nil {
		fmt.Println("CONVERSION FAILED")
		return []string{""}
	}

	return []string{string(book.RenderMonolithicHTML()), string(book.RenderNavHTML())}
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

func (lib *Library) getBookData(name string) (book *azrconvert.Book, err error) {

	book = new(azrconvert.Book)

	bk, err := lib.GetBookRecord(name)

	if err != nil {
		return
	}
	book = lib.getBook(bk)

	return
}

func (lib *Library) generateFile(name string) (fs.File, error) {

	book, _ := lib.getBookData(name)

	var br []byte

	switch filepath.Ext(name) {

	case ".epub":
		br = book.RenderEpub()

	case ".azw3":
		br = book.RenderAZW3()

	case ".tex":
		br = book.RenderPackage("tex")

	case ".txt":
		br = book.RenderPackage("txt")

	case ".html":
		br = book.RenderPackage("html")

	case ".json":
		br = book.RenderPackage("json")

	case ".zip":
		bk, _ := lib.GetBookRecord(name)
		br, _ = download(bk.URI)

	default:
		br = book.RenderMonolithicHTML()

	}

	return lib.cache.CreateFile(name, br)
}

func getID(name string) string {

	if strings.HasSuffix(name, ".mono") {
		name = strings.TrimSuffix(name, ".mono") + ".html"
	}
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

func (lib *Library) GetID(name string) string {
	return getID(name)
}

func (lib *Library) getBook(bk *Record) *azrconvert.Book {

	data, _ := download(bk.URI)
	book := azrconvert.NewEbookFromZip(data)

	if book.Body == nil {
		return book
	}

	book.Body.ClearMetadata()

	book.Body.SetTitle(bk.Title)

	bk.TxtFileName = book.TxtFileName

	book.Body.SetDocID(bk.URI)

	if bk.Subtitle != "" {
		book.Body.SetSubtitle(bk.Subtitle)

		book.SetTitle(bk.Title + "─" + bk.Subtitle + "─")
	} else {
		book.SetTitle(bk.Title)
	}
	book.SetCreator(bk.FullName())
	book.SetPublisher("青空文庫")

	for _, c := range bk.Contributors {

		name := c.B.FullName()

		if c.B.Role == "翻訳者" {
			name = name + " 訳"
		}

		if c.B.Role == "校訂者" {
			name = name + " 校訂"
		}

		if c.B.Role == "編者" {
			name = name + " 編"
		}

		book.Body.AddContributor(name)

	}

	return book
}

func (lib *Library) GetRecordWithID(authorid, bookid string) *Record {

	for _, e := range lib.booksByID[bookid] {

		if e.AuthorID == authorid {
			return e
		}
	}

	return new(Record)
}
