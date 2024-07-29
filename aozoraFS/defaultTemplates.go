package aozorafs

import (
	_ "embed" //use this to hard code templates
)

//go:embed resources/defaultcss.css
var fileDefaultcss string

func defaultCSS(lib *Library) string {
	return fileDefaultcss
}

//go:embed resources/simplecss.css
var fileSimplecss string

func simpleCSS(lib *Library) string {
	return fileSimplecss
}

//go:embed resources/index.html
var fileIndexhtml string

func indexTemplate(lib *Library) string {
	return fileIndexhtml
}

//go:embed resources/recent.html
var fileRecenthtml string

func recentTemplate(lib *Library) string {
	return fileRecenthtml
}

//go:embed resources/author.html
var fileAuthorhtml string

func authorTemplate(lib *Library) string {
	return fileAuthorhtml
}

//go:embed resources/book.html
var fileBookhtml string

func bookTemplate(lib *Library) string {

	return fileBookhtml
}

//go:embed resources/category.html
var fileCategoryhtml string

func categoryTemplate(lib *Library) string {

	return fileCategoryhtml

}

//go:embed resources/searchForm.html
var fileSearchhtml string

func searchFormTemplate() string {
	return fileSearchhtml
}

//go:embed resources/searchResults.html
var fileSearchresultshtml string

func searchResultsTemplate() string {
	return fileSearchresultshtml
}
