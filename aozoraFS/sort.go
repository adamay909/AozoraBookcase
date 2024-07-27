package aozorafs

import (
	"log"
	"sort"
	"strings"
	"time"
)

type lessFunc func(l []*Record) func(i, j int) bool

// sortList sorts list of *Record using the less function.
func sortList(list []*Record, f lessFunc) {
	sort.Slice(list, f(list))
}

// ************
// The lesFunc collection
// ***********
func byAuthor(l []*Record) func(i, j int) bool {

	return func(i, j int) bool {

		if l[i].fullNameS() < l[j].fullNameS() {
			return true
		}
		if l[i].fullNameS() > l[j].fullNameS() {
			return false
		}

		if l[i].AuthorID < l[j].AuthorID {
			return true
		}
		if l[i].AuthorID > l[j].AuthorID {
			return false
		}

		if l[i].TitleSort < l[j].TitleSort {
			return true
		}
		return false
	}
}

func byBookID(l []*Record) func(i, j int) bool {
	return func(i, j int) bool { return l[i].BookID < l[j].BookID }
}

func byTitle(l []*Record) func(i, j int) bool {
	return func(i, j int) bool { return l[i].TitleSort < l[j].TitleSort }
}

func byAuthorName(l []*Record) func(i, j int) bool {
	return func(i, j int) bool { return l[i].fullNameS() < l[j].fullNameS() }
}

func byAuthorID(l []*Record) func(i, j int) bool {
	return func(i, j int) bool { return l[i].AuthorID < l[j].AuthorID }
}

func byModTime(l []*Record) func(i, j int) bool {

	return func(i, j int) bool {

		itime, _ := time.Parse(time.DateOnly, l[i].ModTime)

		jtime, _ := time.Parse(time.DateOnly, l[j].ModTime)

		return itime.Before(jtime)

	}
}

func byAvailableDate(l []*Record) func(i, j int) bool {

	return func(i, j int) bool {

		itime, err := time.Parse(time.DateOnly, l[i].FirstAvailable)

		if err != nil {
			log.Println(err)
		}

		jtime, err := time.Parse(time.DateOnly, l[j].FirstAvailable)

		if err != nil {
			log.Println(err)
		}

		return itime.After(jtime)

	}
}

// 1=chosha
// 2=hensha
// 3=honyakusha
// 4=kouetusha
// 5=sonota
func byRole(l []ContribRole) func(i, j int) bool {

	return func(i, j int) bool {
		return numOf(l[i].Role) < numOf(l[j].Role)
	}
}

func numOf(role string) int {

	switch strings.TrimSpace(role) {

	case "著者":
		return 1

	case "編者":
		return 2

	case "翻訳者":
		return 3

	case "校訂者":
		return 4

	case "その他":
		return 5

	}

	return 0
}
