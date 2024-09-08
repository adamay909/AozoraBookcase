package main

import (
	"fmt"
	"testing"
)

func TestSort(t *testing.T) {

	s := []string{
		"a",
		"bc",
		"def",
	}

	sortPrefixes(s)

	for _, e := range s {
		fmt.Println(e)
	}
}
