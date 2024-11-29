package aozorafs

import (
	"strconv"
	"strings"
)

type customerror string

// ItoHex returns n in hexadecimal notation
func ItoHex(n int) string {

	return strconv.FormatInt(int64(n), 16)
}

func HexToI(h string) (int64, error) {

	return strconv.ParseInt(h, 16, 64)

}

// digitLen returns the number of digts of n
// in a base b notation. E.g., digitLen(10, 999)= 3,
// digitLen(16, 251) = 2
func DigitLen(b, n int) (l int) {

	if n < b {
		return 1
	}

	for l, n = 1, n/b; n != 0; l++ {
		n = n / b
	}
	return l
}

func BytesToStr(b []byte) string {

	return strings.Join(bytesToHex(b), "")

}

func BytesToPercent(b []byte) string {

	return "%" + strings.Join(bytesToHex(b), "%")
}

func PercentDecode(enc string) string {

	var rdec []byte

	renc := []byte(enc)

	for i := 0; i < len(renc); i++ {

		if renc[i] != '%' {

			rdec = append(rdec, renc[i])

			continue
		}

		if len(enc)-i < 3 {

			rdec = append(rdec, renc[i:]...)

			break
		}

		r, err := HexToI(string(renc[i+1 : i+3]))

		if err != nil {

			rdec = append(rdec, renc[i])

			continue

		}

		rdec = append(rdec, byte(int8(r)))

		i = i + 2
	}

	return string(rdec)
}

func bytesToHex(b []byte) []string {

	r := make([]string, len(b))

	for i, e := range b {

		r[i] = ItoHex(int(e))

	}

	return r

}
