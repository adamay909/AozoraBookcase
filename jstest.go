package main

import (
	"errors"
	"strconv"
	"syscall/js"
)

// Thanks to Google Gemini for this
func fetchData(url string) (data []byte, err error) {

	var result []js.Value

	fetchPromise := domWindow.Call("fetch", url)

	result, err = jsAwait(fetchPromise)
	if err != nil {
		return
	}

	response := result[0]
	if !response.Get("ok").Bool() {
		err = errors.New("HTTP error: " + strconv.Itoa(response.Get("status").Int()))
		return
	}

	if response.Get("status").Int() == 404 {
		err = errors.New("HTTP error: " + strconv.Itoa(response.Get("status").Int()))
		return
	}

	bytesPromise := response.Call("bytes")
	bytesResult, err := jsAwait(bytesPromise)
	if err != nil {
		return
	}

	data = bytesOf(bytesResult[0])
	return
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
