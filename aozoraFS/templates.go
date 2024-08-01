package aozorafs

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	_ "embed" //for embedding resources
)

//use this to hard code templates

//go:embed resources/defaultcss.css
var fileDefaultcss string

//go:embed resources/random.html
var randombookhtml string

func (lib *Library) saveCSS() {

	f, err := os.Create(filepath.Join(lib.cache, "ebooks.css"))
	if err != nil {
		log.Println("1", err)
	}
	defer f.Close()
	log.Println("css has", len([]byte(fileDefaultcss)))
	n, err := f.Write([]byte(fileDefaultcss))
	if err != nil {
		log.Println("2", err)
	}
	log.Println("wrote", n, "bytes")
	return
}

//go:embed resources/index.html
var fileIndexhtml string

func (lib *Library) mainIndexTemplate() {

	lib.indexT = template.Must(template.New("index.html").Parse(fileIndexhtml))
}

//go:embed resources/recent.html
var fileRecenthtml string

func (lib *Library) recentTemplate() {

	lib.recentT = template.Must(template.New("recent.html").Parse(fileRecenthtml))
}

//go:embed resources/author.html
var fileAuthorhtml string

func (lib *Library) authorpageTemplate() {

	lib.authorT = template.Must(template.New("author.html").Parse(fileAuthorhtml))
}

//go:embed resources/book.html
var fileBookhtml string

func (lib *Library) bookpageTemplate() {

	NdcOf := func(i string) string {

		return ndcmap()[i]

	}

	NdcPOf := func(i [2]string) string {

		return ndcmap()[i[0]]

	}

	NdcCOf := func(i [2]string) string {

		return ndcmap()[i[1]]

	}

	funcMap := template.FuncMap{"ndc1": NdcPOf,
		"ndc2": NdcCOf, "ndc": NdcOf}

	lib.bookT = template.Must(template.New("book.html").Funcs(funcMap).Parse(fileBookhtml))
}

//go:embed resources/category.html
var fileCategoryhtml string

func (lib *Library) categorypageTemplate() {

	lib.categoryT = template.Must(template.New("book.html").Parse(fileCategoryhtml))
}

//go:embed resources/random.html
var randomBookhtml string

func (lib *Library) randomBookTemplate() {

	NdcOf := func(i string) string {

		return ndcmap()[i]

	}

	NdcPOf := func(i [2]string) string {

		return ndcmap()[i[0]]

	}

	NdcCOf := func(i [2]string) string {

		return ndcmap()[i[1]]

	}

	funcMap := template.FuncMap{"ndc1": NdcPOf,
		"ndc2": NdcCOf, "ndc": NdcOf}

	lib.randomT = template.Must(template.New("book.html").Funcs(funcMap).Parse(randomBookhtml))
}
