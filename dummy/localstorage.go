package main

import (
	"encoding/base64"
	"errors"
	"io/fs"
	"syscall/js"
	"time"
)

type localStorage struct{}

type cacheFile struct {
	name string
	data []byte
}

type cfinfo struct {
	name  string
	size  int
	mode  fs.FileMode
	isdir bool
	sys   struct{}
}

func (s *localStorage) Open(name string) (f *cacheFile, err error) {

	f = new(cacheFile)

	if !fs.ValidPath(name) {

		err = &fs.PathError{Op: "open", Path: name, Err: errors.New("invalid path name")}
		return
	}

	v := js.Global().Get("localStorage").Call("getItem", name)

	if v.IsNull() {

		err = &fs.PathError{Op: "open", Path: name, Err: errors.New("file not found")}
		return
	}

	f.name = name
	f.data, err = b64dec(v.String())

	return
}

func (s *localStorage) Ls() (dir []string) {

	n := js.Global().Get("localStorage").Get("length").Int()

	for i := 0; i < n; i++ {

		dir = append(dir, js.Global().Get("localStorage").Call("key", i).String())

	}

	return

}

func (f *cacheFile) Read(r []byte) (int, error) {

	var err error

	if len(r) >= len(f.data) {

		for i := range f.data {

			r[i] = f.data[i]
		}

		return len(f.data), err

	}

	for i := range r {

		r[i] = f.data[i]
	}

	return len(r), err

}

func (f *cacheFile) Stat() (info *cfinfo, err error) {

	info = new(cfinfo)

	info.name = f.name
	info.size = len(f.data)
	info.mode = fs.ModePerm
	info.isdir = false

	return
}

func (f *cacheFile) Close() (err error) {

	f = nil

	return

}

func (s *localStorage) Write(name string, data []byte) (err error) {

	if !fs.ValidPath(name) {
		err = &fs.PathError{Op: "write", Path: name, Err: errors.New("invalid path name")}
		return
	}

	js.Global().Get("localStorage").Call("setItem", name, b64enc(data))

	return
}

func (s *localStorage) WriteString(name string, data string) (err error) {

	return s.Write(name, []byte(data))

}

func (info *cfinfo) Name() string {

	return info.name
}

func (info *cfinfo) Size() int64 {

	return int64(info.size)

}

func (info *cfinfo) Mode() fs.FileMode {

	return info.mode

}

func (info *cfinfo) ModTime() time.Time {

	return time.Now()

}

func (info *cfinfo) IsDir() bool {

	return info.isdir

}

func (info *cfinfo) Sys() any {

	return info.sys

}

func b64enc(in []byte) string {

	return base64.StdEncoding.EncodeToString(in)

}

func b64dec(in string) ([]byte, error) {

	return base64.StdEncoding.DecodeString(in)

}
