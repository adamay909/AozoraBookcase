package aozorafs

import (
	"bytes"
	"log"

	"math/rand"
)

func (lib *Library) GenRandomBook() []byte {

	w := new(bytes.Buffer)

	B := lib.booklist[rand.Intn(len(lib.booklist))]
	log.Println("showing book", B.BookID)

	lib.randomT.Execute(w, B)

	return w.Bytes()
}
