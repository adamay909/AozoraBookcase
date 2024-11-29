package aozorafs

func IsAsciiEncoded(enc string) bool {

	return len(enc) == len([]rune(enc))

}
