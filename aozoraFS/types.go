package aozorafs

import (
	"html/template"
	"time"
)

// Library is stores basic information. Root is the path to the original
// aozora file tree. Cache is path to the storage directory for server. Resources
// holds the templates, css.
// Catalog is for stroing Library information.
type Library struct {
	src       string
	cache     string
	resources string
	root      string
	//	prefix    string
	//	catalog   records
	booklist []*Record
	//	size      int
	indexT,
	authorT,
	bookT,
	recentT *template.Template
	updater       chan bool
	kids          bool
	strict        bool
	lastUpdated   time.Time
	checkInterval time.Duration
}

// records is for storing index cards.
type records map[string]*Record

// Record is for storing information about individual books.
type Record struct {
	count           int
	BookID          string
	Title           string
	TitleY          string
	TitleSort       string
	Subtitle        string
	SubtitleY       string
	SubtitleSort    string
	OriginalTitle   string
	PublDate        string
	NDC             string
	Category        string
	KanaZukai       string
	WorkCopyright   string
	FirstAvailable  string
	ModTime         string
	AuthorID        string
	NameSei         string
	NameMei         string
	NameSeiY        string
	NameMeiY        string
	NameSeiSort     string
	NameMeiSort     string
	NameSeiR        string
	NameMeiR        string
	Role            string
	DoBirth         string
	DoDeath         string
	AuthorCopyright string
	URI             string
	Kids            bool
}
