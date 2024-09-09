package main

import (
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"syscall/js"
)

var (
	uint8ArrayConstructor,
	fileConstructor,
	arrayConstructor,
	blobConstructor,
	domWindow,
	domDocument,
	domHTML,
	domBody js.Value
)

func init() {

	uint8ArrayConstructor = js.Global().Get("Uint8Array")

	fileConstructor = js.Global().Get("File")

	arrayConstructor = js.Global().Get("Array")

	blobConstructor = js.Global().Get("Blob")

	domWindow = js.Global()

	domDocument = domWindow.Get("document")

	domHTML = domDocument.Get("documentElement")

	domBody = domDocument.Get("body")
}

func uint8arrayOf(data []byte) js.Value {

	jsdata := uint8ArrayConstructor.New(len(data))

	js.CopyBytesToJS(jsdata, data)

	return jsdata

}

func saveFile(file js.Value) {

	href := domWindow.Get("URL").Call("createObjectURL", file)

	link := domDocument.Call("createElement", "a")

	link.Set("href", href)

	link.Set("download", file.Get("name"))

	domBody.Call("appendChild", link)

	link.Call("click")

	domBody.Call("removeChild", link)

	domWindow.Get("URL").Call("revokeObjectURL", href)

	return

}

// Thanks to https://javascript.plainenglish.io/javascript-create-file-c36f8bccb3be for how to do this
func createJSFile(data []byte, name string) js.Value {

	jsdata := uint8arrayOf(data)

	blob := blobConstructor.New(arrayConstructor.New(jsdata), map[string]any{"type": mime.TypeByExtension(filepath.Ext(name))})

	return fileConstructor.New(arrayConstructor.New(blob), filepath.Base(name))

}

func getUrl() string {

	return domWindow.Get("location").Get("href").String()

}

func getHost() string {

	return domWindow.Get("location").Get("host").String()

}

func getHash() string {

	return domWindow.Get("location").Get("hash").String()

}

func setHash(h string) {

	domWindow.Get("location").Set("hash", h)

	return

}

func replaceBody(p string) {

	domBody.Set("innerHTML", p)

	return

}

func fetchData(path *url.URL) (data []byte) {

	loc := path.String()

	log.Println("fetching", loc)

	r, err := http.Get(loc)

	if err != nil {
		log.Println(err)
		return
	}

	data, _ = io.ReadAll(r.Body)

	return data

}

func getElementById(id string) (elem js.Value, err error) {

	elem = domDocument.Call("getElementById", id)

	if elem.IsNull() {
		err = errors.New("no element found with ID " + id)
		log.Println(err)
	}

	return
}

func scrollTo(elem js.Value) {

	elem.Call("scrollIntoView")

	return

}

// addEventListener adds an event listener. The function f is called
// with the event object as its first argument and the arguments given by params.
func addEventListener(elem js.Value, eventType string, f func(event js.Value, args ...any), params ...any) {

	elem.Call("addEventListener", eventType, js.FuncOf(func(this js.Value, margs []js.Value) any {

		f(margs[0], params...)
		return true
	}), true)

	return
}
