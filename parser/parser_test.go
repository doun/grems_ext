package main

import (
	"testing" 
	"strings"
	"bufio"
	"fmt"
)

func TestSeekString(t *testing.T) {
	var s string
	s = `asdf
	asdf11
	dd`
	r := strings.NewReader(s)
	r.Seek(1,0)
	buf := bufio.NewReader(r)
	bb,_,_ := buf.ReadLine()
	fmt.Println(string(bb))

}

