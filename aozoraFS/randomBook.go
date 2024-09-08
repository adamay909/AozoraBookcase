package aozorafs

import (
	"bytes"
	"log"

	"math/rand"
)

func (lib *Library) GenRandomBook() []byte {

	var D struct {
		B *Record
	}

	D.B = lib.RandomBook()

	w := new(bytes.Buffer)

	err := lib.randomT.Execute(w, D)

	if err != nil {
		log.Println(err)
	}

	return w.Bytes()
}

func (lib *Library) RandomBook() *Record {

	return lib.booklist[rand.Intn(len(lib.booklist))]

}
