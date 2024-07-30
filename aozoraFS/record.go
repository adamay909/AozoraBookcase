package aozorafs

import (
	"path/filepath"
	"strings"

	"github.com/adamay909/AozoraConvert/jptools"
	"github.com/adamay909/AozoraConvert/runes"
)

func hasSameAuthor(a, b *Record) bool {

	return a.AuthorID == b.AuthorID

}

func (lib *Library) getBookRecord(authorID, bookID string) (*Record, int) {

	for i, e := range lib.booklist {

		if e.BookID == bookID {
			if e.AuthorID == authorID {
				return e, i
			}
		}
	}

	return lib.booklist[0], 0

}

func (lib *Library) getBookRecordSimple(bookID string) *Record {

	for _, e := range lib.booklist {

		if e.BookID == bookID {
			return e
		}
	}

	return lib.booklist[0]

}

/*PrevBook returns the previous book by the same author.*/
func (lib *Library) PrevBook(b *Record) *Record {

	list := lib.getBooksByAuthor(b.AuthorID)

	k := 0

	for i, e := range list {
		if e.BookID == b.BookID {
			k = i - 1
			break
		}
	}

	if k == -1 {
		k = len(list) - 1
	}

	return list[k]
}

/*NextBook returns the next book by the same author.*/
func (lib *Library) NextBook(b *Record) *Record {

	list := lib.getBooksByAuthor(b.AuthorID)

	k := 0

	for i, e := range list {
		if e.BookID == b.BookID {
			k = i + 1
			break
		}
	}

	if k == len(list) {
		k = 0
	}
	return list[k]
}

func (lib *Library) allAuthors() (list []*Record) {

	list = append(list, lib.booklist[0])

	author0 := lib.booklist[0].AuthorID

	for _, e := range lib.booklist {

		if e.AuthorID != author0 {

			list = append(list, e)

			author0 = e.AuthorID
		}
	}
	return list
}

/*NextAuthor returns the next author where the current author is the author of b. */
func (lib *Library) NextAuthor(b *Record) *Record {

	list := lib.allAuthors()

	for i, e := range list {

		if e.AuthorID == b.AuthorID {

			if i == len(list)-1 {
				return list[0]
			}
			return list[i+1]
		}
	}

	return list[0]
}

/*PrevAuthor returns the previous author where the current author is the author of b. */
func (lib *Library) PrevAuthor(b *Record) *Record {

	list := lib.allAuthors()

	for i, e := range list {

		if e.AuthorID == b.AuthorID {

			if i == 0 {
				return list[len(list)-1]
			}
			return list[i-1]
		}
	}

	return list[0]
}

func (lib *Library) getAuthorsByInitial(s string) (authorList []*Record) {

	var tlist []*Record

	for _, b := range lib.booklist {

		if strings.HasPrefix(b.NameSeiSort, s) {
			tlist = append(tlist, b)
		}
	}

	if len(tlist) == 0 {
		return
	}

	authorList = append(authorList, tlist[0])

	for i := 1; i < len(tlist); i++ {

		if !hasSameAuthor(tlist[i], tlist[i-1]) {
			authorList = append(authorList, tlist[i])
		}
	}
	return
}

func (lib *Library) getBooksByAuthor(aID string) (list []*Record) {

	for _, e := range lib.booklist {

		if e.AuthorID == aID {
			list = append(list, e)
		}
	}
	sortList(list, byTitle)
	return
}

func isKana(s string) bool {
	if len(s) == 0 {
		return false
	}
	t := runes.Runes(s)[0]
	return jptools.CharType(t) == jptools.Hiragana
}

func containsKanji(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range runes.Runes(s) {
		if jptools.CharType(c)&jptools.Kanji != 0 {
			return true
		}
	}

	return false
}

// FullName returns the full name in Japanese
// order (family name first), with a space between
// family and given name.
func (b *Record) FullName() string {
	switch {
	case len(b.NameSei) == 0:
		return b.NameMei
	case len(b.NameMei) == 0:
		return b.NameSei
	default:
		f := b.NameSei + " " + b.NameMei
		if !containsKanji(f) {
			f = b.NameSei + "・" + b.NameMei
		}
		return f
	}
}

// FullNameY returns the yomi of full name.
// It return a blank when there is no kanji
// in the name.
func (b *Record) FullNameY() string {
	if !containsKanji(b.FullName()) {
		return ""
	}
	switch {
	case len(b.NameSeiY) == 0:
		return b.NameMeiY
	case len(b.NameMeiY) == 0:
		return b.NameSeiY
	default:
		return b.NameSeiY + " " + b.NameMeiY
	}
}

func (b *Record) fullNameS() string {
	switch {
	case len(b.NameSeiSort) == 0:
		return b.NameMeiSort
	case len(b.NameMeiSort) == 0:
		return b.NameSeiSort
	default:
		return b.NameSeiSort + " " + b.NameMeiSort
	}
}

func (b *Record) fullNameMeta() string {
	switch {
	case len(b.NameSei) == 0:
		return b.NameMei
	case len(b.NameMei) == 0:
		return b.NameSei
	default:
		return b.NameSei + ", " + b.NameMei
	}
}

// FileName returns the filename of the file associated with b.
func (b *Record) FileName() string {
	return strings.TrimSuffix(filepath.Base(b.URI), filepath.Ext(b.URI))
}

func (b *Record) setCategory() {

	ndc := ndcmap()

	s := strings.TrimLeft(b.NDC, "NDC ")
	if strings.HasPrefix(s, "K") {
		b.Kids = true
		s = strings.TrimLeft(s, "K ")
	}
	codes := strings.Split(s, " ")
	for _, c := range codes {
		if len(c) != 3 {
			continue
		}
		_, ok := ndc[c[:1]]
		if ok {
			_, ok := ndc[c[:2]]
			if ok {
				b.Categories = append(b.Categories, [2]string{c[:1], c[:2]})
			}
		}
	}
	for _, e := range b.Categories {

		b.Category = b.Category + ndc[e[0]] + "--" + ndc[e[1]] + ";"

	}
	b.Category = strings.TrimSuffix(b.Category, ";")
}

func (b *Record) isChildrensBook() bool {

	if b.Kids {
		return b.KanaZukai == "新字新仮名"
	}

	return false
}

// RealBookID returns the book id without any suffixes.
func (b *Record) RealBookID() string {
	return strings.TrimRight(b.BookID, "a")
}

// Dates returns the dates of the author of b.
func (b *Record) Dates() string {

	dob := strings.Split(b.DoBirth, `-`)[0]
	dod := strings.Split(b.DoDeath, `-`)[0]

	return dob + `-` + dod

}
