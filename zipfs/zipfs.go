package zipfs

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"log"
	"time"
)

type Ziparchive struct {
	z    *zip.Reader
	name string
}

type Zipfile struct {
	data []byte
	name string
	stat zinfo
	sys  *zip.File
}

type zinfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
}

func OpenZipArchive(fsys fs.FS, name string) (*Ziparchive, error) {

	log.Println("attemption to open:", name)

	za := new(Ziparchive)

	fd, err := fsys.Open(name)

	if err != nil {
		log.Println("zipfs: error 1")
		return za, err
	}

	info, _ := fd.Stat()

	data := make([]byte, info.Size())

	fd.Read(data)

	log.Println("len of underlying data:", len(data))

	r := bytes.NewReader(data)

	za.z, err = zip.NewReader(r, int64(len(data)))

	return za, err

}

func (za *Ziparchive) CloseArchive() {

	za.z = nil

	za = nil
}

func (za *Ziparchive) Open(name string) (fs.File, error) {

	f := new(Zipfile)
	found := false
	index := 0

	for i, f := range za.z.File {
		if f.FileHeader.Name == name {
			found = true
			index = i
			break
		}
	}

	if !found {
		err := fs.ErrNotExist
		return f, err
	}

	zf := za.z.File[index]
	rc, err := zf.Open()
	defer rc.Close()
	if err != nil {
		return f, err
	}

	//size := zf.FileHeader.UncompressedSize64

	// f.data = make([]byte, size)
	// rc.Read(f.data)
	f.data, _ = io.ReadAll(rc)
	f.name = name
	f.sys = zf

	f.stat.name = name
	f.stat.size = int64(len(f.data))
	f.stat.mode = fs.ModePerm
	f.stat.modTime = zf.FileHeader.Modified
	f.stat.isDir = false

	return f, nil

}

func (z *Ziparchive) Read(name string) ([]byte, error) {

	var b []byte

	fz, err := z.Open(name)

	if err != nil {
		return b, err
	}

	f := fz.(*Zipfile)

	return f.data, err
}

func (z *Ziparchive) ReadMust(name string) []byte {

	d, _ := z.Read(name)

	return d

}
func (f *Zipfile) size() int64 {

	return int64(len(f.data))
}

func (f *Zipfile) Read(r []byte) (n int, err error) {

	log.Println("attempting to read", len(f.data), "bytes")

	if len(r) >= len(f.data) {
		log.Println("r is long enough")
		i := 0
		for i = range f.data {
			r[i] = f.data[i]
		}
		log.Println("copied", i, "bytes")

		return len(f.data), err
	}

	for i := range r {
		r[i] = f.data[i]
	}

	return len(r), err
}

func (f *Zipfile) Stat() (fs.FileInfo, error) {

	return &f.stat, nil

}

func (f *Zipfile) Close() error {

	f = new(Zipfile)

	return nil
}

func (i *zinfo) Name() string {
	return i.name
}

func (i *zinfo) Size() int64 {
	return i.size
}

func (i *zinfo) Mode() fs.FileMode {
	return i.mode
}

func (i *zinfo) ModTime() time.Time {
	return i.modTime
}

func (i *zinfo) IsDir() bool {
	return false
}

func (i *zinfo) Sys() any {
	return nil
}
