package aozorafs

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
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

	os.RemoveAll(filepath.Join(lib.cache, "index.html"))
	os.RemoveAll(filepath.Join(lib.cache, "authors/"))
	os.RemoveAll(filepath.Join(lib.cache, "books/"))

	if clean {
		os.Remove(filepath.Join(lib.root, "aozoradata.zip"))
	}

	os.Mkdir(lib.cache, 0766)

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
func (lib *Library) LoadBooklist() {

	fi, err := os.Stat(filepath.Join(lib.root, "aozoradata.zip"))

	switch {
	case os.IsNotExist(err):
		lib.UpdateDB()

	case err == nil:
		if lib.UpstreamUpdated(fi.ModTime()) {
			lib.UpdateDB()
		}

	default:
		log.Println(err)
		return
	}

	lib.UpdateBooklist()

	lib.updatePages()

	go lib.RefreshBooklist()

}

/*UpstreamUpdated reports whether the upstream database has been updated since it was last updated locally.
 */
func (lib *Library) UpstreamUpdated(t time.Time) bool {

	r, err := http.Head(lib.src + filepath.Join("/index_pages", "list_person_all_extended_utf8.zip"))

	if err != nil {
		log.Println(err)
		return false
	}

	m, err := time.Parse(time.RFC1123, r.Header.Get("Last-Modified"))

	log.Println("Server reports last update time of: ", m)

	return m.After(t)
}

/*UpdateDB downloads the database from upstream.*/
func (lib *Library) UpdateDB() {

	path, err := url.Parse(lib.src + filepath.Join("/index_pages", "list_person_all_extended_utf8.zip"))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("requesting db", path.String())

	data := downloadFile(path)

	err = os.WriteFile(filepath.Join(lib.root, "aozoradata.zip"), data, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("data saved to ", filepath.Join(lib.root, "aozoradata.zip"))
}

/*UpdateBooklist updates the booklist of lib from the locally available database.*/
func (lib *Library) UpdateBooklist() {

	lib.updating = true
	zf, _ := os.ReadFile(filepath.Join(lib.root, "aozoradata.zip"))
	lib.getBooklist(unzip(zf))
	lib.consolidateBookRecords()
	sortList(lib.booklist, byAuthor)
	lib.Categories = ndcmap()
	os.RemoveAll(lib.cache)
	lib.lastUpdated = time.Now()
	log.Println("sorted entries.")
	lib.updating = false
	return
}

func (lib *Library) consolidateBookRecords() {

	sortList(lib.booklist, byBookID)

	for i, j := 0, 0; i < len(lib.booklist); {

		var list []*Record

		e := lib.booklist[i]

		for j = i; j < len(lib.booklist) && lib.booklist[j].BookID == e.BookID; j++ {

			list = append(list, lib.booklist[j])

		}

		for _, l := range list {

			list[0].Contributors = append(list[0].Contributors, ContribRole{l.Role, l.AuthorID, l})
		}
		sort.Slice(list[0].Contributors, byRole(list[0].Contributors))

		for _, e := range list {

			if e == list[0] {
				continue
			}

			e.Contributors = nil

			e.Contributors = append(e.Contributors, list[0].Contributors...)

		}
		i = j
	}

	return

}

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

func (lib *Library) getBooklist(d []byte) {

	rows := strings.Split(string(d), "\n")
	log.Println("database has ", len(rows), "entries")
	headings, err := csv.NewReader(strings.NewReader(rows[0])).Read()
	if err != nil {
		log.Println("error reading Aozora Bunko database", err)
		return
	}

	col := make(map[string]int)

	//get column number for each heading
	for i, h := range headings {
		col[h] = i
	}

	//read into records
	for i := 1; i < len(rows)-1; i++ {

		r := rows[i]
		cells, _ := csv.NewReader(strings.NewReader(r)).Read()

		if len(cells) == 0 {
			break
		}

		book := new(Record)

		book.BookID = cells[col["作品ID"]]
		book.Title = cells[col["作品名"]]
		book.TitleY = cells[col["作品名読み"]]
		book.TitleSort = cells[col["ソート用読み"]]
		book.Subtitle = cells[col["副題"]]
		book.SubtitleY = cells[col["副題読み"]]
		book.OriginalTitle = cells[col["原題"]]
		book.PublDate = cells[col["初出"]]
		book.NDC = cells[col["分類番号"]]
		book.KanaZukai = cells[col["文字遣い種別"]]
		book.WorkCopyright = cells[col["作品著作権フラグ"]]
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
		book.AuthorCopyright = cells[col["人物著作権フラグ"]]
		book.URI = strings.TrimPrefix(cells[col["XHTML/HTMLファイルURL"]], "https://www.aozora.gr.jp")
		book.setCategory()

		if lib.strict {
			if book.WorkCopyright == "あり" || book.AuthorCopyright == "あり" {
				continue
			}
		}

		if lib.kids {
			if !book.isChildrensBook() {
				continue
			}
		}
		lib.booklist = append(lib.booklist, book)

	}
	log.Println("finished parsing db.")
	return
}

func (lib *Library) updatePages() {

	os.RemoveAll(lib.cache + "/authors")
	os.RemoveAll(lib.cache + "/books")
	os.Mkdir(lib.cache, 0766)
	os.Mkdir(lib.cache+"/authors", 0766)
	os.Mkdir(lib.cache+"/books", 0766)
	os.Mkdir(lib.cache+"/files", 0766)
	os.Mkdir(lib.cache+"/categories", 0766)
	os.RemoveAll(lib.cache + "index.html")
	os.RemoveAll(lib.cache + "recent.html")
	lib.saveCSS()
	//lib.genMainIndex()

	lib.lastUpdated = time.Now()
}

/*ReadBooklist is for retrieving the list of books stored in lib.*/
func (lib *Library) ReadBooklist() (o []*Record) {

	o = append(o, lib.booklist...)

	return

}

/*
WriteBooklist is for adding l as the booklist of lib.

ReadBooklist and WriteBooklist are provided for manual inspection and editing of the booklist.
*/
func (lib *Library) WriteBooklist(l []*Record) {

	lib.booklist = nil

	lib.booklist = append(lib.booklist, l...)

	return

}

func (lib *Library) getRecents(n int) []*Record {

	var list []*Record

	list = append(list, lib.booklist...)

	sortList(list, byAvailableDate)

	return list[:n]

}
