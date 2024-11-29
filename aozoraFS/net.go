package aozorafs

var download func(path string) []byte

func SetDownloader(f func(path string) []byte) {

	download = f

}
