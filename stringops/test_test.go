package aozorafs

import (
	_ "embed" //embed
	"fmt"
	"strconv"
	"strings"
	"testing"
)

//go:embed stringops.go
var data string

var substr = "9" //data[:100]

var intstr []string

func init() {
	spdata = Split(data, "\n")

	for i := 0; i < 10000; i++ {

		intstr = append(intstr, itoStr(i, 10))

	}

}

func BenchmarkIndex(b *testing.B) {

	for range b.N {

		Index(data, substr)

	}
}
func BenchmarkIndexS(b *testing.B) {

	for range b.N {

		strings.Index(data, substr)

	}
}

func TestItoHex(t *testing.T) {

	fmt.Println(HexToI("FA10"))

}
func TestJoin(t *testing.T) {

	s1 := "https://www.orihasam.com"
	s2 := "aozora/test/index.html"

	fmt.Println(Split(s1, "/"))
	fmt.Println(HasPrefix(s1, "/"))

	fmt.Println(FilepathJoin(s1, s2))
	fmt.Println(URLJoin(s1, s2))

}

func TestLower(t *testing.T) {

	s := `ABCDEFGHIJKLMNOPQRSTUVWXYZ!@#123`

	sl := ToLower(s)

	fmt.Println(s, sl)

	fmt.Println(ToUpper(sl))

}
func TestTrim(t *testing.T) {

	s := `  ABCDEFGHIJKLMNOPQRSTUVWXYZ!@#123 `
	s = ` 　あさ `
	fmt.Println(TrimLeft(s, "ABC"))

	fmt.Println(TrimRight(s, "3"))

	fmt.Println("*" + s + "*")
	fmt.Println("*" + TrimSpace(s) + "*")

}

func TestBtoS(t *testing.T) {

	s := `ABCDEF abcdef`
	s = ` 　あさ `

	enc := BytesToPercent([]byte(s))

	dec := PercentDecode(enc)

	fmt.Println(s, enc, dec)

}

func TestIntOf(t *testing.T) {

	fmt.Println(65535, ItoHex(65535))

	return

	for i := 0; i > -1; i++ {

		s := ItoHex(i)

		j, err := HexToI(s)

		if err != nil {
			fmt.Println(i, s, err)
			break
		}

		fmt.Println(i, s, j)

		if i != j {
			break
		}

	}
	return

}

func TestNumerals(t *testing.T) {

	fmt.Println(numerals(10))

	fmt.Println(numerals(16))

	fmt.Println(numerals(2))

}

func TestBlock(t *testing.T) {

	s := "/aozora/orihasam/com/testing/index.html"

	elems := Split(s, "/")

	c := Join(elems, "/")

	fmt.Println(s)
	fmt.Println(c)

	if s != c {

		fmt.Println("fail!")
	}

	return

	reader := NewBlockReader(s, "/")

	fmt.Println(s)

	elem, _ := reader.Read()
	elem, _ = reader.Read()

	fmt.Println(elem, "***", reader.Unread())

	_, eof := reader.Read()

	fmt.Println(eof)

}

func TestReplace(t *testing.T) {

	s := "/aozora/orihasam/com/testing/index.html"

	fmt.Println(s, Replace(s, "/", "%", 3))

	s = "aozora/orihasam/com/testing/index.html"

	fmt.Println(s, Replace(s, "/", "%", 2))
}

func TestSplit(t *testing.T) {

	//	data = "/aozora/orihasam/com/testing/index.html"
	//	fmt.Println(data)
	//	fmt.Println(Join(Split(data, "\n"), "\n"))
	fmt.Println(data == Join(Split(data, "\n"), "\n"))
	return

}

func TestBuilder(t *testing.T) {

	//data = "/aozora/orihasam/com/testing/index.html"

	elems := Split(data, "/")

	bd := new(Builder)

	bs := new(strings.Builder)

	var out string

	for _, e := range elems {

		bd.WriteString(e)
		bs.WriteString(e)
		out = out + e
	}

	out1 := bd.String()
	out2 := bs.String()

	fmt.Println("out1,out2", out1 == out2)

	fmt.Println("out1,out", out1 == out)

	return

	fmt.Println(BytesToPercent([]byte(out)))

}

func BenchmarkBuilder(b *testing.B) {

	for range b.N {
		bd := new(Builder)

		for i := range spdata {

			bd.WriteString(spdata[i])

		}

		bd.String()
	}
	return

}
func BenchmarkSBuilder(b *testing.B) {

	for range b.N {
		bd := new(strings.Builder)

		for i := range spdata {

			bd.WriteString(spdata[i])

		}

		bd.String()
	}
}

func BenchmarkSplit(b *testing.B) {
	for range b.N {
		Split(data, "\n")
	}
}

func BenchmarkSSplit(b *testing.B) {
	for range b.N {
		strings.Split(data, "\n")
	}
}

func BenchmarkHasPrefix(b *testing.B) {

	prefix := data[:1000]
	b.ResetTimer()
	for range b.N {
		HasPrefix(data, prefix)
	}
}

func BenchmarkHasPrefixS(b *testing.B) {

	prefix := data[:1000]
	b.ResetTimer()
	for range b.N {
		strings.HasPrefix(data, prefix)
	}
}

var spdata []string

func BenchmarkJoin(b *testing.B) {

	for range b.N {
		Join(spdata, "\n")
	}

}

func BenchmarkSJoin(b *testing.B) {

	for range b.N {
		strings.Join(spdata, "\n")
	}
}

func BenchmarkBlock(b *testing.B) {

	br := NewBlockReader(data, "\n")

	for _, eof := br.Read(); !eof; _, eof = br.Read() {
	}

}

func BenchmarkBlockS(b *testing.B) {

	s := strings.Split(data, "\n")

	for _ = range s {
	}

}
func _TestItoA(b *testing.T) {

	for i := 0; i < 100000; i++ {

		s := itoStr(i, 16)

		j, err := strToInt(s, 16)

		fmt.Println(i, digitLen(16, i), s, j)

		if err != nil {
			fmt.Println("FAIL")
			return
		}

		if i != j {
			fmt.Println("FAIL")
			return
		}
	}
}

func BenchmarkStrtoA(b *testing.B) {

	for _, s := range intstr {

		strToInt(s, 10)

	}

}

func BenchmarkStrtoAS(b *testing.B) {

	for _, s := range intstr {
		strconv.ParseInt(s, 10, 64)

	}

}

func BenchmarkItoA(b *testing.B) {

	for i := 0; i < 100000; i++ {
		itoStr(i, 10)
	}
}

func BenchmarkItoAS(b *testing.B) {

	var i int64
	for i = 0; i < 100000; i++ {
		strconv.FormatInt(i, 10)
	}
}
