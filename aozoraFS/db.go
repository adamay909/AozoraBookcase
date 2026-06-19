package aozorafs

import (
	"log"
	"path"
	"sort"
	"strings"

	str "github.com/adamay909/AozoraBookcase/aozoraFS/stringops"
	"github.com/adamay909/AozoraBookcase/aozoraFS/zipfs"
)

func (lib *Library) ConstructLibrary(data []byte) {
	za, _ := zipfs.ZipArchiveFromData(data)
	defer za.CloseArchive()
	lib.GetBooklist(za.ReadMust("list_person_all_extended_utf8.csv"))
	lib.setupAuthorsList()
	return
}

func (lib *Library) FetchLibraryData() []byte {

	log.Println("getting library catalog information")

	pathStr := path.Join(lib.src, "/index_pages", "list_person_all_extended_utf8.zip")

	log.Println("requesting db", pathStr)

	data, _ := download(pathStr)

	return data
}

func (lib *Library) setupAuthorsList() {
	lib.authorsSorted = make([]*Record, 0, len(lib.booksByAuthor))

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

func (lib *Library) GetBooklist(d []byte) {

	rows := str.NewLineReader(string(d))

	t := len(strings.Split(string(d), "\n"))

	lib.booksByID = make(map[string][]*Record, t)
	lib.booksByAuthor = make(map[string][]*Record, t)

	row, eof := rows.Read()

	headings := getCells(removeBOM(row))

	hmap := make(map[string]int)

	for i := range headings {

		hmap[headings[i]] = i

	}

	col := func(h string) int {

		return hmap[h]

	}

	cells := make([]string, len(headings))

	var uri string
	var book *Record

	//read into records
	for row, eof = rows.Read(); !eof; row, eof = rows.Read() {

		cells = getCells(row)

		if lib.strict {
			if cells[col("作品著作権フラグ")] == "あり" || cells[col("人物著作権フラグ")] == "あり" {
				continue
			}
		}

		//uri = aozoraPath(cells[col("XHTML/HTMLファイルURL")])
		uri = aozoraPath(cells[col("テキストファイルURL")])

		if uri == "" {
			continue
		}

		book = new(Record)

		book.URI = path.Join(lib.src, uri)
		book.NDC = cells[col("分類番号")]
		book.setCategory(lib.Categories)

		book.KanaZukai = cells[col("文字遣い種別")]

		if lib.kids {
			if !book.isChildrensBook() {
				continue
			}
		}

		book.BookID = cells[col("作品ID")]
		book.Title = cells[col("作品名")]
		book.TitleY = cells[col("作品名読み")]
		book.TitleSort = cells[col("ソート用読み")]
		book.Subtitle = cells[col("副題")]
		book.SubtitleY = cells[col("副題読み")]
		//book.OriginalTitle = cells[col("原題")]
		book.PublDate = cells[col("初出")]
		book.FirstAvailable = cells[col("公開日")]
		//book.ModTime = cells[col("最終更新日")]
		book.AuthorID = cells[col("人物ID")]
		book.NameSei = cells[col("姓")]
		book.NameMei = cells[col("名")]
		book.NameSeiY = cells[col("姓読み")]
		book.NameMeiY = cells[col("名読み")]
		book.NameSeiSort = cells[col("姓読みソート用")]
		book.NameMeiSort = cells[col("名読みソート用")]
		//book.NameSeiR = cells[col("姓ローマ字")]
		//book.NameMeiR = cells[col("名ローマ字")]
		book.Role = cells[col("役割フラグ")]
		book.DoBirth = cells[col("生年月日")]
		book.DoDeath = cells[col("没年月日")]

		lib.booklist = append(lib.booklist, book)
		lib.booksByID[book.BookID] = append(lib.booksByID[book.BookID], book)
		lib.booksByAuthor[book.AuthorID] = append(lib.booksByAuthor[book.AuthorID], book)

	}

	log.Println("finished parsing db.")

	return
}

func (lib *Library) getRecents(n int) []*Record {

	if n*100+100 > len(lib.booksByDate) {
		return lib.booksByDate[n*100:]
	} else {
		return (lib.booksByDate[n*100 : n*100+100])
	}
}

func (lib *Library) SortByAvailDate() {

	if len(lib.booksByDate) != 0 {
		return
	}

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

	lib.booksByDate = append(lib.booksByDate, list...)

}

func (lib *Library) LenDistinctBooks() int {

	lib.SortByAvailDate()

	return len(lib.booksByDate)

}

func aozoraPath(fullpath string) string {

	path := strings.TrimPrefix(fullpath, "https://www.aozora.gr.jp")

	if len(path) == len(fullpath) {
		return ""
	}

	return path
}

func (lib *Library) LatestReads() []string {
	return lib.latestReads
}

func (lib *Library) SetLatestReads(list []string) {
	lib.latestReads = nil
	lib.latestReads = append(lib.latestReads, list...)
}

func (lib *Library) SetFavorites(list []string) {
	lib.favorites = make(map[string]struct{})
	for _, item := range list {
		lib.favorites[item] = struct{}{}
	}
}

func (lib *Library) AddFavorite(id string) {
	lib.favorites[id] = struct{}{}
}

func (lib *Library) RemoveFavorite(id string) {
	delete(lib.favorites, id)
}

func (lib *Library) FavoriteList() []string {
	list := make([]string, len(lib.favorites))
	for id, _ := range lib.favorites {
		list = append(list, id)
	}
	return list
}

func (lib *Library) IsFavorite(id string) bool {
	_, ok := lib.favorites[id]
	return ok
}
