package aozorafs

import (
	_ "embed" //for embedding data
	"strings"

	str "github.com/adamay909/AozoraBookcase/aozoraFS/stringops"
)

//go:embed ndc.data
var ndcdata string

func ndcmap() map[string]string {

	ndc := make(map[string]string)

	lines := str.NewLineReader(ndcdata)

	for {

		d, eof := lines.Read()

		if eof {
			break
		}

		idx := strings.Index(d, ",")

		ndc[d[:idx]] = d[idx+1:]

	}

	for key, val := range ndc {

		ndc[val] = key

	}

	return ndc
}
