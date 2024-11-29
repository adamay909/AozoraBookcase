package aozorafs

import "strings"

func removeBOM(s string) string {

	sb := []byte(s)

	for i, b := range []byte{239, 187, 191} {

		if b != sb[i] {
			return s
		}

	}

	return string(sb[3:])
}

func getCells(line string) (cells []string) {

	nextCell := nextCSVcell(line)

	for nc, end := nextCell(); !end; nc, end = nextCell() {

		cells = append(cells, nc)

	}

	return cells

}

func nextCSVcell(line string) func() (string, bool) {

	start, end, oldstart := 0, 0, 0

	return func() (string, bool) {

		if start == len(line) {
			return "", true
		}

		oldstart = start

		if strings.HasPrefix(line[start:], `"`) {

			end = strings.Index(line[start+1:], `",`)

			if end != -1 {

				start = start + end + 3

				return line[oldstart+1 : oldstart+end+1], false

			}
		} else {

			end = strings.Index(line[start:], ",")

			if end != -1 {

				start = start + end + 1

				return line[oldstart : oldstart+end], false
			}
		}

		start = len(line)

		return line[oldstart:], false

	}
}
