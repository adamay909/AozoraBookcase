package aozorafs

import (
	"html/template"
	"os"
	"path/filepath"
)

func (lib *Library) saveCSS() {
	_, err := os.Stat(filepath.Join(lib.resources, "ebooks.css"))
	if os.IsNotExist(err) {
		os.WriteFile(filepath.Join(lib.resources, "ebooks.css"), []byte(defaultCSS(lib)), 0644)
	}
	return
}

func (lib *Library) saveSimpleCSS() {
	_, err := os.Stat(filepath.Join(lib.resources, "simple.css"))
	if os.IsNotExist(err) {
		os.WriteFile(filepath.Join(lib.resources, "ebooks.css"), []byte(simpleCSS(lib)), 0644)
	}
	return
}

func (lib *Library) mainIndexTemplate() {

	p, err := os.ReadFile(filepath.Join(lib.resources, "index.html"))
	if os.IsNotExist(err) {
		os.WriteFile(filepath.Join(lib.resources, "index.html"), []byte(indexTemplate(lib)), 0644)
		p, err = os.ReadFile(filepath.Join(lib.resources, "index.html"))
	}

	lib.indexT = template.Must(template.New("index.html").Parse(string(p)))
}

func (lib *Library) recentTemplate() {

	p, err := os.ReadFile(filepath.Join(lib.resources, "recent.html"))
	if os.IsNotExist(err) {
		os.WriteFile(filepath.Join(lib.resources, "recent.html"), []byte(recentTemplate(lib)), 0644)
		p, err = os.ReadFile(filepath.Join(lib.resources, "recent.html"))
	}

	lib.recentT = template.Must(template.New("recent.html").Parse(string(p)))
}

func (lib *Library) authorpageTemplate() {

	p, err := os.ReadFile(filepath.Join(lib.resources, "author.html"))

	if os.IsNotExist(err) {
		os.WriteFile(filepath.Join(lib.resources, "author.html"), []byte(authorTemplate(lib)), 0644)
		p, err = os.ReadFile(filepath.Join(lib.resources, "author.html"))
	}
	lib.authorT = template.Must(template.New("author.html").Parse(string(p)))
}

func (lib *Library) bookpageTemplate() {

	p, err := os.ReadFile(filepath.Join(lib.resources, "book.html"))

	if os.IsNotExist(err) {
		os.WriteFile(filepath.Join(lib.resources, "book.html"), []byte(bookTemplate(lib)), 0644)
		p, err = os.ReadFile(filepath.Join(lib.resources, "book.html"))
	}
	lib.bookT = template.Must(template.New("book.html").Parse(string(p)))
}

func (lib *Library) categorypageTemplate() {

	p, err := os.ReadFile(filepath.Join(lib.resources, "category.html"))

	if os.IsNotExist(err) {
		os.WriteFile(filepath.Join(lib.resources, "category.html"), []byte(categoryTemplate(lib)), 0644)
		p, err = os.ReadFile(filepath.Join(lib.resources, "category.html"))
	}
	lib.categoryT = template.Must(template.New("book.html").Parse(string(p)))
}
