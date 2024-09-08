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

func (s *localStorage) Open(name string) (fs.File, error) {

	f := new(cacheFile)
	var err error
	if !fs.ValidPath(name) {

		err = &fs.PathError{Op: "open", Path: name, Err: errors.New("invalid path name")}
		return f, err
	}

	v := js.Global().Get("localStorage").Call("getItem", name)

	if v.IsNull() {

		err = &fs.PathError{Op: "open localstorage", Path: name, Err: fs.ErrNotExist}
		return f, err
	}

	f.name = name
	f.data, err = b64dec(v.String())

	return f, err
}

func (s *localStorage) Exists(name string) bool {

	for _, e := range s.Ls() {
		if e == name {
			return true
		}
	}
	return false
}

func (s *localStorage) Path() string {
	return ""
}

func (s *localStorage) RemoveAll(name string) {

	js.Global().Get("localStorage").Call("clear")

	return
}

func (s *localStorage) Stat(name string) (fs.FileInfo, error) {

	var err error

	r := new(cfinfo)

	fd, err := s.Open(name)

	f := fd.(*cacheFile)
	defer f.Close()

	if err != nil {
		return r, err
	}

	inf, err := f.Stat()

	r = inf.(*cfinfo)

	if err != nil {
		return r, err
	}

	return r, err
}

func (s *localStorage) Ls() (dir []string) {

	n := js.Global().Get("localStorage").Get("length").Int()

	for i := 0; i < n; i++ {

		dir = append(dir, js.Global().Get("localStorage").Call("key", i).String())

	}

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

func (s *localStorage) CreateFile(name string, data []byte) (fs.File, error) {

	f := new(cacheFile)

	f.name = name

	f.data = data

	err := s.Write(name, data)

	return f, err
}

func (s *localStorage) CreateEphemeral(name string, data []byte) (fs.File, error) {

	var err error

	f := new(cacheFile)

	f.name = name

	f.data = data

	return f, err
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

func (f *cacheFile) Stat() (fs.FileInfo, error) {

	info := new(cfinfo)

	info.name = f.name
	info.size = len(f.data)
	info.mode = fs.ModePerm
	info.isdir = false

	return info, errors.New("")
}

func (f *cacheFile) Close() (err error) {

	f.name = ""
	f.data = nil

	f = nil

	return

}

func (f *cacheFile) Size() int64 {

	info, _ := f.Stat()

	return info.Size()
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
