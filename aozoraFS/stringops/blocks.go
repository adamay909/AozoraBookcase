package aozorafs

import "strings"

// nextBlock returns *two* functions:
// readFunc: can be used to find the next block of string with data where each
// block is separated by the separator.
// unreadFunc: returns the unread portion of the string where every read is
// understood to consume the separator. E.g., if the separator is "\n",
// the unread part does not start with "\n" because intuitively the unread
// part is the set of remaining lines.
func nextBlock(data string, separator string) (readFunc func() (string, bool), unreadFunc func() string) {

	start, end := 0, -len(separator)

	readFunc = func() (string, bool) {

		start = end + len(separator)

		if start >= len(data) {
			return "", true
		}

		end = start + strings.Index(data[start:], separator)

		if end < start {
			end = len(data)
			return data[start:], false
		}

		return data[start:end], false
	}

	unreadFunc = func() string {

		ostart := end + len(separator)

		if ostart >= len(data) {
			return ""
		}

		start = len(data)

		return data[ostart:]
	}

	return readFunc, unreadFunc

}

// BlockReader can be used to read blocks of text.
type BlockReader struct {
	nextBlock func() (string, bool)
	remainder func() string
	separator string
}

// NewBlockReader returns a new BlockReader with data
// as the underlying data and the specified separator.
func NewBlockReader(data, separator string) *BlockReader {

	b := new(BlockReader)

	b.separator = separator

	b.nextBlock, b.remainder = nextBlock(data, separator)

	return b

}

// Read returns the next block of text. The boolean
// return value is true iff. the end of the underlying
// data has been reached.
func (b *BlockReader) Read() (string, bool) {

	return b.nextBlock()

}

// Unread returns the unread portion of the underlying data of b.
func (b *BlockReader) Unread() string {

	return b.remainder()

}

// NewLineReader returns a BlockReader with "\n" as separator.
func NewLineReader(data string) *BlockReader {

	return NewBlockReader(data, "\n")

}
