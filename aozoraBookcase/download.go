package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func DownloadFile(path *url.URL) []byte {

	r, err := http.Get(path.String())
	if err != nil {
		log.Println("server download Remote:", err)
	}
	data, _ := io.ReadAll(r.Body)
	r.Body.Close()
	return data
}

func GetHeader(path *url.URL) map[string][]string {

	r, err := http.Head(path.String())

	if err != nil {
		log.Println(err)
		return nil
	}

	return r.Header
}
