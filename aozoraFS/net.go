package aozorafs

var download func(path string) ([]byte, error)

func SetDownloader(f func(path string) ([]byte, error)) {

	download = f

}
