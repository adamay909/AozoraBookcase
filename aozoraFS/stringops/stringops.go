/*
Package stringops provides some functionality to assist with working with strings.
*/
package stringops

func IsAsciiEncoded(enc string) bool {

	return len(enc) == len([]rune(enc))

}
