package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	str "strings"
	"time"
)

type GammaRPT struct {
	SVol         float32
	STitle       string
	SType        string
	AcqStartTime time.Time
	LiveSeconds  int32
	RptGenTime   time.Time
	Nuclides     map[string]NuclideActivity
	file         *string
	fReader      *str.Reader
	bufReader    *bufio.Reader
}

type NuclideActivity struct {
	Name  string
	Act   float32
	IsLLD bool
}

func (self *GammaRPT) Parse(filePath string) error {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	file := string(f)
	self.file = &file
	self.fReader = str.NewReader(file)
	self.bufReader = bufio.NewReader(self.fReader)
	volStart := str.Index(file, "Sample Size")
	if volStart < 0 {
		return errors.New("未找到样品体积")
	}
	self.fReader.Seek(int64(volStart), 0)
	vol,err := self.bufReader.ReadLine()
	if err != nil{
		panic(err)
	}
	return nil
}
