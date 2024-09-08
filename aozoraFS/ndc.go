package aozorafs

import (
	_ "embed" //for embedding data
	"strings"
)

//go:embed ndc.data
var ndcdata string

func ndcmap() map[string]string {

	ndc := make(map[string]string)

	lines := strings.Split(ndcdata, "\n")

	for _, l := range lines {

		d := strings.Split(l, ",")

		if len(d) != 2 {
			continue
		}

		ndc[d[0]] = d[1]

	}

	for key, val := range ndc {

		ndc[val] = key

	}

	return ndc
}
