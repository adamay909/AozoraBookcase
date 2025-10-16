package main

import (
	"errors"
	"strconv"
	"syscall/js"
)

// Thanks to Google Gemini for this
func fetchData(url string) []byte {

	fetchPromise := domWindow.Call("fetch", url)

	result, err := jsAwait(fetchPromise)
	if err != nil {
		panic("file download failed")
	}

	response := result[0]
	if !response.Get("ok").Bool() {
		panic("HTTP error: " + strconv.Itoa(response.Get("status").Int()))
	}

	bytesPromise := response.Call("bytes")
	bytesResult, err := jsAwait(bytesPromise)
	if err != nil {
		panic("file download filed")
	}

	data := make([]byte, bytesResult[0].Get("length").Int())
	js.CopyBytesToGo(data, bytesResult[0])
	return data
}

// await the resolution of promise. Err is non-nil when promise
// is not fulfilled.
func jsAwait(promise js.Value) (resp []js.Value, err error) {

	fail := false

	ch := make(chan []js.Value)

	var resolve, reject js.Func

	resolve = js.FuncOf(func(this js.Value, args []js.Value) any {
		ch <- args
		resolve.Release()
		reject.Release()
		return nil
	})

	reject = js.FuncOf(func(this js.Value, args []js.Value) any {
		ch <- args
		fail = true
		resolve.Release()
		reject.Release()
		return nil
	})

	promise.Call("then", resolve, reject)
	resp = <-ch

	if fail {
		err = errors.New("promise rejected")
	}

	return
}
