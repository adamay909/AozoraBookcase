package aozorafs

import "net/url"

var download func(path *url.URL) []byte

func SetDownloader(f func(path *url.URL) []byte) {

	download = f

}

var getHeader func(path *url.URL) map[string][]string

func SetHeader(f func(*url.URL) map[string][]string) {

	getHeader = f

}

func get(header map[string][]string, key string) string {

	r, ok := header[key]

	if !ok {
		return ""
	}

	return r[0]
}
