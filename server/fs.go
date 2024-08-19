package server

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type diskFS struct {
	path string
	sys  fs.FS
}

func NewDiskFS(path string) (s *diskFS) {

	s = new(diskFS)

	s.sys = os.DirFS(path)

	s.path = path

	return

}

func (s *diskFS) Open(name string) (fs.File, error) {

	return s.sys.Open(name)

}

func (s *diskFS) Stat(name string) (info fs.FileInfo, err error) {

	f, err := s.sys.Open(name)

	if err != nil {
		return
	}

	return f.Stat()

}

func (s *diskFS) CreateFile(name string, data []byte) (fs.File, error) {

	dir := filepath.Dir(name)

	err := os.MkdirAll(filepath.Join(s.path, dir), fs.ModePerm)

	if err != nil {
		log.Println(err)
	}

	err = os.WriteFile(filepath.Join(s.path, name), data, fs.ModePerm)

	if err != nil {
		log.Println(err)
	}

	return s.Open(name)

}

func (s *diskFS) CreateEphemeral(name string, data []byte) (fs.File, error) {

	return s.CreateFile(name, data)

}
