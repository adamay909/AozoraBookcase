package aozorafs

import (
	"io/fs"
	"math/rand"
)

func (lib *Library) GenRandomBook() (fs.File, error) {

	bk := lib.RandomBook()

	booklnk := "/books/book_" + bk.AuthorID + "_" + bk.BookID + ".html"

	return lib.genBookPage(booklnk)

}

func (lib *Library) RandomBook() *Record {

	return lib.booklist[rand.Intn(len(lib.booklist))]

}
