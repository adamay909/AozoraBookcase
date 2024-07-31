package aozorafs

import (
	"log"
	"net/http"

	"math/rand"
)

func RandomBook(lib *Library) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		B := lib.booklist[rand.Intn(len(lib.booklist))]
		log.Println("showing book", B.BookID)

		lib.randomT.Execute(w, B)
		return
	}
}
