package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestSeekString(t *testing.T) {
	var s string
	s = `asdf
	asdf11
	dd`
	r := strings.NewReader(s)
	r.Seek(1, 0)
	buf := bufio.NewReader(r)
	bb, _, _ := buf.ReadLine()
	fmt.Println(string(bb))

}

/*func TestFmtScanf(t *testing.T) {
	s := "asdf 123"
	var s1 string
	var i1 int32
	_, e := fmt.Sscanf(s, "%*s %d", &s1, &i1)
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}
	if i1 != 123 {
		t.Fail()
	}
}
*/
func TestCanDecode(t *testing.T) {
	fpath := "./testfiles/1升液体.RPT"
	parser := new(GammaRPT)
	err := parser.Parse(fpath)
	if err != nil {
		t.Fatal(err)
	}

	//vol
	if parser.SVol != 0.001 {
		t.Fatal(parser.SVol)
	}
	//title
	if len(parser.STitle) <= 1 {
		t.Fatal(parser.STitle)
	}
	//time
	if parser.AcqStartTime.Year() != 2008 {
		t.Fatal(parser.AcqStartTime)
	}
	//Ident
	if len(parser.Nuclides) < 2 {
		t.Fatal("只解析到:%v个数据", len(parser.Nuclides))
	}
	if parser.Nuclides["K-40"].Act!=2.772275e4{
		t.Fatal(parser.Nuclides)
	}

}

func TestSplit(t *testing.T) {
	s := "       K-40       0.952      2.772275E+004   3.473816E+003"
	var (
		s1         string
		i1, i2, i3 float32
	)
	fmt.Sscanf(s, "%s %f %f %f", &s1, &i1, &i2, &i3)
}

func TestMDALen(t *testing.T){
	t.Fatal(len("MN-54      834.83    99.97  8.924E+002  8.92E+002 -5.105E+002  4.462E+002"))
}


