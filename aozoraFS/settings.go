package aozorafs

func (lib *Library) SetSrc(name string) {

	lib.src = name

	return

}

func (lib *Library) SetCache(fs LibFS) {

	lib.cache = fs

	return

}

func (lib *Library) SetKids(k bool) {

	lib.kids = k

	return

}
