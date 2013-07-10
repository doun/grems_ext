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

func TestFmtScanf(t *testing.T) {
	s := "asdf 123"
	var s1 string
	var i1 int32
	_, e := fmt.Sscanf(s, "%*s %d",&s1, &i1)
	if e!= nil{
		fmt.Println(e)
		t.Fail()
	}
	if i1 != 123 {
		t.Fail()
	}
}

func TestCanDecodeVol(t *testing.T) {
	fpath := "./testfiles/1升液体.RPT"
	parser := new(GammaRPT)
	parser.Parse(fpath)
	if parser.SVol != 0.001 {
		t.Fail()
	}
}
