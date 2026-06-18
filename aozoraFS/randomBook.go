package aozorafs

import (
	"io/fs"
	"math"
	"math/rand"
)

func (lib *Library) GenRandomBook() (fs.File, error) {

	bk := lib.RandomBook()

	booklnk := "/books/book_" + bk.AuthorID + "_" + bk.BookID + ".html"

	return lib.genBookPage(booklnk)

}

func (lib *Library) RandomBook() *Record {

	if len(lib.authorWeight) == 0 {
		boundary := 0
		lib.authorWeight = make(map[int]string)
		for _, rec := range lib.authorsSorted {
			boundary = boundary + int(math.Round(math.Sqrt(float64(len(lib.booksByAuthor[rec.AuthorID])))*10))
			lib.authorWeight[boundary] = rec.AuthorID
			lib.totalWeight = boundary
		}
	}
	firstRound := rand.Intn(lib.totalWeight) + 1
	hit := 0
	for k := range lib.authorWeight {
		if firstRound >= k && k > hit {
			hit = k
		}
	}
	authorID := lib.authorWeight[hit]

	return lib.booksByAuthor[authorID][rand.Intn(len(lib.booksByAuthor[authorID]))]
	//	 return lib.booklist[rand.Intn(len(lib.booklist))]

}
