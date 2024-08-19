package server

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func DownloadFile(path *url.URL) []byte {

	r, err := http.Get(path.String())
	log.Println("server: success downloading")
	if err != nil {
		log.Println("server download Remote:", err)
	}
	data, _ := io.ReadAll(r.Body)
	r.Body.Close()
	return data
}
