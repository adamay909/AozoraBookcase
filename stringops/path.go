package aozorafs

import "strings"

func FilepathBase(path string) string {

	if path == "" {
		return "."
	}

	path = removeTrailingSlash(path)

	if path == "" {
		return `/`
	}

	elem := strings.Split(path, `/`)

	return elem[len(elem)-1]

}

func FilepathExt(path string) string {

	base := FilepathBase(path)

	if base == "." {
		return ""
	}

	if base == "/" {
		return ""
	}

	idx := strings.LastIndex(base, `.`)

	if idx == -1 {
		return ""
	}

	return base[idx:]

}

func FilepathDir(path string) string {

	idx := strings.LastIndex(path, "/")

	if idx == -1 {
		return "."
	}

	if path == "." {
		return "."
	}

	if path == ".." {
		return ".."
	}

	path = removeTrailingSlash(path[:idx])

	if path == "" {
		path = "/"
	}

	return path

}

func FilepathJoin(elem ...string) string {

	return FilepathClean(strings.Join(elem, "/"))

}

func FilepathClean(path string) string {

	elem := strings.Split(path, "/")

	nelem := make([]string, len(elem)+1)
	c := 0

	if strings.HasPrefix(path, "/") {
		nelem[0] = ""
		c++
	}

	for i, e := range elem {

		if e == "" {
			continue
		}

		if e == "." {
			continue
		}

		if i+1 < len(elem) {

			if elem[i+1] == ".." {
				continue
			}
		}

		nelem[c] = e
		c++

	}

	return strings.Join(nelem[:c], "/")
}

func URLJoin(base string, elem ...string) string {

	return base + FilepathJoin(elem...)

}

func removeTrailingSlash(path string) string {

	return strings.TrimRight(path, "/")
}
