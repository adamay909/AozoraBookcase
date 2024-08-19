package server

import (
	_ "embed" //for embedding resources

	aozorafs "github.com/adamay909/AozoraBookcase/aozoraFS"
)

//use this to hard code templates

//go:embed resources/defaultcss.css
var fileDefaultcss string

//go:embed resources/random.html
var randombookhtml string

//go:embed resources/index.html
var fileIndexhtml string

//go:embed resources/recent.html
var fileRecenthtml string

//go:embed resources/author.html
var fileAuthorhtml string

//go:embed resources/book.html
var fileBookhtml string

//go:embed resources/category.html
var fileCategoryhtml string

//go:embed resources/random.html
var randomBookhtml string

//go:embed resources/searchForm.html
var fileSearchhtml string

//go:embed resources/searchResults.html
var fileSearchresultshtml string

func SetTemplates() {

	_ = aozorafs.NewLibrary()

	aozorafs.FileDefaultcss = fileDefaultcss

	aozorafs.Randombookhtml = randombookhtml

	aozorafs.FileIndexhtml = fileIndexhtml

	aozorafs.FileRecenthtml = fileRecenthtml

	aozorafs.FileAuthorhtml = fileAuthorhtml

	aozorafs.FileBookhtml = fileBookhtml

	aozorafs.FileCategoryhtml = fileCategoryhtml

	aozorafs.RandomBookhtml = randomBookhtml

	aozorafs.FileSearchhtml = fileSearchhtml

	aozorafs.FileSearchresultshtml = fileSearchresultshtml

}
