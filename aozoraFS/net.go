package aozorafs

import "net/url"

var download func(path *url.URL) []byte

func SetDownloader(f func(path *url.URL) []byte) {

	download = f

}
