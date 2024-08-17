package aozorafs

import (
	"bytes"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	_ "embed" //for embedding resources
)

var FileDefaultcss string

var Randombookhtml string

var FileIndexhtml string

var FileRecenthtml string

var FileAuthorhtml string

var FileBookhtml string

var FileCategoryhtml string

var RandomBookhtml string

var FileSearchhtml string

var FileSearchresultshtml string

func (lib *Library) ImportTemplates(dir fs.ReadDirFS) {

	entry, err := dir.ReadDir(".")

	if len(entry) != 1 {
		log.Println("templates must be inside a single subdirectory")
		return
	}

	if err != nil {
		log.Println(err)
		return
	}

	dirname := entry[0].Name()
	entry, err = dir.ReadDir(dirname)
	if err != nil {
		log.Println("templates must be inside a single subdirectory")
		return
	}

	for k := range entry {

		log.Println("checking", entry[k].Name())

		f, err := dir.Open(filepath.Join(dirname, entry[k].Name()))
		defer f.Close()

		if err != nil {
			log.Println(err)
			return
		}

		info, err := f.Stat()

		if err != nil {
			log.Println(err)
			return
		}

		data := make([]byte, info.Size())
		f.Read(data)
		t := strings.Split(string(data), "\n")[0]

		tn, err := templateName(t)

		if err != nil {
			log.Println("not a template file:", info.Name())
			continue
		}

		NdcOf := func(i string) string {

			return ndcmap()[i]

		}

		NdcPOf := func(i [2]string) string {

			return ndcmap()[i[0][:1]]

		}

		NdcCOf := func(i [2]string) string {

			return ndcmap()[i[1]]

		}

		funcMap := template.FuncMap{"ndc1": NdcPOf,
			"ndc2": NdcCOf, "ndc": NdcOf}

		//Now define the templates

		log.Println("template name is ", tn)
		switch tn {

		case "defaultcss":
			buf := new(bytes.Buffer)

			err := template.Must(template.New("css").Parse(string(data))).Execute(buf, "")
			if err != nil {
				log.Println(err)
			}
			_, err = lib.cache.CreateFile("ebooks.css", buf.Bytes())
			if err != nil {
				log.Println(err)
			}
			log.Println("saved css")

		case "randombook":
			lib.randomT = template.Must(template.New("random.html").Funcs(funcMap).Parse(string(data)))

		case "index":
			lib.indexT = template.Must(template.New("index.html").Parse(string(data)))

		case "recent":
			lib.recentT = template.Must(template.New("recent.html").Parse(string(data)))

		case "author":
			lib.authorT = template.Must(template.New("author.html").Parse(string(data)))

		case "book":
			lib.bookT = template.Must(template.New("book.html").Funcs(funcMap).Parse(string(data)))

		case "category":
			lib.categoryT = template.Must(template.New("category.html").Parse(string(data)))

		case "search":

		case "searchresult":
			lib.searchresultT = template.Must(template.New("searchresult.html").Parse(string(data)))

		default:

		}

	}
}

func templateName(src string) (name string, err error) {

	src = strings.TrimSpace(src)

	name = strings.TrimPrefix(src, `{{/*`)

	if name == src {
		return "", errors.New("not a template file")
	}

	name = strings.TrimSuffix(name, `*/}}`)

	if name == src {
		return "", errors.New("not a template file")
	}

	return strings.TrimSpace(name), err
}
