package aozorafs

/*RefreshBooklist checks for updates and if necessary refreshes the database periodically as specified by lib.checkInt. If lib.checkInt <=0, then database is never refreshed. */
/*func (lib *Library) RefreshBooklist() {

	update := func() {
		lib.UpdateDB()
		lib.UpdateBooklist()
		lib.updatePages()
	}

	if lib.checkInterval <= 0 {
		return
	}
	for {
		time.Sleep(lib.checkInterval)

		if lib.UpstreamUpdated(lib.lastUpdated) {
			update()
		}
	}
}
*/
/*
// LoadBooklist adds (and possibly updates) the list of books for lib.
func (lib *Library) _LoadBooklist() {

	fi, err := lib.cache.Stat("aozoradata.zip")

	switch {
	case err != nil:
		log.Println("no local aozora database")
		lib.UpdateDB()

	default:
		if lib.UpstreamUpdated(fi.ModTime()) {
			lib.UpdateDB()
		}

	}

	lib.UpdateBooklist()

	lib.updatePages()

	go lib.RefreshBooklist()

}
*/
/*UpstreamUpdated reports whether the upstream database has been updated since it was last updated locally.
 */
/*func (lib *Library) _UpstreamUpdated(t time.Time) bool {

	loc, err := url.JoinPath(lib.src, "/index_pages", "list_person_all_extended_utf8.zip")

	if err != nil {
		log.Println(err)
		return false
	}

	path, _ := url.Parse(loc)

	r := getHeader(path)

	m, err := time.Parse(time.RFC1123, get(r, "Last-Modified"))

	return m.After(t)
}
*/
/*UpdateDB downloads the database from upstream.*/
/*func (lib *Library) _UpdateDB() {

	pathString, err := url.JoinPath(lib.src, "/index_pages", "list_person_all_extended_utf8.zip")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("requesting db", pathString)

	path, _ := url.Parse(pathString)

	data := download(path)

	f, err := lib.cache.CreateFile("aozoradata.zip", data)

	defer f.Close()

	if err != nil {
		log.Println(err)
		return
	}
}
*/
/*UpdateBooklist updates the booklist of lib from the locally available database.*/
/*func (lib *Library) _UpdateBooklist() {

	log.Println("updating book list")

	lib.updating = true

	zf, err := zipfs.OpenZipArchive(lib.cache, "aozoradata.zip")

	defer zf.CloseArchive()

	if err != nil {
		log.Println(err)
		return
	}

	zd := zf.ReadMust("list_person_all_extended_utf8.csv")

	lib.getBooklist(zd)
	lib.setupAuthorsList()
	lib.updating = false
	return
}
*/

/*
func (lib *Library) updatePages() {

		log.Println("Updating pages")

		lib.removePages(`index.html`, `recent.html`)
		lib.removePages(lib.allUpdatedPages()...)
		lib.lastUpdated = time.Now()
		log.Println("pages updated")
	}

func (lib *Library) _allUpdatedPages() (list []string) {

		for _, b := range lib.booklist {

			var fnames []string

			fnames = append(fnames, `authors/author_`+b.AuthorID+`.html`)
			fnames = append(fnames, `books/book_`+b.AuthorID+`_`+b.BookID+`.html`)

			for _, c := range b.Categories {

				fnames = append(fnames, `categories/ndc_`+c[0]+`.html`)
				fnames = append(fnames, `categories/ndc_`+c[1]+`.html`)

			}

			for _, f := range fnames {
				info, err := lib.cache.Stat(f)
				if err == nil {

					if !info.ModTime().Before(lib.lastUpdated) {
						log.Println(f, "needs updating")
						list = append(list, f)
					}
				}
			}

		}

		return list
	}

func (lib *Library) removePages(pages ...string) {


	lib.cache.RemoveAll()

	return
}
*/

/*
var getHeader func(path *url.URL) map[string][]string

func SetHeader(f func(*url.URL) map[string][]string) {

	getHeader = f

}

func get(header map[string][]string, key string) string {

	r, ok := header[key]

	if !ok {
		return ""
	}

	return r[0]
}

*/
