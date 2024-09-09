package aozorafs

import (
	"encoding/csv"
	"log"
	"math/rand"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adamay909/AozoraBookcase/zipfs"
)

// LoadBooklist adds (and possibly updates) the list of books for lib.
func (lib *Library) LoadBooklist() {

	fi, err := lib.cache.Stat("aozoradata.zip")

	switch {
	case err != nil:
		log.Println("no local aozora database")
		lib.UpdateDB()

	default:
		if lib.UpstreamUpdated(fi.ModTime()) {
			lib.UpdateDB()
		}

	}

	lib.UpdateBooklist()

	lib.updatePages()

	go lib.RefreshBooklist()

}

/*UpstreamUpdated reports whether the upstream database has been updated since it was last updated locally.
 */
func (lib *Library) UpstreamUpdated(t time.Time) bool {

	loc, err := url.JoinPath(lib.src, "/index_pages", "list_person_all_extended_utf8.zip")

	if err != nil {
		log.Println(err)
		return false
	}

	path, _ := url.Parse(loc)

	r := getHeader(path)

	m, err := time.Parse(time.RFC1123, get(r, "Last-Modified"))

	return m.After(t)
}

/*UpdateDB downloads the database from upstream.*/
func (lib *Library) UpdateDB() {

	pathString, err := url.JoinPath(lib.src, "/index_pages", "list_person_all_extended_utf8.zip")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("requesting db", pathString)

	path, _ := url.Parse(pathString)

	data := download(path)

	f, err := lib.cache.CreateFile("aozoradata.zip", data)

	defer f.Close()

	if err != nil {
		log.Println(err)
		return
	}
}

/*UpdateBooklist updates the booklist of lib from the locally available database.*/
func (lib *Library) UpdateBooklist() {

	log.Println("updating book list")

	lib.updating = true

	zf, err := zipfs.OpenZipArchive(lib.cache, "aozoradata.zip")

	defer zf.CloseArchive()

	if err != nil {
		log.Println(err)
		return
	}

	zd := zf.ReadMust("list_person_all_extended_utf8.csv")

	lib.getBooklist(zd)
	lib.setupAuthorsList()
	lib.updating = false
	return
}

func (lib *Library) FetchLibrary() {

	log.Println("getting library catalog information")

	pathString, err := url.JoinPath(lib.src, "/index_pages", "list_person_all_extended_utf8.zip")

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("requesting db", pathString)

	path, _ := url.Parse(pathString)

	data := download(path)

	za, _ := zipfs.ZipArchiveFromData(data)

	defer za.CloseArchive()

	zd := za.ReadMust("list_person_all_extended_utf8.csv")

	lib.getBooklist(zd)

	lib.setupAuthorsList()

	return
}

func (lib *Library) setupAuthorsList() {

	for _, e := range lib.booksByAuthor {
		lib.authorsSorted = append(lib.authorsSorted, e[0])
	}

	sortList(lib.authorsSorted, byAuthor)

	for k, b := range lib.authorsSorted {

		lib.posOfAuthor[b.AuthorID] = k

	}

	return
}

func (lib *Library) consolidateRecords(bookID string) {

	log.Println("found", len(lib.booksByID[bookID]), "books with ID", bookID)

	if lib.booksByID[bookID][0].consolidated {

		return
	}

	for _, l := range lib.booksByID[bookID] {

		lib.booksByID[bookID][0].Contributors = append(lib.booksByID[bookID][0].Contributors, ContribRole{l.Role, l.AuthorID, l})
	}

	sort.Slice(lib.booksByID[bookID][0].Contributors, byRole(lib.booksByID[bookID][0].Contributors))

	for k, e := range lib.booksByID[bookID] {

		e.consolidated = true

		if k == 0 {
			continue
		}

		e.Contributors = nil

		e.Contributors = append(e.Contributors, lib.booksByID[bookID][0].Contributors...)

	}

	return

}

func (lib *Library) getBooklist(d []byte) {

	rows := strings.Split(string(d), "\n")
	log.Println("database has ", len(rows), "entries")

	headings, err := csv.NewReader(strings.NewReader(rows[0])).Read()
	if err != nil {
		log.Println("error reading Aozora Bunko database:", err)
		return
	}

	col := make(map[string]int)

	//get column number for each heading
	for i, h := range headings {
		col[h] = i
	}

	var book *Record

	//read into records
	for i := 1; i < len(rows)-1; i++ {

		r := rows[i]
		cells, _ := csv.NewReader(strings.NewReader(r)).Read()

		if len(cells) == 0 {
			break
		}

		book = new(Record)

		book.WorkCopyright = cells[col["作品著作権フラグ"]]
		book.AuthorCopyright = cells[col["人物著作権フラグ"]]

		if lib.strict {
			if book.WorkCopyright == "あり" || book.AuthorCopyright == "あり" {
				continue
			}
		}

		book.NDC = cells[col["分類番号"]]
		book.setCategory(lib.Categories)

		if lib.kids {
			if !book.isChildrensBook() {
				continue
			}
		}

		book.BookID = cells[col["作品ID"]]
		book.Title = cells[col["作品名"]]
		book.TitleY = cells[col["作品名読み"]]
		book.TitleSort = cells[col["ソート用読み"]]
		book.Subtitle = cells[col["副題"]]
		book.SubtitleY = cells[col["副題読み"]]
		book.OriginalTitle = cells[col["原題"]]
		book.PublDate = cells[col["初出"]]
		book.KanaZukai = cells[col["文字遣い種別"]]
		book.FirstAvailable = cells[col["公開日"]]
		book.ModTime = cells[col["最終更新日"]]
		book.AuthorID = cells[col["人物ID"]]
		book.NameSei = cells[col["姓"]]
		book.NameMei = cells[col["名"]]
		book.NameSeiY = cells[col["姓読み"]]
		book.NameMeiY = cells[col["名読み"]]
		book.NameSeiSort = cells[col["姓読みソート用"]]
		book.NameMeiSort = cells[col["名読みソート用"]]
		book.NameSeiR = cells[col["姓ローマ字"]]
		book.NameMeiR = cells[col["名ローマ字"]]
		book.Role = cells[col["役割フラグ"]]
		book.DoBirth = cells[col["生年月日"]]
		book.DoDeath = cells[col["没年月日"]]
		book.URI, _ = url.JoinPath(lib.src, strings.TrimPrefix(cells[col["XHTML/HTMLファイルURL"]], "https://www.aozora.gr.jp"))

		lib.booklist = append(lib.booklist, book)

	}
	rows = nil

	for _, e := range lib.booklist {
		lib.booksByID[e.BookID] = append(lib.booksByID[e.BookID], e)
		lib.booksByAuthor[e.AuthorID] = append(lib.booksByAuthor[e.AuthorID], e)
	}

	lib.nextrandom = rand.Intn(len(lib.booklist))

	log.Println("finished parsing db.")
	return
}

func (lib *Library) updatePages() {

	log.Println("Updating pages")

	lib.removePages(`index.html`, `recent.html`)
	lib.removePages(lib.allUpdatedPages()...)
	lib.lastUpdated = time.Now()
	log.Println("pages updated")
}

func (lib *Library) allUpdatedPages() (list []string) {

	for _, b := range lib.booklist {

		var fnames []string

		fnames = append(fnames, `authors/author_`+b.AuthorID+`.html`)
		fnames = append(fnames, `books/book_`+b.AuthorID+`_`+b.BookID+`.html`)

		for _, c := range b.Categories {

			fnames = append(fnames, `categories/ndc_`+c[0]+`.html`)
			fnames = append(fnames, `categories/ndc_`+c[1]+`.html`)

		}

		for _, f := range fnames {
			info, err := lib.cache.Stat(f)
			if err == nil {

				if !info.ModTime().Before(lib.lastUpdated) {
					log.Println(f, "needs updating")
					list = append(list, f)
				}
			}
		}

	}

	return list
}

func (lib *Library) removePages(pages ...string) {

	for _, p := range pages {

		lib.cache.RemoveAll(filepath.Join(p))

	}

	return
}

/*ReadBooklist is for retrieving the list of books stored in lib.*/
func (lib *Library) ReadBooklist() (o []*Record) {

	o = append(o, lib.booklist...)

	return

}

/*
WriteBooklist is for adding l as the booklist of lib.
*/
func (lib *Library) WriteBooklist(l []*Record) {

	lib.booklist = nil

	lib.booklist = append(lib.booklist, l...)

	return

}

func (lib *Library) getRecents(n int) []*Record {

	var list []*Record

	listed := make(map[string]bool)

	for _, e := range lib.booklist {

		if _, ok := listed[e.BookID]; ok {
			continue
		}

		list = append(list, e)
		listed[e.BookID] = true

	}

	sortList(list, byAvailableDate)

	if len(list) < n {
		return list
	}

	return list[:n]

}
