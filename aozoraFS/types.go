package aozorafs

import (
	"io/fs"
	"text/template"
	"time"
)

type LibFile interface {
	fs.File
	Write(b []byte) (n int, err error)
	// Sync() error
}

type LibFS interface {
	fs.StatFS
	CreateFile(name string, data []byte) (fs.File, error)
	CreateEphemeral(name string, data []byte) (fs.File, error)
	Exists(name string) bool
	RemoveAll()
	Path() string
}

// Library is stores basic information. Root is the path to the original
// aozora file tree. Cache is path to the storage directory for server. Resources
// holds the templates, css.
// Catalog is for stroing Library information.
type Library struct {
	src string
	//cache string
	cache         LibFS
	root          string
	resources     string
	booklist      []*Record
	booksByID     map[string][]*Record
	booksByAuthor map[string][]*Record
	booksByDate   []*Record
	authorsSorted []*Record
	posOfAuthor   map[string]int
	indexT,
	authorT,
	bookT,
	categoryT,
	recentT,
	randomT,
	searchresultT,
	readingT *template.Template
	updating      bool
	kids          bool
	strict        bool
	lastUpdated   time.Time
	checkInterval time.Duration
	Categories    map[string]string
	nextrandom    int
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
	Categories      [][3]string
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
	Contributors    []ContribRole
	consolidated    bool
}

// ContribRole is for storing various contributors to a book.
type ContribRole struct {
	Role     string
	AuthorID string
	B        *Record
}
