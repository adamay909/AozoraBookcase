package aozorafs

import (
	"bytes"
	"errors"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

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

			return lib.Categories[i]

		}

		NdcDOf := func(i [3]string) string {

			return lib.Categories[i[0]]

		}

		NdcCOf := func(i [3]string) string {

			return lib.Categories[i[1]]

		}

		NdcSOf := func(i [3]string) string {

			return lib.Categories[i[2]]

		}

		NdcNDOf := func(i [3]string) string {

			return i[0]

		}

		NdcNCOf := func(i [3]string) string {

			return i[1]

		}

		NdcNSOf := func(i [3]string) string {

			return i[2]

		}

		NdcMust := func(i string) string {

			if len(i) == 1 {

				return lib.Categories[i]
			}
			if i[:1] == "9" && len(i) == 3 {
				return lib.Categories[i[:1]] + " : " + lib.Categories[i[:2]] + " : " + lib.Categories[i[:3]]
			} else {
				return lib.Categories[i[:1]] + " : " + lib.Categories[i[:2]]
			}

		}

		funcMap := template.FuncMap{"ndc1": NdcDOf,
			"ndc2": NdcCOf, "ndc3": NdcSOf, "ndc": NdcOf, "ndcm": NdcMust,
			"ndcn2": NdcNCOf, "ndcn3": NdcNSOf, "ndcn1": NdcNDOf}

		//Now define the templates

		switch tn {

		case "defaultcss":
			buf := new(bytes.Buffer)

			err := template.Must(template.New("css").Parse(string(data))).Execute(buf, "")
			if err != nil {
				log.Println(err)
			}
			_, err = lib.cache.CreateEphemeral("ebooks.css", buf.Bytes())
			if err != nil {
				log.Println(err)
			}

		case "readingpanecss":
			buf := new(bytes.Buffer)

			err := template.Must(template.New("css").Parse(string(data))).Execute(buf, "")
			if err != nil {
				log.Println(err)
			}
			_, err = lib.cache.CreateEphemeral("readingpane.css", buf.Bytes())
			if err != nil {
				log.Println(err)
			}

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

		case "reading":
			lib.readingT = template.Must(template.New("reading.html").Parse(string(data)))

		case "search":

		case "searchresult":
			lib.searchresultT = template.Must(template.New("searchresult.html").Funcs(funcMap).Parse(string(data)))

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
